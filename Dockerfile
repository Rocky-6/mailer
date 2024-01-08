FROM golang:1.21-bullseye as builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app


FROM debian:bullseye-slim
RUN apt update && apt install ca-certificates -y

COPY --from=builder /usr/local/bin/app /app

CMD [ "/app" ]