package app

import (
	"context"
	"log"
)

// NavRequest asks the app to read a directory and deliver it to a panel.
type NavRequest struct {
	Panel PanelID
	Dir   string
}

const (
	navChannelSize = 8
	rootDir        = "/"
)

type App struct {
	ui        UI
	connector Connector

	navCh chan NavRequest
}

func New(ui UI, conn Connector) *App {
	return &App{
		ui:        ui,
		connector: conn,
		navCh:     make(chan NavRequest, navChannelSize),
	}
}

// Navigate enqueues a navigation request for the given panel. It is
// non-blocking so the UI event loop is never stalled by a full queue.
func (a *App) Navigate(panel PanelID, dir string) {
	select {
	case a.navCh <- NavRequest{Panel: panel, Dir: dir}:
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

		case nav := <-a.navCh:
			a.handleNav(nav)
		}
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
