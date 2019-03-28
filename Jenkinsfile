node{
  checkout scm

    environment {
            DOCKER_LOGIN     = credentials('docker-login')
            DOCKER_PASSWORD = credentials('docker-password')
        }
        stages {
            stage('Unit Test') {
                steps {
                   sh """ docker run --rm \${PWD}:/go/src/github.com/thiagotrennepohl/fortune-backend golang go test ./..."""
                }
            }
            stage('Docker build') {
                steps {
                sh """
                    docker login -u \${DOCKER_LOGIN} -p \${DOCKER_PASSWORD}
                    docker build -t . thiagotr/fortune-backend .
                    docker push thiagotr/fortune-backend
                    """
                }
            }
        }
    }

}