#!/usr/bin/env groovy

node('docker') {
  String applicationName = "orchlift"
  String buildNumber = "0.1.${env.BUILD_NUMBER}"
  String goPath = "/go/src/github.com/timcurless/${applicationName}"

  stage('Checkout From GitHub') {
    checkout scm
  }

  stage('Create Binaries') {
    docker.image("golang:1.8.0-alpine").inside("-v ${pwd()}:${goPath}") {
      // Only going to build for Linux in this example
      sh "GOOS=linux GOARCH=amd64 go build -o binaries/amd64/${buildNumber}/linux/${applicationName}-${buildNumber}.linux.amd64"
    }
  }

  stage('Archive Artifacts') {
    archiveArtifacts artifacts: 'binaries/**', fingerprint: true
  }
}
