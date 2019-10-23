FROM arm32v7/golang:alpine AS build-env
COPY . $GOPATH/src/github.com/techdad/infping
WORKDIR $GOPATH/src/github.com/techdad/infping
RUN apk --update add git
ARG GOOS=linux
ARG GOARCH=arm
ARG GOARM=7
RUN go get -v && go build -o /infping

# final stage
FROM arm32v7/alpine:latest
COPY --from=build-env /infping /
RUN apk add --no-cache ca-certificates fping
ENTRYPOINT ["/infping"]
