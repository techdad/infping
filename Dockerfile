FROM golang:alpine AS build-env
COPY . /src
RUN cd /src && go build -o infping

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /src/infping /
ENTRYPOINT ["/infping"]
