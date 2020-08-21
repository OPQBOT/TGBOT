echo  `go version`
CGO_ENABLED=1;
export CGO_ENABLED
go build -ldflags "-s -w" -o ./TGCli ./
