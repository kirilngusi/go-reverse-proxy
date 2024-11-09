package main

import (
	_ "github.com/kirilngusi/go-reverse-proxy/assets/frpc"
	"github.com/kirilngusi/go-reverse-proxy/cmd/frpc/sub"
	"github.com/kirilngusi/go-reverse-proxy/pkg/util/system"
)

func main() {
	system.EnableCompatibilityMode()
	sub.Execute()
}
