package main

import (
	app2 "2FA-PHP/app"
	"context"
	"embed"
	"fmt"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

func main() {
	// Initialize database
	db, err := app2.NewDatabase("two_fa.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
	} else {
		defer db.Close()
	}

	// Create a new context
	ctx := context.Background()
	// Create an instance of the app structure
	app := app2.NewApp(ctx, db)
	// Create application with options
	err = wails.Run(&options.App{
		Title:         "2FA-PHP",
		Width:         332,
		Height:        555,
		DisableResize: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Frameless:        false,
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		OnStartup:        app.Startup,
		Bind: []interface{}{
			app,
		},
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: true,
				HideTitle:                  false,
				HideTitleBar:               false,
				FullSizeContent:            false,
				UseToolbar:                 true,
				HideToolbarSeparator:       true,
			},
			WebviewIsTransparent: true,
			WindowIsTranslucent:  false,
			About: &mac.AboutInfo{
				Title:   "2FA PHP",
				Message: "© 2024 phuocph1903@gmail.com",
				Icon:    icon,
			},
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
