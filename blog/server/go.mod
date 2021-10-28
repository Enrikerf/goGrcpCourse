module github.com/enrikerf/grcpGoProof/blog/server

go 1.17

replace proto v1.0.0 => ../proto

require (
	google.golang.org/grpc v1.41.0
	proto v1.0.0
)

