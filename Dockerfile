FROM golang:1.22-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0
ARG VERSION=dev
RUN go build -trimpath -ldflags="-s -w -X main.version=${VERSION}" -o /service-cars ./cmd

FROM gcr.io/distroless/static:nonroot
WORKDIR /app

COPY --from=build /service-cars /app/service-cars

ENV APP_PORT=8080 \
    DB_HOST=db \
    DB_PORT=5432 \
    DB_USER=cars \
    DB_PASSWORD=cars \
    DB_NAME=cars \
    DB_SSLMODE=disable

EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/service-cars"]