node {
    def root = tool name: 'Go1.8', type: 'go'
    ws("${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/src/github.com/eltorocorp/drygopher") {
        withEnv(["GOROOT=${root}", "GOPATH=${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/", "PATH+GO=${root}/bin"]) {
            env.PATH="${GOPATH}/bin:$PATH"
            
            stage 'Checkout'
        
            git url: 'https://github.com/eltorocorp/drygopher.git'
        
            stage 'preTest'
            sh 'go version'
            sh 'go get -u github.com/golang/dep/...'
            sh 'dep init'
            
            stage 'Test'
            sh 'go vet'
            sh 'go test -cover'
            
            stage 'Build'
            sh 'go build .'
        }
    }
}