package migration

/*
func TestMigrationService_Up(t *testing.T) {
	service := NewMigrationService(3)

	successTask := model.MigrationTask{
		Name: "SuccessTask",
		Execute: func(interface{}) error {
			return nil
		},
		Rollback: func(interface{}) error {
			return nil
		},
	}

	failTask := model.MigrationTask{
		Name: "FailTask",
		Execute: func(interface{}) error {
			return errors.New("execution failed")
		},
		Rollback: func(interface{}) error {
			return nil
		},
	}

	service.AddTask(successTask)
	service.AddTask(failTask)

	err := service.Up()
	if err == nil || err.Error() != "execution failed" {
		t.Fatalf("Expected execution failure, got %v", err)
	}
}

func TestMigrationService_Down(t *testing.T) {
	service := NewMigrationService(3)

	successTask := model.MigrationTask{
		Name: "SuccessTask",
		Execute: func(interface{}) error {
			return nil
		},
		Rollback: func(interface{}) error {
			return nil
		},
	}

	failTask := model.MigrationTask{
		Name: "FailTask",
		Execute: func(interface{}) error {
			return nil
		},
		Rollback: func(interface{}) error {
			return errors.New("rollback failed")
		},
	}

	service.AddTask(successTask)
	service.AddTask(failTask)

	err := service.Down()
	if err == nil || err.Error() != "rollback failed" {
		t.Fatalf("Expected rollback failure, got %v", err)
	}

}
*/
