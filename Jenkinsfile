pipeline {
    agent {
        docker {
            image 'golang:1.10'
            reuseNode true
            customWorkspace '/var/lib/jenkins/workspace/drygopher'
            args '-v /var/lib/jenkins/workspace/drygopher:/go/src/github.com/eltorocorp/drygopher'
        }
    }
    stages {
        stage('Build') {
            steps {
                echo 'Building...'
                sh 'cd /go/src/github.com/eltorocorp/drygopher && make build'
            }
        }
        stage('Test') {
            steps {
                echo 'Testing...'
                sh 'cd /go/src/github.com/eltorocorp/drygopher && make test || cat coverage.out'
            }
        }
    }
}