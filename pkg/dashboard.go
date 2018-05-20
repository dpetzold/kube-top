package main

import ui "github.com/dpetzold/termui"

func ShowDashboard() {

	NodePanel.Height = NODE_PANEL_HEIGHT
	EventsPanel.Height = EVENTS_PANEL_HEIGHT
	ContainerPanel.Height = ui.TermHeight() - EVENTS_PANEL_HEIGHT - NODE_PANEL_HEIGHT - 1

	ui.Body.Rows = []*ui.Row{
		nodeRow(),
		ui.NewRow(
			ui.NewCol(12, 0, ContainerPanel),
		),
		ui.NewRow(
			ui.NewCol(12, 0, EventsPanel),
		),
		ui.NewRow(
			ui.NewCol(12, 0, Footer()),
		),
	}
}
