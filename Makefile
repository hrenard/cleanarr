gen-radarr:
	docker run --user 1000 --rm -v ${PWD}:/local openapitools/openapi-generator-cli generate \
		-g go \
		--package-name radarr_client \
		--additional-properties=isGoSubmodule=true \
		-o /local/radarr-client \
		--skip-validate-spec \
		-i https://raw.githubusercontent.com/Radarr/Radarr/develop/src/Radarr.Api.V3/swagger.json