
make.PHONY:
proto:
	protoc greet/proto/greet.proto --go_out=plugins=grpc:.
	protoc calculator/calculatorpb/calculator.proto --go_out=plugins=grpc:.