package main

import (
	"fmt"
	"time"

	ui "github.com/dpetzold/termui"
	"k8s.io/apimachinery/pkg/util/duration"
)

func EventsPanel() *ui.Table {
	p := ui.NewTable()
	p.Height = 20
	p.BorderLabel = "Events"
	p.TextAlign = ui.AlignLeft
	p.Separator = false
	p.Headers = true
	p.Analysis()
	p.SetSize()
	return p
}

func updateEvents(eventsPanel *ui.Table) {
	eventRows := [][]string{
		[]string{"Last Seen", "Count", "Name", "Kind", "Type", "Reason", "Message"},
	}

	events, err := kubeClient.Events(Namespace)
	if err != nil {
		panic(err.Error())
	}

	for _, e := range events {
		eventRows = append(eventRows, []string{
			duration.ShortHumanDuration(time.Now().Sub(e.LastTimestamp.Time)),
			fmt.Sprintf("%d", e.Count),
			e.ObjectMeta.Name[0:20],
			e.InvolvedObject.Kind,
			e.Type,
			e.Reason,
			e.Message,
		})
	}

	if len(eventRows) > 18 {
		eventRows = eventRows[0:18]
	}

	eventsPanel.Rows = eventRows
	eventsPanel.Analysis()
	eventsPanel.SetSize()
}