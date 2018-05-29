
node {
    String goPath = "/go/src/github.com/eltorocorp/drygopher"
    docker.image("golang:1.10").inside("-v ${pwd()}:${goPath} -u root") {
        try {
            checkout([$class: 'GitSCM', branches: [[name: '*/master']], doGenerateSubmoduleConfigurations: false, extensions: [[$class: 'CleanBeforeCheckout']], submoduleCfg: [], userRemoteConfigs: [[credentialsId: 'c01e62dd-e191-42c0-8b86-bbbef49c0292', url: 'https://github.com/eltorocorp/drygopher.git']]])

            stage('Pre-Build') {
                sh "curl -sX POST 'http://badges.awsp.eltoro.com?project=drygopher&item=build&value=pending&color=blue'"
                sh "cd ${goPath} && make prebuild"
            }

            stage('Build') {
                sh "cd ${goPath} && make build"
            }

            stage('Test') {
                sh "cd ${goPath} && make test"
            }

            stage("Post-Build") {
                sh "cd ${goPath} && ls"
                def coverage = sh(script: "cd ${goPath} && cat coveragepct", returnStdout: true)
                def coverageUri = "\'http://badges.awsp.eltoro.com?project=drygopher&item=coverage&value=${coverage}&color=yellow\'"
                sh "echo ${coverageUri}" 
                sh "curl -sX POST ${coverageUri}"
                sh "curl -sX POST 'http://badges.awsp.eltoro.com?project=drygopher&item=build&value=passing&color=green'"
                currentBuild.result = 'SUCCESS'
            }
        } catch (Exception err) {
            sh "echo ${err}"
            sh "curl -sX POST 'http://badges.awsp.eltoro.com?project=drygopher&item=build&value=failing&color=red'"
            currentBuild.result = 'FAILURE'
        }

    }
}
