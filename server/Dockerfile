FROM golang:1.20.5

WORKDIR /app

COPY bin/dcache-server .
COPY example.json .

CMD ["./dcache-server"]
