FROM golang:1.25.3 AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o built-binary ./cmd/auth-service/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/built-binary .
# Copy any additional files needed for the application
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/docs/swagger.json ./docs/swagger.json

# Expose port if needed
EXPOSE ${PORT}

# Set environment variables if needed
# ENV VAR_NAME=value

CMD ["./built-binary"]
# Entry point if needed
# ENTRYPOINT ["./built-binary"]