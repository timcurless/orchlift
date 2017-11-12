#!/usr/bin/env groovy

pipeline {
  agent {
    docker {
      image 'buildmystartup/custom-image-with-go'
      label 'docker'
    }
  }
  stages {
    stage('Create Binaries') {
      steps {
        sh "mkdir -p /gopath/src/github.com/timcurless/orchlift && \
            export GOPATH=/gopath && \
            cd /gopath/src/github.com/timcurless/orchlift"
        git "https://github.com/timcurless/orchlift.git"
        sh "go get -u github.com/golang/dep/cmd/dep && \
            /gopath/bin/dep ensure && \
            GOOS=linux GOARCH=amd64 go build -o binaries/amd64/${env.BUILD_NUMBER}/linux/orchlift-${env.BUILD_NUMBER}.linux.amd64"
      }
    }
  }
}
