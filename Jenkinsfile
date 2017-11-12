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
        git "https://github.com/timcurless/orchlift.git"
        sh "ls -al && \
            apk --no-cache add git && \
            go get -u github.com/golang/dep/cmd/dep && \
            ls -al src/github.com src/github.com/timcurless && \
            cd ${env.GOPATH}/src/github.com/timcurless/orchlift && \
            ${env.GOPATH}/bin/dep ensure && \
            GOOS=linux GOARCH=amd64 go build -o binaries/amd64/${env.BUILD_NUMBER}/linux/orchlift-${env.BUILD_NUMBER}.linux.amd64"
      }
    }
  }
}
