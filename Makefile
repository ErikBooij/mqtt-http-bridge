DOCKER_IMAGE_NAME = erikbooij/mqtt-http-bridge
DOCKER_IMAGE_TAG = dev

CONFIG_FILE = config.dev.yaml

.PHONY: build-css
build-css:
	bunx tailwindcss -i ./src/frontend/css/main.css -o ./src/frontend/dist/main.css

.PHONY: build-image
build-image:
	docker build . -t $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) --progress=plain

.PHONY: build-js
build-js:
	bun build \
			--target browser \
			--format esm \
			--sourcemap=external \
			--outdir src/frontend/dist/ \
			src/frontend/js/*.ts

.PHONY: push-image
push-image:
	docker push $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

.PHONY: run-image
run-image:
	docker run --rm -e CONFIG_FILE=$(CONFIG_FILE) -v ./$(CONFIG_FILE):/app/$(CONFIG_FILE) -p 8081:8080 $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

.PHONY: watch-backend
watch-backend:
	air -c .air.toml

.PHONY: watch-css
watch-css:
	bunx tailwindcss -i ./src/frontend/css/main.css -o ./src/frontend/dist/main.css --watch=always

.PHONY: watch-js
watch-js:
	bun build \
    		--target browser \
    		--format esm \
    		--sourcemap=external \
    		--watch \
    		--outdir src/frontend/dist/ \
    		src/frontend/js/*.ts
