FROM golang:alpine as builder

ADD . /go/src/github.com/thiagotrennepohl/fortune-backend
ENV GO111MODULE on
WORKDIR /go/src/github.com/thiagotrennepohl/fortune-backend
RUN apk add --update git && go mod download && CGO_ENABLED=0 go build -a -installsuffix main.go -o main

# final stage
FROM alpine
WORKDIR /app
COPY --from=builder /go/src/github.com/thiagotrennepohl/fortune-backend/main /app/main
ENTRYPOINT ./main