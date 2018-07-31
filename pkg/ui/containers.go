package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/dpetzold/termui"
)

func NewContainersPanel(height int) *termui.Table {
	p := termui.NewTable()
	p.Height = height
	p.BorderLabel = "Containers"
	p.TextAlign = termui.AlignLeft
	p.Separator = false
	p.Headers = true
	p.Analysis()
	return p
}

func ContainerFooter() *termui.Par {
	text := strings.Join([]string{
		"p:Pod/Container",
		"c:Cpu",
		"C:CpuMax",
		"e:Memory",
		"E:MemoryMax",
		"s:Status",
		"a:Age",
		"b:back",
	}, " ")
	return Footer(text)
}

func ShowContainers() {
	ContainerPanel.Height = termui.TermHeight() - 1
	termui.Body.Rows = []*termui.Row{
		termui.NewRow(termui.NewCol(12, 0, ContainerPanel)),
		termui.NewRow(termui.NewCol(12, 0, ContainerFooter())),
	}
	ActiveWindow = ContainersWindow
}

func updateContainers(containersPanel *termui.Table) {

	containers := Containers

	sort.Slice(containers, func(i, j int) bool {
		return CompareStruct(containers, SortField, i, j, SortOrder)
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
			TimeToDurationStr(c.Age),
		})
	}

	max_rows := 1
	if ContainerPanel.Height > 4 {
		max_rows = ContainerPanel.Height - 3
	}

	if len(rows) > max_rows {
		rows = rows[0:max_rows]
	}

	containersPanel.Rows = rows
	containersPanel.Analysis()
}
