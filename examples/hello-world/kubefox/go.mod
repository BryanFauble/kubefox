module github.com/xigxog/kubefox/examples/hello-world/kubefox

go 1.21

replace github.com/xigxog/kubefox => ../../../

// TODO update when kubefox is released
require github.com/xigxog/kubefox v0.2.5-alpha.0.20231030185832-519fa63e00a6

require (
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/uuid v1.4.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231016165738-49dd2c1f3d0b // indirect
	google.golang.org/grpc v1.59.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
)
