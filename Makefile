
make.PHONY:
proto:
	protoc greet/proto/greet.proto --go_out=plugins=grpc:.
	protoc calculator/calculatorpb/calculator.proto --go_out=plugins=grpc:.
	protoc blog/proto/blog.proto --go_out=plugins=grpc:.
clearGit:
	git branch --merged | egrep -v "(^\*|master|main|dev)" | xargs git branch -d