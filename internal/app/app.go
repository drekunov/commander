package app

import (
	"context"
	"log"
	"sync"
)

// NavRequest asks the app to read a directory and deliver it to a panel.
type NavRequest struct {
	Panel PanelID
	Dir   string
}

const rootDir = "/"

type App struct {
	ui        UI
	connector Connector

	// navMu guards pending, the single latest navigation request. A newer
	// request replaces an older one instead of being dropped.
	navMu   sync.Mutex
	pending *NavRequest
	navSig  chan struct{}
}

func New(ui UI, conn Connector) *App {
	return &App{
		ui:        ui,
		connector: conn,
		navSig:    make(chan struct{}, 1),
	}
}

// Navigate stores the latest navigation request and wakes the run loop. It is
// non-blocking and never discards the newest request.
func (a *App) Navigate(panel PanelID, dir string) {
	a.navMu.Lock()
	a.pending = &NavRequest{Panel: panel, Dir: dir}
	a.navMu.Unlock()

	select {
	case a.navSig <- struct{}{}:
	default:
	}
}

func (a *App) Run(ctx context.Context) error {
	log.Println(getBuildInfo())

	for _, panel := range []PanelID{PanelLeft, PanelRight} {
		a.loadPanel(ctx, panel, rootDir)
	}

	for {
		select {
		case <-ctx.Done():
			return nil

		case <-a.navSig:
			if nav, ok := a.takePending(); ok {
				a.handleNav(nav)
			}
		}
	}
}

// takePending returns and clears the latest request, reporting whether one was
// waiting.
func (a *App) takePending() (NavRequest, bool) {
	a.navMu.Lock()
	defer a.navMu.Unlock()

	if a.pending == nil {
		return NavRequest{}, false
	}

	nav := *a.pending
	a.pending = nil

	return nav, true
}

// loadPanel reads dir and delivers it to the panel. A read failure is shown
// through the info dialog and the application keeps running.
func (a *App) loadPanel(ctx context.Context, panel PanelID, dir string) {
	data, err := a.connector.ReadDir(dir)
	if err != nil {
		a.ui.Info(ctx, a.connector.Name(), dir, err.Error())

		return
	}

	a.ui.SetData(panel, dir, data)
}

// handleNav reads the requested directory and delivers it to the requesting
// panel. On failure the existing listing is untouched and the error is shown.
func (a *App) handleNav(nav NavRequest) {
	data, err := a.connector.ReadDir(nav.Dir)
	if err != nil {
		a.ui.Error(a.connector.Name(), err.Error())

		return
	}

	a.ui.SetData(nav.Panel, nav.Dir, data)
}
