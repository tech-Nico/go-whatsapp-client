package ui

import (
	"github.com/gdamore/tcell"
	log "github.com/sirupsen/logrus"
	"gitlab.com/tslocum/cview"
)

func (thisUI *UI) BuildLogView() *cview.TextView {
	logView := cview.NewTextView()
	logView.SetDynamicColors(true).
		SetScrollable(true).
		SetTitle("Logs").
		SetTitleColor(tcell.ColorSkyblue).
		SetBorder(true)

	log.SetOutput(logView)
	return logView
}
