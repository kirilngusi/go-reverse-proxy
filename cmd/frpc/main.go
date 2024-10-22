package main

import (
	_ "github.com/fatedier/frp/assets/frpc"
	"github.com/fatedier/frp/cmd/frpc/sub"
	"github.com/fatedier/frp/pkg/util/system"
)

func main() {
	system.EnableCompatibilityMode()
	sub.Execute()
}
