FROM golang:1.16 as builder

RUN apt update

WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download

COPY ./src/* ./

ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64
RUN go build \
    -o shimoapp \
    -ldflags '-s -w'

FROM centos:7 as runner

WORKDIR /app

RUN yum install -y libreoffice libreoffice-langpack-ja
RUN yum install -y ipa-gothic-fonts ipa-pgothic-fonts

COPY --from=builder /go/src/shimoapp ./shimoapp
COPY ./files/* ./files/
COPY ./views/* ./views/

RUN if [ ! -d ./outfiles ]; then mkdir ./outfiles; fi

EXPOSE 8080

ENTRYPOINT ["./shimoapp"]
