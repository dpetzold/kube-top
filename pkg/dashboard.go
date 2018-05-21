package main

import (
	"strings"

	ui "github.com/dpetzold/termui"
)

func DashboardFooter() *ui.Par {

	text := strings.Join([]string{
		"(D)ashboard",
		"(C)ontainers",
		"(E)vents",
		"(N)odes",
		"(Q)uit",
	}, "    ")

	return Footer(text)
}

func ShowDashboard() {

	Globals.NodePanel.Height = NODE_PANEL_HEIGHT
	Globals.EventsPanel.Height = EVENTS_PANEL_HEIGHT
	Globals.ContainerPanel.Height = ui.TermHeight() - EVENTS_PANEL_HEIGHT - NODE_PANEL_HEIGHT - 1
	Globals.ActiveWindow = DashboardWindow

	ui.Body.Rows = []*ui.Row{
		nodeRow(),
		ui.NewRow(
			ui.NewCol(12, 0, Globals.ContainerPanel),
		),
		ui.NewRow(
			ui.NewCol(12, 0, Globals.EventsPanel),
		),
		ui.NewRow(
			ui.NewCol(12, 0, DashboardFooter()),
		),
	}
}
