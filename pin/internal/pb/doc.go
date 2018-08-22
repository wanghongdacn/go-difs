package pb

//go:generate protoc --gogo_out=. header.proto

// kludge to get vendoring right in protobuf output
//go:generate sed -i s,github.com/,github.com/Harold-the-Axeman/dacc-iam-filesystem/Godeps/_workspace/src/github.com/,g header.pb.go
