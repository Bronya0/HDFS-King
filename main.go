package main

import (
	"context"
	"embed"

	"hdfs-king/backend"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()
	connMgr := backend.NewConnectionManager()
	hdfsSvc := backend.NewHdfsService()

	err := wails.Run(&options.App{
		Title:  "HDFS King",
		Width:  1280,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			connMgr.Startup(ctx)
			hdfsSvc.Startup(ctx)
		},
		Bind: []interface{}{
			app,
			connMgr,
			hdfsSvc,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
