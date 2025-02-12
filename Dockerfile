FROM oven/bun:1.1.42 AS fe-builder

WORKDIR /app

RUN apt update && apt install make -y

COPY bun.lockb \
    package.json \
    tailwind.config.js \
    Makefile \
    ./
COPY src/frontend src/frontend

RUN bun install

RUN rm -rf src/frontend/dist/*
RUN make build-css
RUN make build-js

FROM golang:1.24-alpine AS be-builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

COPY --from=fe-builder /app/src/frontend/dist ./src/frontend/dist

# Build the Go app
RUN go build -o main src/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=be-builder /app/main ./main

EXPOSE 8080

CMD ["./main"]
