package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
)

// Constants
const htmlAbout = `Welcome on <b>Astilectron</b> demo!<br>
This is using the bootstrap and the bundler.`

// Vars injected via ldflags by bundler
var (
	AppName            = "Elementary"
	BuiltAt            string
	VersionAstilectron string
	VersionElectron    string
)

// Application Vars
var (
	debug         = flag.Bool("d", true, "enables the debug mode")
	activeWindow  *astilectron.Window
	mainWindow    *astilectron.Window
	startupWindow *astilectron.Window
	fileWindows   []*astilectron.Window
	app           *astilectron.Astilectron
	menu          *astilectron.Menu
	l             *log.Logger
)

func main() {
	// Parse flags
	flag.Parse()

	// Create logger
	l = log.New(log.Writer(), log.Prefix(), log.Flags())
	l.SetFlags(log.LstdFlags | log.Lshortfile)

	// Run bootstrap
	l.Printf("Running app built at %s\n", BuiltAt)
	if err := bootstrap.Run(options(l)); err != nil {
		l.Fatal(fmt.Errorf("running bootstrap failed: %w", err))
	}
}

func menuItems() []*astilectron.MenuItemOptions {
	return []*astilectron.MenuItemOptions{
		{
			Label: astikit.StrPtr("File"),
			SubMenu: []*astilectron.MenuItemOptions{
				{Role: astilectron.MenuItemRoleAbout, Label: astikit.StrPtr("About " + AppName)},
				{Type: astilectron.MenuItemTypeSeparator},
				{
					Label:       astikit.StrPtr("Preferences"),
					OnClick:     func(e astilectron.Event) bool { return false },
					Accelerator: astilectron.NewAccelerator("CommandOrControl", ","),
				},
				{Type: astilectron.MenuItemTypeSeparator},
				{Role: astilectron.MenuItemRoleServices},
				{Type: astilectron.MenuItemTypeSeparator},
				{Role: astilectron.MenuItemRoleHide, Label: astikit.StrPtr("Hide " + AppName)},
				{Role: astilectron.MenuItemRoleHideOthers},
				{Role: astilectron.MenuItemRoleUnhide},
				{Role: astilectron.MenuItemRoleReload},
				{Type: astilectron.MenuItemTypeSeparator},
				{Role: astilectron.MenuItemRoleQuit, Label: astikit.StrPtr("Quit " + AppName)},
				{
					Label: astikit.StrPtr("Debug Window"),
					OnClick: func(e astilectron.Event) (deleteListener bool) {
						startupWindow.OpenDevTools()
						startupWindow.Resize(1000, 200)
						for _, fileWindow := range fileWindows {
							fileWindow.OpenDevTools()
						}
						return false
					},
				},
			},
		},
		{
			Label: astikit.StrPtr("File"),
			SubMenu: []*astilectron.MenuItemOptions{
				{Label: astikit.StrPtr("Open forensicstore"), Accelerator: astilectron.NewAccelerator("CommandOrControl", "o"), OnClick: func(e astilectron.Event) (deleteListener bool) {
					bootstrap.SendMessage(mainWindow, "menu-open", nil)
					return false
				}},
				{Label: astikit.StrPtr("New forensicstore"), Accelerator: astilectron.NewAccelerator("CommandOrControl", "n"), OnClick: func(e astilectron.Event) (deleteListener bool) {
					bootstrap.SendMessage(mainWindow, "menu-new", nil)
					return false
				}},
				{Label: astikit.StrPtr("Import disk image"), Accelerator: astilectron.NewAccelerator("CommandOrControl", "i"), Enabled: astikit.BoolPtr(false)},
				{Type: astilectron.MenuItemTypeSeparator},
				{Label: astikit.StrPtr("Close"), Accelerator: astilectron.NewAccelerator("CommandOrControl", "w"), OnClick: func(e astilectron.Event) (deleteListener bool) {
					if activeWindow != nil {
						activeWindow.Close()
					}
					return false
				}},
			},
		},
		{
			// Label: astikit.StrPtr("Window"),
			Role: astilectron.MenuItemRoleWindowMenu,
			/*SubMenu: []*astilectron.MenuItemOptions{
				{Role: astilectron.MenuItemRoleMinimize},
				{Role: astilectron.MenuItemRoleZoom},
			},*/
		},
		{Role: astilectron.MenuItemRoleHelp},
	}
}

func options(l *log.Logger) bootstrap.Options {
	return bootstrap.Options{
		Asset:    Asset,
		AssetDir: AssetDir,
		AstilectronOptions: astilectron.Options{
			AppName:            "Elementary",
			AppIconDarwinPath:  "resources/icon.icns",
			AppIconDefaultPath: "resources/icon.png",
			SingleInstance:     true,
			VersionAstilectron: VersionAstilectron,
			VersionElectron:    VersionElectron,
		},
		Logger:      l,
		MenuOptions: menuItems(),
		OnWait: func(a *astilectron.Astilectron, ws []*astilectron.Window, m *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
			mainWindow = ws[0]
			mainWindow.On(astilectron.EventNameWindowEventFocus, func(e astilectron.Event) (deleteListener bool) {
				activeWindow = nil
				return false
			})
			startupWindow = ws[1]
			startupWindow.On(astilectron.EventNameWindowEventFocus, func(e astilectron.Event) (deleteListener bool) {
				activeWindow = nil
				return false
			})
			app = a
			menu = m
			return nil
		},
		RestoreAssets: RestoreAssets,
		Windows: []*bootstrap.Window{{
			Homepage:       "main.html",
			MessageHandler: open,
			Options:        &astilectron.WindowOptions{Show: astikit.BoolPtr(false)},
		}, {
			Homepage:       "open.html",
			MessageHandler: open,
			Options: &astilectron.WindowOptions{
				BackgroundColor: astikit.StrPtr("#333"),
				Center:          astikit.BoolPtr(true),
				Height:          astikit.IntPtr(170),
				Width:           astikit.IntPtr(480),
				TitleBarStyle:   astilectron.TitleBarStyleHiddenInset,
				Resizable:       astikit.BoolPtr(false),
				Minimizable:     astikit.BoolPtr(false),
			},
		}},
	}
}

var windowCount = 0

func storeWindow(store string) error {
	url := filepath.Join(app.Paths().DataDirectory(), "resources", "app", "index.html")
	o := &astilectron.WindowOptions{
		Title:         astikit.StrPtr(filepath.Base(store)),
		Center:        astikit.BoolPtr(true),
		Height:        astikit.IntPtr(800),
		Width:         astikit.IntPtr(1024),
		TitleBarStyle: astilectron.TitleBarStyleHiddenInset,
		WebPreferences: &astilectron.WebPreferences{
			Preload: astikit.StrPtr(filepath.Join(app.Paths().DataDirectory(), "resources", "app", "static", "js", "preload.js")),
		},
	}

	window, err := app.NewWindow(url, o)
	if err != nil {
		return err
	}
	window.OnMessage(handleMessages(window, handleStoreMessages(store), astikit.AdaptStdLogger(l)))
	window.On(astilectron.EventNameWindowEventClosed, func(e astilectron.Event) (deleteListener bool) {
		windowCount--
		return true
	})
	window.On(astilectron.EventNameWindowEventFocus, func(e astilectron.Event) (deleteListener bool) {
		activeWindow = window
		return false
	})

	if fileWindows == nil {
		fileWindows = []*astilectron.Window{}
	}
	fileWindows = append(fileWindows, window)

	if startupWindow.IsShown() {
		startupWindow.Hide()
	}
	windowCount++

	defer window.OpenDevTools()

	return window.Create()
}

type MessageIn struct {
	Name    string          `json:"name"`
	Payload json.RawMessage `json:"payload,omitempty"`
	Method  string          `json:"method,omitempty"`
}

type MessageHandler func(w *astilectron.Window, m MessageIn) (payload interface{}, err error)

func handleMessages(w *astilectron.Window, messageHandler MessageHandler, l astikit.SeverityLogger) astilectron.ListenerMessage {
	return func(m *astilectron.EventMessage) (v interface{}) {
		// Unmarshal message
		var i MessageIn
		var err error
		if err = m.Unmarshal(&i); err != nil {
			l.Error(fmt.Errorf("unmarshaling message %+v failed: %w", *m, err))
			return
		}

		// Handle message
		var p interface{}
		if p, err = messageHandler(w, i); err != nil {
			l.Error(fmt.Errorf("handling message %+v failed: %w", i, err))
		}

		// Return message
		if p != nil {
			o := &bootstrap.MessageOut{Name: i.Name + ".callback", Payload: p}
			if err != nil {
				o.Name = "error"
			}
			v = o
		}
		return
	}
}
