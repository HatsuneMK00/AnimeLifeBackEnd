FROM golang:1.20 as builder
LABEL authors="makise"
LABEL stage=builder

WORKDIR /app/AnimeLifeBackEnd
COPY . ./

RUN export GO111MODULE=on && export GOPATH=$GOPATH:/app
RUN go env -w GOPROXY=https://goproxy.cn,direct && go mod download
RUN CGO_ENABLED=0 go build -tags=release -o backend main.go


FROM alpine:latest
LABEL authors="makise"

WORKDIR /app/AnimeLifeBackEnd
COPY --from=builder /app/AnimeLifeBackEnd/backend ./

EXPOSE 8080

ENTRYPOINT ["./backend"]