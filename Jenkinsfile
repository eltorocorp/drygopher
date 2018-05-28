pipeline {
    agent {
        docker {
            image 'golang:1.10'
            customWorkspace '/var/lib/jenkins/workspace/drygopher'
            args '-v /var/lib/jenkins/workspace/drygopher:/go/src/github.com/eltorocorp/drygopher:rw -u root'
        }
    }
    options {
        timestamps()
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
    post {
        always {
            dir('.') {
                deleteDir()
            }
        }
    }
}