FROM golang:1.18

WORKDIR /examples/go-backpressure

COPY . .

EXPOSE 4000

CMD [ "go", "run", "./main.go", "./gauge.go" ]