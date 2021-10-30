module github.com/enrikerf/grcpGoProof/blog/server

go 1.17

replace proto v1.0.0 => ../proto

require (
	go.mongodb.org/mongo-driver v1.7.3 // indirect
	google.golang.org/grpc v1.41.0
	proto v1.0.0
)
