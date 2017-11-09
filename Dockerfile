FROM golang:1.8.0-alpine

COPY . /go/src/github.com/timcurless/orchlift/

WORKDIR /go/src/github.com/timcurless/orchlift/

RUN apk --no-cache add git && \
    go get -u github.com/golang/dep/cmd/dep && \
    dep ensure && \
    GOOS=linux GOARCH=amd64 go build -o binaries/amd64/linux/orchlift-1.linux.amd64

ENTRYPOINT ["/usr/local/go/src/github.com/timcurless/orchlift/binaries/amd64/linux/orchlift-1.linux.amd64"]
