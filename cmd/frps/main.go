package main

import (
	_ "github.com/fatedier/frp/assets/frps"
	_ "github.com/fatedier/frp/pkg/metrics"
	"github.com/fatedier/frp/pkg/util/system"
)

func main() {
	system.EnableCompatibilityMode()
	Execute()
}
