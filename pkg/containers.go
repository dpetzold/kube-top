package main

import (
	"fmt"
	"sort"
	"strings"

	ui "github.com/dpetzold/termui"
)

func NewContainersPanel(height int) *ui.Table {
	p := ui.NewTable()
	p.Height = height
	p.BorderLabel = "Containers"
	p.TextAlign = ui.AlignLeft
	p.Separator = false
	p.Headers = true
	p.Analysis()
	return p
}

func ContainerFooter() *ui.Par {
	text := strings.Join([]string{
		"Sort:",
		"p:Pod/Container",
		"c:Cpu",
		"C:CpuMax",
		"e:Memory",
		"E:MemoryMax",
		"a:Age",
		"b:back",
	}, " ")
	return Footer(text)
}

func ShowContainers() {
	Globals.ContainerPanel.Height = ui.TermHeight() - 1
	ui.Body.Rows = []*ui.Row{
		ui.NewRow(ui.NewCol(12, 0, Globals.ContainerPanel)),
		ui.NewRow(ui.NewCol(12, 0, ContainerFooter())),
	}
	Globals.ActiveWindow = ContainersWindow
}

func updateContainers(containersPanel *ui.Table) {

	containers := Globals.Containers

	sort.Slice(containers, func(i, j int) bool {
		return cmp_struct(containers, Globals.SortField, i, j, Globals.SortOrder)
	})

	rows := [][]string{
		[]string{
			"Pod/Container",
			"Ready",
			"Status",
			"Restarts",
			"Cpu",
			"CpuMax",
			"Memory",
			"MemoryMax",
			"Age",
		},
	}

	for _, c := range containers {
		ready := "0"
		if c.Ready {
			ready = "1"
		}

		rows = append(rows, []string{
			c.Name,
			ready,
			c.Status,
			fmt.Sprintf("%d", c.Restarts),
			c.CpuUsage.String(),
			c.CpuMax.String(),
			c.MemoryUsage.String(),
			c.MemoryMax.String(),
			timeToDurationStr(c.Age),
		})
	}

	max_rows := Globals.ContainerPanel.Height - 3
	if len(rows) > max_rows {
		rows = rows[0:max_rows]
	}

	containersPanel.Rows = rows
	containersPanel.Analysis()
}
