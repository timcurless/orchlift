#!/usr/bin/env groovy

node('docker') {
  String applicationName = "orchlift"
  String buildNumber = "0.1.${env.BUILD_NUMBER}"

  stage('Checkout From GitHub') {
    checkout scm
  }

  stage('Create Binaries') {
    docker.image("golang:1.8.0-alpine").inside("-v ${pwd()}") {
      sh "export GOPATH=$WORKSPACE && \
          apk --no-cache add git && \
          go get -u github.com/golang/dep/cmd/dep && \
          ls -al && \
          cd $GOPATH/src/github.com/timcurless/orchlift && \
          $GOPATH/bin/dep ensure && \
          GOOS=linux GOARCH=amd64 go build -o binaries/amd64/${buildNumber}/linux/${applicationName}-${buildNumber}.linux.amd64"
    }
  }

  stage('Archive Artifacts') {
    archiveArtifacts artifacts: 'binaries/**', fingerprint: true
  }
}
