API_REPOSITORY_ROOT := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
UI_REPOSITORY_ROOT  := $(shell dirname ${API_REPOSITORY_ROOT})/knodeledge-ui

OPEN_API_DOCS_SERVER := openapi-docs-server
OPEN_API_INDEX       := docs/openapi/_index.yaml
OPEN_API_BUNDLE      := tmp/openapi-bundle.yaml
OPEN_API_DOCS_DIST   := tmp/openapi-docs.html

OPEN_API_GO_GENERATOR := go-gin-server
OPEN_API_GO_PACKAGE   := model
OPEN_API_GO_DST       := interal/${OPEN_API_GO_PACKAGE}

OPEN_API_NODE_GENERATOR := typescript-fetch
OPEN_API_NODE_DST       := src/openapi

setup: dependencies generate
	cp .pre-commit .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit

dependencies:
	go get -v ./...
	go install go.uber.org/mock/mockgen@latest

run:
	go run cmd/app/main.go

format:
	go fmt ./...

lint:
	go vet ./...

build:
	go build -o tmp/app cmd/app/main.go

test:
	firebase emulators:exec 'go test ./...'

generate:
	go generate ./...

start-firestore-emulator:
	firebase emulators:start --only firestore

start-docs-server:
	@docker run --detach --name ${OPEN_API_DOCS_SERVER} -v "${API_REPOSITORY_ROOT}:/api" -p 8081:8081 \
		redocly/cli preview-docs \
		/api/${OPEN_API_INDEX} \
		--host 0.0.0.0 \
		--port 8081 \
		1> /dev/null

	@echo "OpenAPI docs server started at http://localhost:8081"

stop-docs-server:
	@docker stop ${OPEN_API_DOCS_SERVER} 1> /dev/null
	@docker rm ${OPEN_API_DOCS_SERVER} 1> /dev/null

	@echo "OpenAPI docs server stopped"

build-docs:
	docker run --rm -v "${API_REPOSITORY_ROOT}:/api" \
		redocly/cli build-docs \
		/api/${OPEN_API_INDEX} \
		--output /api/${OPEN_API_DOCS_DIST}

gen-openapi-go:
	rm -Rf \
		${API_REPOSITORY_ROOT}/${OPEN_API_BUNDLE} \
		${API_REPOSITORY_ROOT}/${OPEN_API_GO_DST}
	
	docker run --rm -v "${API_REPOSITORY_ROOT}:/api" \
		redocly/cli bundle \
		/api/${OPEN_API_INDEX} \
		--output /api/${OPEN_API_BUNDLE}

	docker run --rm -v "${API_REPOSITORY_ROOT}:/api" \
		openapitools/openapi-generator-cli generate \
		--input-spec /api/${OPEN_API_BUNDLE} \
		--generator-name ${OPEN_API_GO_GENERATOR} \
		--global-property models,modelDocs=false \
		--additional-properties apiPath="",packageName=${OPEN_API_GO_PACKAGE} \
		--output /api/${OPEN_API_GO_DST}

gen-openapi-node:
	rm -Rf \
		${API_REPOSITORY_ROOT}/${OPEN_API_BUNDLE} \
		${UI_REPOSITORY_ROOT}/${OPEN_API_NODE_DST}
	
	docker run --rm -v "${API_REPOSITORY_ROOT}:/api" \
		redocly/cli bundle \
		/api/${OPEN_API_INDEX} \
		--output /api/${OPEN_API_BUNDLE}

	docker run --rm -v "${API_REPOSITORY_ROOT}:/api" -v "${UI_REPOSITORY_ROOT}:/ui" \
		openapitools/openapi-generator-cli generate \
		--input-spec /api/${OPEN_API_BUNDLE} \
		--generator-name ${OPEN_API_NODE_GENERATOR} \
		--output /ui/${OPEN_API_NODE_DST}
	
	rm -Rf ${UI_REPOSITORY_ROOT}/${OPEN_API_NODE_DST}/.openapi-generator*
