# Build stage
FROM golang:1.23-alpine AS build
WORKDIR /src

COPY backend/go.mod backend/go.sum* ./
RUN go mod download

COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd/server

# Run stage
FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=build /out/server /server
EXPOSE 8080
ENTRYPOINT ["/server"]