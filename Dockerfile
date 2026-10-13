# Build stage
FROM golang:1.27-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/grpc_server ./cmd/grpc_server

# Runtime stage
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /out/grpc_server /grpc_server

EXPOSE 55555

USER nonroot:nonroot

ENTRYPOINT ["/grpc_server"]