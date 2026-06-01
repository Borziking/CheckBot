package main

import (
	"strings"
)

func statusColor(value string) (bg, fg string) {
	v := strings.TrimSpace(strings.ToUpper(value))
	switch {
	case v == "":
		return "#A7F3D0", "#065F46"
	case v == "СС":
		return "#FCA5A5", "#7F1D1D"
	case v == "УД":
		return "#93C5FD", "#1E3A8A"
	case v == "Б":
		return "#FDBA74", "#7C2D12"
	case strings.HasPrefix(v, "-"):
		return "#FCD34D", "#78350F"
	default:
		return "#67E8F9", "#155E75"
	}
}

type legendItem struct {
	bg, label string
}

func drawTable(table TableData) {
	cols := len(table.Headers)
	rowsN := len(table.Rows)
	if cols == 0 {
		return
	}

	nameW := 340.0
	cellW := 46.0
	cellH := 42.0
	cellGap := 5.0
	headerH := 48.0
	step := cellW + cellGap
	r := 8.0

	dayAreaW := float64(cols-1) * step
	contentW := nameW + cellGap + dayAreaW
	if contentW < nameW {
		contentW = nameW
	}

	legend := []legendItem{
		{"#A7F3D0", "Отработано"},
		{"#67E8F9", "Переработка"},
		{"#FCD34D", "Ушёл раньше"},
		{"#FCA5A5", "СС"},
		{"#93C5FD", "УД"},
		{"#FDBA74", "Больничный"},
	}
	legendH := 60.0

	contentH := headerH + cellGap + float64(rowsN)*(cellH+cellGap) + legendH

	dc, ox, oy := newCard(contentW, contentH, 36, 56, "График учёта времени")

	dayX := func(j int) float64 { return ox + nameW + cellGap + float64(j-1)*step }

	fillRoundedHex(dc, ox, oy, nameW, headerH, r, colHeader)
	if cols > 0 {
		textLeft(dc, table.Headers[0], ox, oy, nameW, headerH, 16, fontBold, 17, colOnDark)
	}
	for j := 1; j < cols; j++ {
		x := dayX(j)
		fillRoundedHex(dc, x, oy, cellW, headerH, r, colHeader)
		label := strings.ReplaceAll(table.Headers[j], "\n", "")
		text(dc, label, x, oy, cellW, headerH, fontMedium, 15, colOnDark)
	}

	rowTop := oy + headerH + cellGap
	for i, row := range table.Rows {
		y := rowTop + float64(i)*(cellH+cellGap)

		nameBg := colCard
		if i%2 == 1 {
			nameBg = colRowAlt
		}
		fillRoundedHex(dc, ox, y, nameW, cellH, r, nameBg)
		name := ""
		if len(row) > 0 {
			name = row[0]
		}
		textLeft(dc, fitString(dc, name, fontMedium, 16, nameW-32), ox, y, nameW, cellH, 16, fontMedium, 16, colText)

		for j := 1; j < cols; j++ {
			x := dayX(j)
			val := ""
			if j < len(row) {
				val = strings.TrimSpace(row[j])
			}
			bg, fg := statusColor(val)
			fillRoundedHex(dc, x, y, cellW, cellH, r, bg)
			if val != "" {
				text(dc, val, x, y, cellW, cellH, fontMedium, 15, fg)
			}
		}
	}

	ly := rowTop + float64(rowsN)*(cellH+cellGap) + 14
	lx := ox
	sw := 18.0
	for _, it := range legend {
		fillRoundedHex(dc, lx, ly, sw, sw, 5, it.bg)
		useFont(dc, fontRegular, 15)
		setHex(dc, colMuted)
		dc.DrawStringAnchored(it.label, lx+sw+8, ly+sw/2, 0, 0.5)
		w, _ := dc.MeasureString(it.label)
		lx += sw + 8 + w + 26
	}

	dc.SavePNG("table.png")
}
