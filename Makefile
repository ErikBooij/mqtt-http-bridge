.PHONY: build-css
build-css:
	bunx tailwindcss -i ./src/frontend/css/main.css -o ./src/frontend/dist/main.css

.PHONY: build-js
build-js:
	bun build \
			--target browser \
			--format esm \
			--sourcemap=external \
			--outdir src/frontend/dist/ \
			src/frontend/js/*.ts

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
