package ui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	log "github.com/sirupsen/logrus"
)

func (thisUI *UI) BuildLogView() *tview.TextView {
	logView := tview.NewTextView()
	logView.SetDynamicColors(true).
		SetScrollable(true).
		SetTitle("Logs").
		SetTitleColor(tcell.ColorSkyblue).
		SetBorder(true)

	log.SetOutput(logView)
	return logView
}
