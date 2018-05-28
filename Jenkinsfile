pipeline {
    agent {
        docker {
            image 'golang:1.10'
        }
    }
    environment {
        GOPATH = '${PWD}'
    }
    stages {
        stage('setup') {
            steps {
                echo 'setup...'
                sh 'pwd'
                sh 'go env'
            }
        }
        stage('build') {
            steps {
                echo 'building...'
                sh 'ls /go/src/github.com/eltorocorp/drygopher'
                sh 'cd /go/src/github.com/eltorocorp/drygopher && make build'
            }
        }
        stage('test') {
            steps {
                echo 'testing...'
                sh 'cd /go/src/github.com/eltorocorp/drygopher && make test'
            }
        }

    }
}