package ui

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/dpetzold/kube-top/pkg/global"
	"github.com/dpetzold/termui"
)

func NewNodePanel() *termui.List {
	p := termui.NewList()
	p.Border = false
	p.Height = global.NODE_PANEL_HEIGHT
	return p
}

func nodeRow() *termui.Row {
	return termui.NewRow(
		termui.NewCol(4, 0, global.NodePanel),
		termui.NewCol(4, 0, global.CpuColumn...),
		termui.NewCol(4, 0, global.MemoryColumn...),
	)
}

func NodesFooter() *termui.Par {
	return DashboardFooter()
}

func ShowNodeWindow() {
	global.NodePanel.Height = termui.TermHeight() - 1
	termui.Body.Rows = []*termui.Row{
		nodeRow(),
		termui.NewRow(
			termui.NewCol(12, 0, NodesFooter()),
		),
	}
	global.ActiveWindow = global.NodesWindow
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

	global.CpuColumn = nil
	global.MemoryColumn = nil

	nodeResources := global.NodeResources

	sort.Slice(nodeResources, func(i, j int) bool {
		return CompareStruct(nodeResources, "CpuUsage", i, j, false)
	})

	columnMax := 3
	if global.ActiveWindow == global.NodesWindow {
		columnMax = 10
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

		if len(global.CpuColumn) < columnMax {
			global.CpuColumn = append(global.CpuColumn, cpuGauge)
			global.MemoryColumn = append(global.MemoryColumn, memoryGauge)
		}

		memoryGauge.Percent = r.PercentMemory
		memoryGauge.Label = fmt.Sprintf(gauge_fmt, r.PercentMemory, r.MemoryUsage.String())

		cpuGauge.Percent = r.PercentCpu
		cpuGauge.Label = fmt.Sprintf(gauge_fmt, r.PercentCpu, r.CpuUsage.String())
	}

	// XXX: Move this out
	if global.ActiveWindow == global.DashboardWindow || global.ActiveWindow == global.NodesWindow {
		termui.Body.Rows[0] = nodeRow()
	}

	max_rows := columnMax * 3

	// log.Infof("%v %v %v", len(nodeResources), len(items), max_rows)

	if len(items) > max_rows {
		items = items[0:max_rows]
	}

	nodePanel.Items = items

	return nil
}
