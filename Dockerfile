FROM golang:1.25-alpine

WORKDIR /app

COPY . .
RUN go mod tidy && go build -o notification-center .

EXPOSE 9999

CMD ["./notification-center"]
