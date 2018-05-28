
node {
    String applicationName = "drygopher"
    String goPath = "/go/src/github.com/eltorocorp/${applicationName}"

    docker.image("golang:1.10").inside("-v ${pwd()}:${goPath} -u root") {
        try {
            stage 'Pre-Build'
            sh "curl -X POST 'http://badges.awsp.eltoro.com?project=drygopher&item=build&value=pending&color=blue'"
            sh "cd ${goPath} && make prebuild"

            stage 'Build'
            sh "cd ${goPath} && make build"

            stage 'Test'
            sh "cd ${goPath} && make test"

            stage "Post-Build"
            sh "curl -X POST 'http://badges.awsp.eltoro.com?project=drygopher&item=build&value=passing&color=green'"
            currentBuild.result = 'SUCCESS'
        } catch (Exception err) {
            sh "curl -X POST 'http://badges.awsp.eltoro.com?project=drygopher&item=build&value=failing&color=red'"
            currentBuild.result = 'FAILURE'
        }
    }
}
