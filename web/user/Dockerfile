#FROM golang:1.15
#
#RUN go env -w GO111MODULE=on
#RUN go env -w GOPROXY=https://goproxy.cn,direct
#
#COPY . /go/src/go-project
#
#WORKDIR /go/src/go-project/service/user
#
#RUN go install
#
#ENTRYPOINT ["/go/bin/user"]
#多阶段编译
FROM golang:1.15-alpine AS builder

RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct

COPY . /go/src/go-project

WORKDIR /go/src/go-project/web/user

RUN go install

FROM alpine:3.13

COPY --from=builder /go/bin/user /bin/user-server

COPY --from=builder /go/src/go-project/web/user/config /bin/config

WORKDIR /bin

ENTRYPOINT ["./user-server"]