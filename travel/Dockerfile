FROM golang:1.12
LABEL maintainer="nitin.i.joy@gmail.com"
RUN mkdir /work

ENV GOPATH=/work/travel
ENV GOBIN=/work/travel/bin 

ADD ./travel /work/travel
ADD ./common /work/common

WORKDIR /work/travel
RUN go get
RUN go install
CMD ./bin/travel

EXPOSE 8082/tcp