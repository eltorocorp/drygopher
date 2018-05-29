
def setCoverageBadge(goPath) {
    def coverage = sh(script: "cd ${goPath} && cat coveragepct", returnStdout: true)
    def coverageUri = "\'http://badges.awsp.eltoro.com?project=drygopher&item=coverage&value=${coverage}&color=yellow\'"
    sh "curl -sX POST ${coverageUri}"
}

def setBuildStatusBadge(status, color) {
    def statusUri = "\'http://badges.awsp.eltoro.com?project=drygopher&item=build&value=${status}&color=${color}\'"
    sh "curl -sX POST ${statusUri}"
}

node {
    String goPath = "/go/src/github.com/eltorocorp/drygopher"
    docker.image("golang:1.10").inside("-v ${pwd()}:${goPath} -u root") {

        try {
            stage('Pre-Build') {
                setBuildStatusBadge('pending', 'blue')
                sh "chmod -R 0777 ${goPath}"
                checkout([$class: 'GitSCM', branches: [[name: '*/master']], doGenerateSubmoduleConfigurations: false, extensions: [], submoduleCfg: [], userRemoteConfigs: [[credentialsId: 'c01e62dd-e191-42c0-8b86-bbbef49c0292', url: 'git@github.com:eltorocorp/drygopher.git']]])
                sh "cd ${goPath} && make prebuild"
            }

            stage('Build') {
                sh "cd ${goPath} && make build"
            }

            stage('Test') {
                sh "cd ${goPath} && make test"
            }

            stage("Post-Build") {
                setBuildStatusBadge('passing', 'green')
                currentBuild.result = 'SUCCESS'
            }
        } catch (Exception err) {
            sh "echo ${err}"
            setBuildStatusBadge('failing', 'red')
            currentBuild.result = 'FAILURE'
        } finally {
            setCoverageBadge(goPath)
        }
    }
}
