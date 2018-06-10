package ui

import (
	"strings"

	"github.com/dpetzold/kube-top/pkg/global"
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

	global.NodePanel.Height = global.NODE_PANEL_HEIGHT
	global.EventsPanel.Height = global.EVENTS_PANEL_HEIGHT
	global.ContainerPanel.Height = termui.TermHeight() - global.EVENTS_PANEL_HEIGHT - global.NODE_PANEL_HEIGHT - 1
	global.ActiveWindow = global.DashboardWindow

	termui.Body.Rows = []*termui.Row{
		nodeRow(),
		termui.NewRow(
			termui.NewCol(12, 0, global.ContainerPanel),
		),
		termui.NewRow(
			termui.NewCol(12, 0, global.EventsPanel),
		),
		termui.NewRow(
			termui.NewCol(12, 0, DashboardFooter()),
		),
	}
}
