pipeline {
    agent {
        docker {
            image 'golang:1.10'
            reuseNode true
        }
    }
    stages {
        stage('Prepare') {
            steps {
                echo 'Preparing build environment...'
                sh 'go get -u github.com/vektra/mockery/.../'
                sh 'go get -u github.com/golang/dep/cmd/dep'
                sh 'mkdir -p $GOPATH/src/drygopher && mv $(pwd) $GOPATH/src/drygopher'
            }
        }
        stage('Build') {
            steps {
                echo 'Building...'
                sh 'ls -a $GOPATH/src/drygopher'
                sh 'cd $GOPATH/src/drygopher && make build'
            }
        }
        stage('Test') {
            steps {
                echo 'Testing...'
                sh 'cd $GOPATH/src/drygopher && make test'
            }
        }
    }
}