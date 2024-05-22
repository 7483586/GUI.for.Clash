//go:build windows

package bridge

import (
	"embed"
	"log"
	"os"

	"github.com/energye/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var assets embed.FS

func CreateTray(a *App, icon []byte, fs embed.FS) {
	assets = fs
	systray.Run(func() {
		systray.SetIcon([]byte(icon))
		systray.SetTitle("GUI.for.Clash")
		systray.SetTooltip("GUI.for.Clash")
		systray.SetOnClick(func(menu systray.IMenu) {
			runtime.EventsEmit(a.Ctx, "onTrayClick")
			runtime.WindowShow(a.Ctx)
		})
		systray.SetOnRClick(func(menu systray.IMenu) {
			runtime.EventsEmit(a.Ctx, "onTrayRClick")
			menu.ShowMenu()
		})
	}, nil)
}

func (a *App) UpdateTray(tray TrayContent) {
	if tray.Icon != "" {
		ico, _ := assets.ReadFile("frontend/dist/" + tray.Icon)
		systray.SetIcon(ico)
	}
	if tray.Title != "" {
		systray.SetTitle(tray.Title)
	}
	if tray.Tooltip != "" {
		systray.SetTooltip(tray.Tooltip)
	}
}

func (a *App) UpdateTrayMenus(menus []MenuItem) {
	log.Printf("UpdateTrayMenus")

	systray.ResetMenu()

	for _, menu := range menus {
		createMenuItem(menu, a, nil)
	}
}

func createMenuItem(menu MenuItem, a *App, parent *systray.MenuItem) {
	if menu.Hidden {
		return
	}
	switch menu.Type {
	case "item":
		var m *systray.MenuItem
		if parent == nil {
			m = systray.AddMenuItem(menu.Text, menu.Tooltip)
		} else {
			m = parent.AddSubMenuItem(menu.Text, menu.Tooltip)
		}
		m.Click(func() { runtime.EventsEmit(a.Ctx, menu.Event) })

		if menu.Checked {
			m.Check()
		}

		for _, child := range menu.Children {
			createMenuItem(child, a, m)
		}
	case "separator":
		systray.AddSeparator()
	}
}

func (a *App) ExitApp() {
	systray.Quit()
	runtime.Quit(a.Ctx)
	os.Exit(0)
}
