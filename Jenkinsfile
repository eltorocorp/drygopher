pipeline {
    agent {
        docker {
            image 'golang:1.10'
            args '-v .:/go/src/github.com/eltorocorp/drygopher:rw -u root'
        }
    }
    stages {
        stage('Prepare') {
            steps {
                echo 'Preparing build environment...'
                sh 'go get github.com/vektra/mockery/.../'
                sh 'go get github.com/golang/dep/cmd/dep'
            }
        }
        stage('Build') {
            steps {
                echo 'Building...'
                sh 'cd /go/src/github.com/eltorocorp/drygopher && make build'
            }
        }
        stage('Test') {
            steps {
                echo 'Testing...'
                sh 'cd /go/src/github.com/eltorocorp/drygopher && make test'
            }
        }
    }
}