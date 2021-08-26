REPO?=/go/src/github.com/jusongchen/REST-app
VERSION_PACKAGE?=github.com/jusongchen/REST-app/pkg/exampleapp
APP?=demoapp
PORT?=12073
BIN_DIR?=./bin

RELEASE?=0.1.0
RELEASE_LATEST?=latest
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date --rfc-3339=seconds | sed 's/ /T/')
DOCKERUSER?=$(shell whoami)
DOCKERHUB?=registry.hub.docker.com
CONTAINER_IMAGE?=${DOCKERHUB}/${DOCKERUSER}/${APP}:${RELEASE}
CONTAINER_IMAGE_LATEST?=${DOCKERHUB}/${DOCKERUSER}/${APP}:${RELEASE_LATEST}
BUILD_CONTAINER?=${APP}-build
TEST_CONTAINER?=${APP}-test


GOBUILD_DOCKER_IMG?=golang:1.17


GOOS?=darwin
GOARCH?=amd64

clean:
	rm -f ${APP}
	go clean

build: clean
		CC=gcc CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build \
	    -mod vendor \
		-ldflags " -w -X ${VERSION_PACKAGE}.Release=${RELEASE} \
		-X ${VERSION_PACKAGE}.Commit=${COMMIT} -X ${VERSION_PACKAGE}.BuildTime=${BUILD_TIME}" \
		-o ${BIN_DIR}/${APP} \
		./cmd/swagger-example/

local-run: build
		ROOT_CA_PATH="pkg/common/falcon/testdata/rootCA.pem"					\
		CLIENT_CERT_PATH="pkg/common/falcon/testdata/client-certs/cert.pem"		\
		CLIENT_KEY_PATH="pkg/common/falcon/testdata/client-certs/key.pem"		\
		LOG_FORMAT="text" SWAGGER_UI_PATH="./dockerfiles/swaggerUI" TNS_ADMIN="./mtls_client" ${BIN_DIR}/${APP} serve 

build-mac: clean
		CC=gcc CGO_ENABLED=0 GOOS=darwin GOARCH=${GOARCH} go build \
	    -mod vendor \
		-ldflags " -w -X ${VERSION_PACKAGE}.Release=${RELEASE} \
		-X ${VERSION_PACKAGE}.Commit=${COMMIT} -X ${VERSION_PACKAGE}.BuildTime=${BUILD_TIME}" \
		-o ${BIN_DIR}/${APP} \
		./cmd/swagger-example/		

local-run-mac: build-mac
		LOG_FORMAT="text" SWAGGER_UI_PATH="./dockerfiles/swaggerUI" TNS_ADMIN="./mtls_client" ${BIN_DIR}/${APP} serve 

docker-build: clean

	docker build -t $(BUILD_CONTAINER):$(RELEASE) -f Dockerfile \
		--build-arg GOBUILD_DOCKER_IMG=$(GOBUILD_DOCKER_IMG) \
		--build-arg APP=$(APP) \
		--build-arg REPO=$(REPO)  \
		--build-arg VERSION_PACKAGE=$(VERSION_PACKAGE) \
		--build-arg RELEASE=$(RELEASE) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		.
	docker rm -f $(BUILD_CONTAINER) || true
	docker create --name $(BUILD_CONTAINER) $(BUILD_CONTAINER):$(RELEASE)
	mkdir -p ${BIN_DIR}/. && docker cp $(BUILD_CONTAINER):$(REPO)/$(APP) ${BIN_DIR}/.
	docker rm -f $(BUILD_CONTAINER)

	
container: 
	docker build -t $(APP):$(RELEASE) -f Dockerfile .

run-skip-build: 
	docker stop $(APP):$(RELEASE) || true && docker rm $(APP):$(RELEASE) || true
	docker run   -p ${PORT}:${PORT} --rm \
		-e "KUBECONFIG=/home/demoapp/kubeconfig" \
		-e PORT	\
		$(APP):$(RELEASE)

run: container run-skip-build
	

coverage:
	go test -race -coverprofile=cover.out ./... && go tool cover -html=cover.out

test:
	go test -race ./... 

test-v:
	go test -v -race ./...

push: container
	docker tag  $(APP):$(RELEASE)   $(CONTAINER_IMAGE)
	docker push $(CONTAINER_IMAGE)
	docker tag  $(APP):$(RELEASE)   $(CONTAINER_IMAGE_LATEST)
	docker push $(CONTAINER_IMAGE_LATEST)

deploy-local:
	docker ps --filter name='demoapp-dev' 
	docker stop demoapp-dev ||true
	docker run -d --name demoapp-dev -p ${PORT}:${PORT} --rm -e "PORT=${PORT}" $(APP):$(RELEASE)

