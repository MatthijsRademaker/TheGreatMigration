FROM golang:1.26.2-alpine

WORKDIR /app

COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o backend .

EXPOSE 8080

CMD ["./backend"]
