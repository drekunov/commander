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

	// navMu guards pending, the latest batch of navigation requests. Ordinary
	// navigation replaces the batch with one request; a refresh stores both
	// panels' requests so neither is dropped.
	navMu   sync.Mutex
	pending []NavRequest
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
	a.setPending([]NavRequest{{Panel: panel, Dir: dir}})
}

// Refresh stores a batch of navigation requests and wakes the run loop, so a
// caller can re-read several panels without one request replacing another.
func (a *App) Refresh(requests []NavRequest) {
	a.setPending(requests)
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
			if requests, ok := a.takePending(); ok {
				for _, nav := range requests {
					a.handleNav(nav)
				}
			}
		}
	}
}

// takePending returns and clears the latest batch, reporting whether one was
// waiting.
func (a *App) takePending() ([]NavRequest, bool) {
	a.navMu.Lock()
	defer a.navMu.Unlock()

	if len(a.pending) == 0 {
		return nil, false
	}

	requests := a.pending
	a.pending = nil

	return requests, true
}

// setPending replaces the pending batch and wakes the run loop. It is
// non-blocking and never discards the batch it was given.
func (a *App) setPending(requests []NavRequest) {
	a.navMu.Lock()
	a.pending = requests
	a.navMu.Unlock()

	select {
	case a.navSig <- struct{}{}:
	default:
	}
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
