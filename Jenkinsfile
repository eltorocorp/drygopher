
node {
    String applicationName = "drygopher"
    String goPath = "/go/src/github.com/eltorocorp/${applicationName}"

    docker.image("golang:1.10").inside("-v ${pwd()}:${goPath} -u root") {
        try {
            stage 'PreBuild'
            sh "cd ${goPath} && make prebuild"

            stage 'Build'
            sh "cd ${goPath} && make build"

            stage 'Test'
            sh "cd ${goPath} && make test"

            currentBuild.result = 'SUCCESS'
        } catch (Exception err) {
            currentBuild.result = 'FAILURE'
        }
        
        if (currentBuild.result == 'SUCCESS') {
            def badgeServerURL = 'badges.awsp.eltoro.com'
            def payload = 'project=drygopher&item=build_status&value=passing&color=green'
            sh "curl -X POST --data-urlencode \'payload=${payload}\' ${badgeServerURL}"
        } else {
            echo 'KAHN!'
        }
    }
}
