FROM golang:1.17 as builder

ARG VERSION_PACKAGE=github.com/jusongchen/REST-app/pkg/
ARG APP=demoapp
ARG REPO=/go/src/github.com/jusongchen/REST-app
ARG GOOS=linux
ARG GOARCH=amd64
ARG RELEASE="0.1.0"

ENV GO111MODULE=on 

WORKDIR $REPO 
COPY . ./

RUN cd $REPO && CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} RELEASE=${RELEASE} go build \
    -mod vendor \
    -ldflags "-s -w -X ${VERSION_PACKAGE}.Release=${RELEASE} \
    -X ${VERSION_PACKAGE}.Commit=$(git rev-parse --short HEAD) -X ${VERSION_PACKAGE}.BuildTime=$(date --rfc-3339=seconds | sed 's/ /T/')" \
    -installsuffix cgo \
    -o ${APP} \
    ./cmd/swagger-example/


# FROM alpine:latest
FROM ubuntu:latest

RUN useradd -r -u 7447 -g 0 demoapp \
    && mkdir -p /home/demoapp/bin 

ARG REPO=/go/src/github.com/jusongchen/REST-app
ENV APP=demoapp PORT=12073 SWAGGER_UI_PATH=/home/demoapp/swaggerUI 
EXPOSE $PORT

ENV PATH=$PATH:/home/demoapp/bin

ADD --chown=demoapp:root ./dockerfiles/swaggerUI $SWAGGER_UI_PATH
ADD --chown=demoapp:root migration/* /home/demoapp/

COPY --from=builder $REPO/$APP /home/demoapp/bin

USER demoapp
WORKDIR /home/demoapp
ENTRYPOINT /home/demoapp/bin/$APP serve
# CMD ["/bin/bash"]
