pipeline {
    agent {
        docker {
            image 'golang:1.10'
        }
    }
    stages {
        stage('setup') {
            steps {
                echo 'setup...'
                script {
                    def workspace = pwd()
                }
                sh "GOPATH = ${workspace}" 
            }
        }
        stage('build') {
            steps {
                echo 'building...'
                sh 'go env'
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