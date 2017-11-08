FROM golang:1.8.0-alpine

COPY . /usr/local/go/src/github.com/timcurless/orchlift/

WORKDIR /usr/local/go/src/github.com/timcurless/orchlift/

RUN apk --no-cache add curl git && \
    curl https://glide.sh/get | sh && \
    glide install && \
    GOOS=linux GOARCH=amd64 go build -o binaries/amd64/linux/orchlift-1.linux.amd64

ENTRYPOINT ["/usr/local/go/src/github.com/timcurless/orchlift/binaries/amd64/linux/orchlift-1.linux.amd64"]
