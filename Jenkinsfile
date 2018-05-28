#!/usr/bin/env groovy

// this will start an executor on a Jenkins agent with the docker label
node('docker') {
    String applicationName = "drygopher"
    String buildNumber = "${env.BUILD_NUMBER}"
    String goPath = "/go/src/github.com/eltorocorp/${applicationName}"

    stage('Checkout from GitHub') {
        checkout scm
    }

    stage("Build and Test") {
        docker.image("golang:1.10").inside("-v ${pwd()}:${goPath}") {
            sh "cd ${goPath} && make build"
            sh "cd ${goPath} && make test"
        }
    }
}
