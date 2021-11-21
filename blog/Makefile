proto_gen:
	protoc --proto_path=proto proto/*.proto --go_out=gen/
	protoc --proto_path=proto proto/*.proto --go-grpc_out=gen/
#	protoc -I . --grpc-gateway_out ./gen/go \
#        --grpc-gateway_opt logtostderr=true \
#        --grpc-gateway_opt paths=source_relative \
#        --grpc-gateway_opt generate_unbound_methods=true \
#        proto/blog.proto
	protoc -I ./proto \
	  --go_out ./gen/proto --go_opt paths=source_relative \
	  --go-grpc_out ./gen/proto --go-grpc_opt paths=source_relative \
	  --grpc-gateway_out ./gen/proto --grpc-gateway_opt paths=source_relative \
	  ./proto/blog.proto