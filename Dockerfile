FROM golang

RUN mkdir /app && curl https://glide.sh/get | sh

COPY . /usr/local/go/src/github.com/timcurless/orchlift/

WORKDIR /usr/local/go/src/github.com/timcurless/orchlift/

RUN glide install && go build -o /usr/local/bin/orchlift

ENTRYPOINT ["orchlift"]
