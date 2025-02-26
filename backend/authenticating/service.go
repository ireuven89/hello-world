package authenticating

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/ireuven89/hello-world/backend/authenticating/model"
	"github.com/ireuven89/hello-world/backend/consumerring"
	"github.com/ireuven89/hello-world/backend/producerring"
	p_model "github.com/ireuven89/hello-world/backend/producerring/model"
	"github.com/ireuven89/hello-world/backend/utils"
)

type Service interface {
	Register(username, password string) error
	Login(username, password string) (string, error)
	VerifyToken(tokenString string) (string, error)
	Migrate(ctx context.Context, startTime time.Time) error
}

type AuthRepo interface {
	Save(username, password string) error
	Find(username string) (model.User, error)
	FindAll(page model.Page) ([]model.User, error)
	Delete(id string) error
}

// AuthService is the core authenticating service
type AuthService struct {
	userStore AuthRepo
	logger    *zap.Logger
	mongo     *mongo.Database
	producer  *kafka.Producer
}

var jwtSecretKey = []byte("your_secret_key")
var batchSize = int64(100)

// NewAuthService creates a new AuthService
func NewAuthService(userStore AuthRepo, logger *zap.Logger) *AuthService {
	return &AuthService{userStore: userStore, logger: logger}
}

// Register registers a new user
func (service *AuthService) Register(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		service.logger.Error("failed to register with error", zap.Error(err))
		return err
	}
	return service.userStore.Save(username, string(hashedPassword))
}

// Login authenticates a model and returns a JWT token
func (service *AuthService) Login(username, password string) (string, error) {
	user, err := service.userStore.Find(username)
	if err != nil {
		service.logger.Error("failed to login with error", zap.Error(err))
		return "", err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // Token expires in 24 hours
	})

	return token.SignedString(jwtSecretKey)
}

// VerifyToken verifies and decodes a JWT token
func (service *AuthService) VerifyToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return "", errors.New("username not found in token")
	}

	return username, nil
}

func (service *AuthService) Migrate(ctx context.Context, startTime time.Time) error {
	execuore := consumerring.NewExecutor("my-migration", service.execute, service.rollback)
	cfg, err := utils.LoadConfig("authenticating", os.Getenv("local"))

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"broker": cfg.Databases["kafka"],
	})
	if err != nil {
		return err
	}
	worker := consumerring.NewWorkerService(consumer, execuore, 5)

	if err != nil {
		return err

	}
	producer := producerring.NewMigrationService(service.mongo, service.producer, "migration")
	err = producer.CreateMigration(ctx, p_model.Migration{
		CreatedAt:     time.Now(),
		MigrationName: "printing",
	})
	if err != nil {
		return err
	}

	delay := time.Until(startTime)

	if delay <= 0 {
		go worker.ProcessTasks(ctx, "migrations")
	}

	go func() {
		timer := time.NewTimer(delay)
		select {
		case <-timer.C:
			fmt.Println("Starting migration:", "printing")
			go worker.ProcessTasks(ctx, "migrations")
		case <-ctx.Done():
			timer.Stop()
		}
	}()

	return nil
}

func (service *AuthService) execute(sqlUser interface{}) error {
	ctx := context.Background()
	user, ok := sqlUser.(model.User)
	if !ok {
		return errors.New("failed to create user")
	}

	res, err := service.mongo.Collection("users").InsertOne(ctx, bson.M{
		"name":     user.Username,
		"password": user.Password,
	})

	if err != nil {
		service.logger.Error("failed inserting", zap.Error(err))
		return err
	}

	if res.InsertedID == "" {
		service.logger.Error("failed inserting")
		return err
	}

	err = service.userStore.Delete(user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (service *AuthService) rollback(sqlUser interface{}) error {
	ctx := context.Background()
	user, ok := sqlUser.(model.User)

	if !ok {
		return errors.New("failed parsing")
	}
	res, err := service.mongo.Collection("users").DeleteOne(ctx, bson.M{"name": user.Username})

	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("failed deleting")
	}

	return nil
}

// createTasks - takes the id from db and converts them
func (service *AuthService) createTasks(producer *producerring.MigrationService) error {
	page := model.Page{Page: 0, PageSize: batchSize}

	for {
		result, err := service.userStore.FindAll(page)
		if err != nil {
			return err
		}

		if len(result) == 0 {
			service.logger.Info("published all tasks")
			break
		}
		for _, user := range result {
			err = producer.PublishTask(p_model.MigrationTask{
				Params: map[string]interface{}{
					"id":       user.Username,
					"password": user.Password,
				},
			})

			if err != nil {
				service.logger.Error("failed publishing task")
			}
		}

		page.Page++
	}

	return nil
}
