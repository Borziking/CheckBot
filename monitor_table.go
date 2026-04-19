package main

import (
	"strings"

	"github.com/fogleman/gg"
)

func drawMonitorTable(rows [][]string) {
	if len(rows) == 0 {
		return
	}

	nameColW := 150.0
	cellW := 110.0
	cellH := 40.0

	headers := rows[0]
	dataRows := rows[1:]

	filtered := [][]string{}
	for _, row := range dataRows {
		if len(row) > 0 && strings.TrimSpace(row[0]) != "" {
			filtered = append(filtered, row)
		}
	}

	cols := len(headers)
	width := nameColW + cellW*float64(cols-1)
	height := cellH * float64(len(filtered)+1)

	dc := gg.NewContext(int(width), int(height))
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.LoadFontFace("Roboto-Regular.ttf", 11)

	headerBg := [3]float64{0.13, 0.37, 0.17}
	checkedBg := [3]float64{0.72, 0.93, 0.72}
	emptyBg := [3]float64{0.98, 0.98, 0.88}
	nameBg := [3]float64{0.95, 0.95, 0.95}
	grayBg := [3]float64{0.85, 0.85, 0.85}

	drawCell := func(x, y, w, h float64, bg [3]float64, text string, light bool) {
		dc.SetRGB(bg[0], bg[1], bg[2])
		dc.DrawRectangle(x, y, w, h)
		dc.Fill()

		dc.SetRGB(0.75, 0.75, 0.75)
		dc.SetLineWidth(0.5)
		dc.DrawRectangle(x, y, w, h)
		dc.Stroke()

		if light {
			dc.SetRGB(1, 1, 1)
		} else {
			dc.SetRGB(0, 0, 0)
		}
		dc.DrawStringAnchored(text, x+w/2, y+h/2, 0.5, 0.5)
	}

	drawCell(0, 0, nameColW, cellH, headerBg, headers[0], true)
	for j := 1; j < cols; j++ {
		x := nameColW + cellW*float64(j-1)
		cleanHeader := strings.ReplaceAll(headers[j], "\n", "")
		drawCell(x, 0, cellW, cellH, headerBg, cleanHeader, true)
	}

	for i, row := range filtered {
		y := float64(i+1) * cellH

		drawCell(0, y, nameColW, cellH, nameBg, row[0], false)

		for j := 1; j < cols; j++ {
			x := nameColW + cellW*float64(j-1)
			var bg [3]float64
			var label string

			if j >= len(row) {
				bg = emptyBg
				label = ""
			} else {
				val := strings.ToUpper(strings.TrimSpace(row[j]))
				switch val {
				case "TRUE":
					bg = checkedBg
					label = "+"
				case "FALSE":
					bg = emptyBg
					label = ""
				default:
					bg = grayBg
					label = val
				}
			}
			drawCell(x, y, cellW, cellH, bg, label, false)
		}
	}

	dc.SavePNG("monitor.png")
}
