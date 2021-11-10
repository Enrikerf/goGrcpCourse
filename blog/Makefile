proto_gen:
	protoc --proto_path=proto proto/*.proto --go_out=gen/
	protoc --proto_path=proto proto/*.proto --go-grpc_out=gen/
	#protoc --proto_path=proto proto/*.proto --go_out=plugins=grpc:.