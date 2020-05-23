FROM golang:1.14

WORKDIR /go/src/app
COPY discord.go .

RUN go get github.com/prometheus/client_golang/prometheus
RUN go get github.com/prometheus/client_golang/prometheus/promauto
RUN go get github.com/prometheus/client_golang/prometheus/promhttp

CMD ["go", "run", "discord.go"]

EXPOSE 8080
