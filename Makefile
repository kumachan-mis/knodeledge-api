OPEN_API_DOCS_SERVER   := openapi-docs-server
OPEN_API_INDEX         := docs/openapi/_index.yaml
OPEN_API_BUNDLE        := tmp/openapi-bundle.yaml
OPEN_API_GENERATOR     := go-gin-server
OPEN_API_MODEL_DST     := interal/model
OPEN_API_MODEL_PACKAGE := model

start-openapi-docs-server:
	@docker run --detach --name ${OPEN_API_DOCS_SERVER} -v "${PWD}:/local" -p 8081:8081 \
		redocly/cli preview-docs \
		--api /local/${OPEN_API_INDEX} \
		--host 0.0.0.0 \
		--port 8081 \
		1> /dev/null

	@echo "OpenAPI docs server started at http://localhost:8081"

stop-openapi-docs-server:
	@docker stop ${OPEN_API_DOCS_SERVER} 1> /dev/null
	@docker rm ${OPEN_API_DOCS_SERVER} 1> /dev/null

	@echo "OpenAPI docs server stopped"

gen-openapi-models:
	rm -Rf ${OPEN_API_BUNDLE} ${OPEN_API_MODEL_DST}
	
	docker run --rm -v "${PWD}:/local" \
		redocly/cli bundle \
		/local/${OPEN_API_INDEX} \
		--output /local/${OPEN_API_BUNDLE}

	docker run --rm -v "${PWD}:/local" \
		openapitools/openapi-generator-cli generate \
		--input-spec /local/${OPEN_API_BUNDLE} \
		--generator-name ${OPEN_API_GENERATOR} \
		--global-property models,modelDocs=false \
		--additional-properties apiPath="",packageName=${OPEN_API_MODEL_PACKAGE} \
		--output /local/${OPEN_API_MODEL_DST}
