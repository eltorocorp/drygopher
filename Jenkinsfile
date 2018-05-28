pipeline {
    agent {
        dockerfile true
    }
    stages {
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