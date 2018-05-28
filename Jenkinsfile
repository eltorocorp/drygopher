pipeline {
    agent {
        dockerfile true
    }
    stages {
        stage('build') {
            steps {
                echo 'building...'
                sh 'go/src/github.com/eltorocorp/drygopher/make build'
            }
        }
        stage('test') {
            steps {
                echo 'testing...'
                sh 'go/src/github.com/eltorocorp/drygopher/make test'
            }
        }

    }
}