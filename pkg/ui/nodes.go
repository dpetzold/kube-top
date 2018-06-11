package ui

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/dpetzold/termui"
)

func NewNodePanel() *termui.List {
	p := termui.NewList()
	p.Border = false
	p.Height = NODE_PANEL_HEIGHT
	return p
}

func nodeRow() *termui.Row {
	return termui.NewRow(
		termui.NewCol(4, 0, NodePanel),
		termui.NewCol(4, 0, CpuColumn...),
		termui.NewCol(4, 0, MemoryColumn...),
	)
}

func NodesFooter() *termui.Par {
	return DashboardFooter()
}

func ShowNodeWindow() {
	NodePanel.Height = termui.TermHeight() - 1
	termui.Body.Rows = []*termui.Row{
		nodeRow(),
		termui.NewRow(
			termui.NewCol(12, 0, NodesFooter()),
		),
	}
	ActiveWindow = NodesWindow
}

func GaugePanel(label string, barColor termui.Attribute) *termui.Gauge {

	p := termui.NewGauge()
	p.BarColor = barColor
	p.BorderFg = termui.ColorWhite
	p.BorderLabelFg = termui.ColorCyan
	p.BorderLabel = label
	p.Height = 3
	p.LabelAlign = termui.AlignRight
	return p
}

func updateNodes(nodePanel *termui.List) error {

	var items []string

	resource_fmt := "%s[%s:](fg-cyan) [%s](fg-cyan,fg-bold)%s"
	gauge_fmt := "%d%% (%s)"

	CpuColumn = nil
	MemoryColumn = nil

	nodeResources := NodeResources

	sort.Slice(nodeResources, func(i, j int) bool {
		return CompareStruct(nodeResources, "CpuUsage", i, j, false)
	})

	columnMax := NODE_DISPLAY_COUNT
	if ActiveWindow == NodesWindow {
		columnMax = NODE_DISPLAY_MAX
	}

	for _, r := range nodeResources {

		cpuGauge := GaugePanel("Cpu", termui.ColorRed)
		memoryGauge := GaugePanel("Mem", termui.ColorCyan)

		resources := fmt.Sprintf(resource_fmt, "  ", "Cpu", r.CpuCapacity.String(), "  ")
		resources += fmt.Sprintf(resource_fmt, "", "Mem", r.MemoryCapacity.String(), "  ")
		resources += fmt.Sprintf(resource_fmt, "", "Pods", strconv.FormatInt(int64(r.Pods), 10), "")

		items = append(items, []string{
			fmt.Sprintf(" [%s](fg-white,fg-bold)", r.Name),
			resources,
			"",
		}...)

		if len(CpuColumn) < columnMax {
			CpuColumn = append(CpuColumn, cpuGauge)
			MemoryColumn = append(MemoryColumn, memoryGauge)
		}

		r.MemoryCapacity.Quantity.Sub(*r.MemoryUsage.Quantity)
		r.CpuCapacity.Quantity.Sub(*r.CpuUsage.Quantity)

		memoryGauge.Percent = r.PercentMemory
		memoryGauge.Label = fmt.Sprintf(gauge_fmt, r.PercentMemory, r.MemoryCapacity.String())

		cpuGauge.Percent = r.PercentCpu
		cpuGauge.Label = fmt.Sprintf(gauge_fmt, r.PercentCpu, r.CpuCapacity.String())
	}

	// XXX: Move this out
	if ActiveWindow == DashboardWindow || ActiveWindow == NodesWindow {
		termui.Body.Rows[0] = nodeRow()
	}

	max_rows := columnMax * NODE_DISPLAY_COUNT

	if len(items) > max_rows {
		items = items[0:max_rows]
	}

	nodePanel.Items = items

	return nil
}
