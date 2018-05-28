    node {
        checkout scm
        def image = docker.build("image:${env.BUILD_ID}").withRun('-u root')
        image.inside {
            sh 'make build'
        }
    }