pipeline {
    agent {
        docker {
            image 'golang:1.10'
            reuseNode true
            args '-v env[WORKSPACE]/env[JOB_NAME]_env[GIT_BRANCH]-env.[BUILD_ID]:/go/src/drygopher'
        }
    }
    stages {
        stage('Prepare') {
            steps {
                echo 'Preparing build environment...'
                sh 'go get -u github.com/vektra/mockery/.../'
                sh 'go get -u github.com/golang/dep/cmd/dep'
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