FROM golang

# Configure Go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin
RUN go get -u github.com/aidenMin/eip-controller

WORKDIR $GOPATH

ENTRYPOINT ["./bin/eip-controller"]
