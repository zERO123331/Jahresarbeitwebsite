# syntax=docker/dockerfile:1

FROM golang:1.26
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ./cmd ./cmd/
COPY ./internal ./internal/
COPY ./ui ./ui/
COPY ./migrations ./migrations

RUN CGO_ENABLED=0 GOOS=linux go build -o ./application ./cmd/web


EXPOSE 4000

CMD ["./application"]

