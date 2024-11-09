package main

import (
	_ "github.com/kirilngusi/go-reverse-proxy/assets/frps"
	_ "github.com/kirilngusi/go-reverse-proxy/pkg/metrics"
	"github.com/kirilngusi/go-reverse-proxy/pkg/util/system"
)

func main() {
	system.EnableCompatibilityMode()
	Execute()
}
