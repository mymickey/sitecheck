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

//go:embed build/menubar.svg
var menubarIcon []byte

func init() {
	application.RegisterEvent[BenchmarkReport]("benchmark-finished")
	application.RegisterEvent[DNSTestReport]("dns-benchmark-finished")
	application.RegisterEvent[MyIPReport]("myip-finished")
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

	menuController := NewMenuController(app, siteCheck, menubarIcon, settings)
	scheduler := NewBenchmarkScheduler(siteCheck)
	scheduler.Start(settings.IntervalMinutes, settings.DNSIntervalHours)
	app.OnShutdown(scheduler.Stop)

	siteCheck.onSettingsSaved = func(settings Settings) {
		menuController.UpdateSettings(settings)
		scheduler.Start(settings.IntervalMinutes, settings.DNSIntervalHours)
		app.Event.Emit("settings-updated", settings)
	}
	siteCheck.onBenchmarkFinish = func(report BenchmarkReport) {
		menuController.UpdateReport(report)
		app.Event.Emit("benchmark-finished", report)
	}
	siteCheck.onDNSFinish = func(report DNSTestReport) {
		menuController.UpdateDNSReport(report)
		app.Event.Emit("dns-benchmark-finished", report)
	}
	siteCheck.onMyIPFinish = func(report MyIPReport) {
		menuController.UpdateMyIPReport(report)
		app.Event.Emit("myip-finished", report)
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
		go siteCheck.Benchmark(TriggerScheduler)
		go siteCheck.BenchmarkDNS(TriggerScheduler)
	}()

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
