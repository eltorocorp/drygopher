    node {
        checkout scm
        def image = docker.build("image:${env.BUILD_ID}")
        image.inside {
            sh 'make'
        }
    }