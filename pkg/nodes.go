package main

import (
	"fmt"
	"sort"
	"strconv"

	ui "github.com/dpetzold/termui"
)

func NewNodePanel() *ui.List {
	p := ui.NewList()
	p.Border = false
	p.Height = NODE_PANEL_HEIGHT
	return p
}

func nodeRow() *ui.Row {
	return ui.NewRow(
		ui.NewCol(4, 0, Globals.NodePanel),
		ui.NewCol(4, 0, Globals.CpuColumn...),
		ui.NewCol(4, 0, Globals.MemoryColumn...),
	)
}

func NodesFooter() *ui.Par {
	return DashboardFooter()
}

func ShowNodeWindow() {
	Globals.NodePanel.Height = ui.TermHeight() - 1
	ui.Body.Rows = []*ui.Row{
		nodeRow(),
		ui.NewRow(
			ui.NewCol(12, 0, NodesFooter()),
		),
	}
	Globals.ActiveWindow = NodesWindow
}

func GaugePanel(label string, barColor ui.Attribute) *ui.Gauge {

	p := ui.NewGauge()
	p.BarColor = barColor
	p.BorderFg = ui.ColorWhite
	p.BorderLabelFg = ui.ColorCyan
	p.BorderLabel = label
	p.Height = 3
	p.LabelAlign = ui.AlignRight
	return p
}

func updateNodes(nodePanel *ui.List) error {

	var items []string

	resource_fmt := "%s[%s:](fg-cyan) [%s](fg-cyan,fg-bold)%s"
	gauge_fmt := "%d%% (%s)"

	Globals.CpuColumn = nil
	Globals.MemoryColumn = nil

	nodeResources := Globals.NodeResources

	sort.Slice(nodeResources, func(i, j int) bool {
		return cmp_struct(nodeResources, "CpuUsage", i, j, false)
	})

	columnMax := 3
	if Globals.ActiveWindow == NodesWindow {
		columnMax = 10
	}

	for _, r := range nodeResources {

		cpuGauge := GaugePanel("Cpu", ui.ColorRed)
		memoryGauge := GaugePanel("Mem", ui.ColorCyan)

		resources := fmt.Sprintf(resource_fmt, "  ", "Cpu", r.CpuCapacity.String(), "  ")
		resources += fmt.Sprintf(resource_fmt, "", "Mem", r.MemoryCapacity.String(), "  ")
		resources += fmt.Sprintf(resource_fmt, "", "Pods", strconv.FormatInt(int64(r.Pods), 10), "")

		items = append(items, []string{
			fmt.Sprintf(" [%s](fg-white,fg-bold)", r.Name),
			resources,
			"",
		}...)

		if len(Globals.CpuColumn) < columnMax {
			Globals.CpuColumn = append(Globals.CpuColumn, cpuGauge)
			Globals.MemoryColumn = append(Globals.MemoryColumn, memoryGauge)
		}

		memoryGauge.Percent = r.PercentMemory
		memoryGauge.Label = fmt.Sprintf(gauge_fmt, r.PercentMemory, r.MemoryUsage.String())

		cpuGauge.Percent = r.PercentCpu
		cpuGauge.Label = fmt.Sprintf(gauge_fmt, r.PercentCpu, r.CpuUsage.String())
	}

	// XXX: Move this out
	if Globals.ActiveWindow == DashboardWindow || Globals.ActiveWindow == NodesWindow {
		ui.Body.Rows[0] = nodeRow()
	}

	max_rows := columnMax * 3

	log.Infof("%v %v %v", len(nodeResources), len(items), max_rows)

	if len(items) > max_rows {
		items = items[0:max_rows]
	}

	nodePanel.Items = items

	return nil
}
