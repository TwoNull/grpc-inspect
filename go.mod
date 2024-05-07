module github.com/twonull/grpc-inspect

go 1.20

require (
	github.com/Philipp15b/go-steam/v3 v3.0.0
	github.com/davecgh/go-spew v1.1.1
	google.golang.org/grpc v1.63.2
	google.golang.org/protobuf v1.34.0
)

require (
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240429193739-8cf5692501f6 // indirect
)

replace github.com/Philipp15b/go-steam/v3 v3.0.0 => github.com/csfloat/go-steam/v3 v3.0.4
