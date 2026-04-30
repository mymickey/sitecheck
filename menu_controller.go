package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

type MenuController struct {
	app            *application.App
	service        *SiteCheckService
	window         *application.WebviewWindow
	tray           *application.SystemTray
	logo           []byte
	settings       Settings
	report         BenchmarkReport
	icons          map[string][]byte
	isBenchmarking bool
	menuWindow     *application.WebviewWindow
	settingsWidth  int
	settingsHeight int
	mu             sync.Mutex
}

func NewMenuController(app *application.App, service *SiteCheckService, logo []byte, settings Settings) *MenuController {
	controller := &MenuController{
		app:      app,
		service:  service,
		logo:     logo,
		settings: settings,
		icons:    make(map[string][]byte),
		settingsWidth:  1024,
		settingsHeight: 768,
	}

	icon, _ := RenderTrayIcon(logo)
	controller.tray = app.SystemTray.New().SetTemplateIcon(icon)
	controller.tray.SetTooltip("SiteCheck")
	controller.tray.SetLabel("-- | --")

	// Create the popover window
	controller.menuWindow = app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:             "menu",
		Width:            300,
		Height:           260,
		URL:              "/?mode=menu",
		Frameless:        true,
		AlwaysOnTop:      true,
		Hidden:           true,
		BackgroundColour: application.NewRGBA(0, 0, 0, 0),
		Mac: application.MacWindow{
			Backdrop: application.MacBackdropTranslucent,
			TitleBar: application.MacTitleBarHidden,
		},
	})

	controller.tray.AttachWindow(controller.menuWindow)

	controller.menuWindow.OnWindowEvent(events.Mac.WindowDidBecomeKey, func(event *application.WindowEvent) {
		go controller.RunBenchmark()
	})

	controller.menuWindow.OnWindowEvent(events.Mac.WindowDidResignKey, func(event *application.WindowEvent) {
		controller.menuWindow.Hide()
	})

	go controller.refreshIcons(settings.Targets)
	return controller
}

func (c *MenuController) UpdateSettings(settings Settings) {
	c.mu.Lock()
	c.settings = settings
	c.mu.Unlock()
	c.rebuildMenu()
	go c.refreshIcons(settings.Targets)
}

func (c *MenuController) UpdateReport(report BenchmarkReport) {
	c.mu.Lock()
	c.report = report
	c.mu.Unlock()

	if icon, err := RenderTrayIcon(c.logo); err == nil {
		c.tray.SetTemplateIcon(icon)
	}

	fast, slow := "--", "--"
	if report.Summary.HasResults {
		fast = fmt.Sprintf("%dms", report.Summary.FastestMS)
		slow = fmt.Sprintf("%dms", report.Summary.SlowestMS)
	}
	c.tray.SetLabel(fmt.Sprintf("%s | %s", fast, slow))
}

func (c *MenuController) ShowSettings() {
	if c.menuWindow != nil {
		c.menuWindow.Hide()
	}
	if c.window != nil {
		c.window.Show().Focus()
		return
	}

	c.window = c.app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:             "settings",
		Title:            "SiteCheck Settings",
		Width:            c.settingsWidth,
		Height:           c.settingsHeight,
		MinWidth:         512,
		MinHeight:        384,
		MaxWidth:         1024,
		MaxHeight:        768,
		URL:              "/",
		BackgroundColour: application.NewRGB(250, 250, 250),
		Mac: application.MacWindow{
			Backdrop:                application.MacBackdropNormal,
			TitleBar:                application.MacTitleBarHiddenInset,
			InvisibleTitleBarHeight: 34,
		},
	})

	c.window.OnWindowEvent(events.Common.WindowDidResize, func(event *application.WindowEvent) {
		width, height := c.window.Size()
		c.settingsWidth = width
		c.settingsHeight = height
	})

	c.window.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
		width, height := c.window.Size()
		c.settingsWidth = width
		c.settingsHeight = height
		c.window = nil
	})
}

func (c *MenuController) RunBenchmark() {
	c.mu.Lock()
	if c.isBenchmarking {
		c.mu.Unlock()
		return
	}
	c.isBenchmarking = true
	c.mu.Unlock()

	_, _ = c.service.Benchmark()

	c.mu.Lock()
	c.isBenchmarking = false
	c.mu.Unlock()
}

func (c *MenuController) rebuildMenu() {
	// Native menu is no longer used. The tray window handles UI reactively.
}

func (c *MenuController) cachedIcon(iconURL string) []byte {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.icons[iconURL]
}

func (c *MenuController) refreshIcons(targets []Target) {
	client := &http.Client{Timeout: 2 * time.Second}
	updated := false

	for _, target := range targets {
		if target.IconURL == "" || len(c.cachedIcon(target.IconURL)) > 0 {
			continue
		}
		icon, err := fetchMenuIcon(client, target.IconURL)
		if err != nil || len(icon) == 0 {
			continue
		}

		c.mu.Lock()
		c.icons[target.IconURL] = icon
		c.mu.Unlock()
		updated = true
	}

	if updated {
		c.rebuildMenu()
	}
}

func fetchMenuIcon(client *http.Client, iconURL string) ([]byte, error) {
	response, err := client.Get(iconURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	contentType := response.Header.Get("Content-Type")
	if contentType != "" && !strings.HasPrefix(contentType, "image/") {
		return nil, nil
	}
	return io.ReadAll(io.LimitReader(response.Body, 64*1024))
}

func menuTargetLabel(target Target, result ProbeResult) string {
	if result.Status == StatusAvailable {
		return fmt.Sprintf("%s  -  %dms", target.Name, result.LatencyMS)
	}
	if result.Status == StatusUnavailable {
		return fmt.Sprintf("%s  -  Unavailable", target.Name)
	}
	return target.Name
}
