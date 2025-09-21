FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o todo-app .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/todo-app .
COPY --from=builder /app/web ./web
RUN apk add --no-cache libc6-compat

EXPOSE 9092
ENV TODO_PORT=9092
ENV TODO_DBFILE=/data/tasks.db

VOLUME /data

CMD ["./todo-app"]
