package ui

import (
	"strings"

	"github.com/dpetzold/termui"
)

func DashboardFooter() *termui.Par {

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

	NodePanel.Height = NODE_PANEL_HEIGHT
	EventsPanel.Height = EVENTS_PANEL_HEIGHT
	ContainerPanel.Height = termui.TermHeight() - EVENTS_PANEL_HEIGHT - NODE_PANEL_HEIGHT - 1
	ActiveWindow = DashboardWindow

	termui.Body.Rows = []*termui.Row{
		nodeRow(),
		termui.NewRow(
			termui.NewCol(12, 0, ContainerPanel),
		),
		termui.NewRow(
			termui.NewCol(12, 0, EventsPanel),
		),
		termui.NewRow(
			termui.NewCol(12, 0, DashboardFooter()),
		),
	}
}
