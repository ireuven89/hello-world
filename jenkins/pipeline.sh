pipeline {
    agent any
    stages {
      stage('Clone') {
        steps{
          cd /home/projects
          git clone ''
          cd hello-world/backend
        }
      }
      stage('Build') {
            steps {
                go build .
            }
        }
      stage('Test') {
            steps {
              cd /tests
              go test ./...
            }
        }
      stage('Deploy') {
            steps {
               docker build .
               docker push '/'
            }
        }
    }
}