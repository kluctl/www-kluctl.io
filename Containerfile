FROM docker.io/golang

ENV CGO_ENABLED=1

RUN go install -tags extended github.com/gohugoio/hugo@v0.120.4

WORKDIR /src

ENTRYPOINT [ "hugo" ]
