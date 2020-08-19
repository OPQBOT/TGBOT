module TGCLI

go 1.13

replace tdlib => ../tdlib

replace NetworkFramework => ../NetworkFramework

replace github.com/astaxie/beego => ../github.com/astaxie/beego

require (
	NetworkFramework v0.0.0-00010101000000-000000000000
	github.com/astaxie/beego v0.0.0-00010101000000-000000000000
	github.com/cjoudrey/gluahttp v0.0.0-20200626084403-ae897a63b78b
	github.com/junhsieh/goexamples v0.0.0-20190721045834-1c67ae74caa6 // indirect
	github.com/tengattack/gluasql v0.0.0-20181229041402-2e5ed630c4cf
	github.com/yuin/gopher-lua v0.0.0-20200807101526-d70801a73ebe
	layeh.com/gopher-luar v1.0.8
	tdlib v0.0.0-00010101000000-000000000000
)
