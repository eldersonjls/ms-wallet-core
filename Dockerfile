FROM golang:1.20

WORKDIR /app/

RUN apt-get update && apt-get install -y librdkafka-dev

COPY . .

RUN go get -v github.com/rubenv/sql-migrate/...

CMD sql-migrate up && go run main.go