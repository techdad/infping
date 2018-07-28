FROM golang:alpine AS build-env
COPY . $GOPATH/src/github.com/nickvanw/infping
WORKDIR $GOPATH/src/github.com/nickvanw/infping
RUN apk --update add git && go get -v && go build -o /infping

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /infping /
ENTRYPOINT ["/infping"]
