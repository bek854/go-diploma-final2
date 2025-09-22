FROM golang:1.21-alpine as builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -o todo-app .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/todo-app .
COPY --from=builder /app/web ./web
RUN apk add --no-cache libc6-compat

EXPOSE 8082
ENV TODO_PORT=8082
ENV TODO_DBFILE=/data/tasks.db

VOLUME /data

CMD ["./todo-app"]