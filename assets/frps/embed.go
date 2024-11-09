package frpc

import (
	"embed"

	"github.com/kirilngusi/go-reverse-proxy/assets"
)

//go:embed static/*
var content embed.FS

func init() {
	assets.Register(content)
}
