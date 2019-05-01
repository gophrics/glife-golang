run:
	./location

build:
	protoc -I/usr/local/include -I. \
	-I/Users/nitinissacjoy/go/src \
	-I/Users/nitinissacjoy/go/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--go_out=plugins=grpc:. \
	proto/location.proto
	
	protoc -I/usr/local/include -I. \
	-I/Users/nitinissacjoy/go/src \
	-I/Users/nitinissacjoy/go/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--grpc-gateway_out=logtostderr=true:. \
	proto/location.proto

	go build -o ../build/location main.go
