FROM golang:1.12
LABEL maintainer="nitin.i.joy@gmail.com"
RUN mkdir /work

ENV GOPATH=/work/social
ENV GOBIN=/work/social/bin 

ADD ./social /work/social
ADD ./common /work/common

RUN apt-get update
RUN apt-get install -y libssl1.1
RUN wget https://github.com/neo4j-drivers/seabolt/releases/download/v1.7.4/seabolt-1.7.4-Linux-ubuntu-18.04.deb
RUN dpkg -i seabolt-1.7.4-Linux-ubuntu-18.04.deb
RUN rm seabolt-1.7.4-Linux-ubuntu-18.04.deb

WORKDIR /work/social
RUN go get
RUN go install
RUN ./bin/social

EXPOSE 8083/tcp