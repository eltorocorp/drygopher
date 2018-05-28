FROM golang:1.10

RUN go get github.com/vektra/mockery/.../
RUN go get github.com/golang/dep/cmd/dep
RUN mkdir -p /go/src/github.com/eltorocorp/drygopher
ADD . /go/src/github.com/eltorocorp/drygopher

ENTRYPOINT make