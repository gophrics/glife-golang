build:
	protoc -I/usr/local/include -I. \
	-I/Users/nitinissacjoy/go/src \
	-I/Users/nitinissacjoy/go/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--go_out=plugins=grpc:. \
	proto/profile/profile.proto
	
	protoc -I/usr/local/include -I. \
	-I/Users/nitinissacjoy/go/src \
	-I/Users/nitinissacjoy/go/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--grpc-gateway_out=logtostderr=true:. \
	proto/profile/profile.proto