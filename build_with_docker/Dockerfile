FROM golang:1.5

ENV GOPATH /go
ENV PATH /go/bin:$PATH

# for glide
ENV GO15VENDOREXPERIMENT 1

# install glide and gox
RUN go get github.com/Masterminds/glide
RUN go get github.com/mitchellh/gox

WORKDIR /go/src/github.com/reiki4040
RUN git clone https://github.com/reiki4040/go-memo.git

WORKDIR go-memo
RUN git checkout -b feature/build_with_docer origin/feature/build_with_docer

