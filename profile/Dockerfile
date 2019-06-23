FROM golang:1.12
LABEL maintainer="nitin.i.joy@gmail.com"
RUN mkdir /work

ENV GOPATH=/work/profile
ENV GOBIN=/work/profile/bin 

ADD ./profile /work/profile
ADD ./common /work/common

WORKDIR /work/profile
RUN go get
RUN go install
CMD ./run_background.sh

EXPOSE 8080/tcp