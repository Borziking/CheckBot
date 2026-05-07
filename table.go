package main

import (
	"strings"

	"github.com/fogleman/gg"
)

func cellColor(value string) (r, g, b float64) {
	v := strings.TrimSpace(strings.ToUpper(value))

	switch {
	case v == "СС":
		return 1, 0.4, 0.4
	case v == "УД":
		return 0.4, 0.6, 1
	case v == "Б":
		return 1, 0.7, 0.3
	case strings.HasPrefix(v, "-"):
		return 1, 1, 0.5
	default:
		return 0.4, 0.8, 1
	}
}

func drawTable(table TableData) {
	cellW0 := 200.0
	cellW := 50.0
	cellH := 30.0
	cols := len(table.Headers)
	rows := len(table.Rows)

	width := cellW0 + cellW*float64(cols-1)
	height := cellH * float64(rows+1)

	dc := gg.NewContext(int(width), int(height))
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.LoadFontFace("Roboto-Regular.ttf", 11)

	colX := func(j int) float64 {
		if j == 0 {
			return 0
		}
		return cellW0 + cellW*float64(j-1)
	}

	colWidth := func(j int) float64 {
		if j == 0 {
			return cellW0
		}
		return cellW
	}

	for j, header := range table.Headers {
		x := colX(j)
		w := colWidth(j)

		dc.SetRGB(0.2, 0.6, 0.4)
		dc.DrawRectangle(x, 0, w, cellH)
		dc.Fill()

		dc.SetRGB(0, 0, 0)
		dc.SetLineWidth(0.5)
		dc.DrawRectangle(x, 0, w, cellH)
		dc.Stroke()

		dc.SetRGB(1, 1, 1)
		dc.DrawStringAnchored(header, x+w/2, cellH/2, 0.5, 0.5)
	}

	for i, row := range table.Rows {
		for j, cell := range row {
			x := colX(j)
			w := colWidth(j)
			y := float64(i+1) * cellH

			if j == 0 {
				dc.SetRGB(0.95, 0.95, 0.95)
			} else if strings.TrimSpace(cell) == "" {
				dc.SetRGB(0.7, 1, 0.7)
			} else {
				r, g, b := cellColor(cell)
				dc.SetRGB(r, g, b)
			}

			dc.DrawRectangle(x, y, w, cellH)
			dc.Fill()

			dc.SetRGB(0.7, 0.7, 0.7)
			dc.SetLineWidth(0.5)
			dc.DrawRectangle(x, y, w, cellH)
			dc.Stroke()

			dc.SetRGB(0, 0, 0)
			dc.DrawStringAnchored(cell, x+w/2, y+cellH/2, 0.5, 0.5)
		}
	}

	dc.SavePNG("table.png")
}
