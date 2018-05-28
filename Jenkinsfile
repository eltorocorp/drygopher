#!/usr/bin/env groovy

// this will start an executor on a Jenkins agent with the docker label
node {
    String applicationName = "drygopher"
    String buildNumber = "${env.BUILD_NUMBER}"
    String goPath = "/go/src/github.com/eltorocorp/${applicationName}"

    // stage('Checkout from GitHub') {
    //     checkout scm
    // }

    docker.image("golang:1.10").inside("-v ${pwd()}:${goPath} -u root") {
        stage 'PreBuild'
        sh "cd ${goPath} && make prebuild"

        stage 'Build'
        sh "cd ${goPath} && make build"

        stage 'Test'
        sh "cd ${goPath} && make test"
    }
}
