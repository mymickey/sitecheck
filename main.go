package main

import (
	"embed"
	_ "embed"
	"log"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var appIcon []byte

func init() {
	application.RegisterEvent[BenchmarkReport]("benchmark-finished")
	application.RegisterEvent[Settings]("settings-updated")
}

func main() {
	store := NewSettingsStore(DefaultSettingsPath())
	settings, err := store.Load()
	if err != nil {
		settings = DefaultSettings()
	}

	monitor := NewConnectivityMonitor(nil)
	telemetry := NewTelemetryClient(aptabaseAppKey)
	siteCheck := NewSiteCheckService(store, monitor, telemetry)

	app := application.New(application.Options{
		Name:        "SiteCheck",
		Description: "Menubar network connectivity monitor",
		Icon:        appIcon,
		Services: []application.Service{
			application.NewService(siteCheck),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ActivationPolicy: application.ActivationPolicyAccessory,
		},
	})

	menuController := NewMenuController(app, siteCheck, appIcon, settings)
	scheduler := NewBenchmarkScheduler(siteCheck)
	scheduler.Start(settings.IntervalMinutes)
	app.OnShutdown(scheduler.Stop)

	siteCheck.onSettingsSaved = func(settings Settings) {
		menuController.UpdateSettings(settings)
		scheduler.Start(settings.IntervalMinutes)
		app.Event.Emit("settings-updated", settings)
	}
	siteCheck.onBenchmarkFinish = func(report BenchmarkReport) {
		menuController.UpdateReport(report)
		app.Event.Emit("benchmark-finished", report)
	}
	siteCheck.onShowSettings = func() {
		menuController.ShowSettings()
	}
	siteCheck.onQuit = func() {
		app.Quit()
	}

	go func() {
		telemetry.Track("app_started", nil)
		time.Sleep(time.Second)
		_, _ = siteCheck.Benchmark()
	}()

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
