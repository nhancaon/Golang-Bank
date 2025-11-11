# build stage
FROM golang:1.24-alpine3.22 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.22
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY db/migration ./db/migration
COPY start.sh .
RUN chmod +x /app/start.sh

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
