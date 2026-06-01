package render

import "strings"

func Monitor(rows [][]string, out string) error {
	if len(rows) == 0 {
		return nil
	}

	headers := rows[0]
	cols := len(headers)

	data := [][]string{}
	for _, row := range rows[1:] {
		if len(row) > 0 && strings.TrimSpace(row[0]) != "" {
			data = append(data, row)
		}
	}

	nameW := 230.0
	cellW := 128.0
	cellH := 46.0
	cellGap := 6.0
	headerH := 64.0
	step := cellW + cellGap
	r := 9.0

	contentW := nameW + cellGap + float64(cols-1)*step
	contentH := headerH + cellGap + float64(len(data))*(cellH+cellGap)

	dc, ox, oy := newCard(contentW, contentH, 36, 56, "Мониторинг")

	colX := func(j int) float64 { return ox + nameW + cellGap + float64(j-1)*step }

	fillRoundedHex(dc, ox, oy, nameW, headerH, r, colHeader)
	textLeft(dc, fitString(dc, headers[0], fontBold, 16, nameW-32), ox, oy, nameW, headerH, 16, fontBold, 16, colOnDark)
	for j := 1; j < cols; j++ {
		x := colX(j)
		fillRoundedHex(dc, x, oy, cellW, headerH, r, colHeader)
		label := strings.ReplaceAll(headers[j], "\n", " ")
		useFont(dc, fontMedium, 13)
		setHex(dc, colOnDark)
		dc.DrawStringWrapped(label, x+cellW/2, oy+headerH/2, 0.5, 0.5, cellW-12, 1.2, 1)
	}

	rowTop := oy + headerH + cellGap
	for i, row := range data {
		y := rowTop + float64(i)*(cellH+cellGap)

		nameBg := colCard
		if i%2 == 1 {
			nameBg = colRowAlt
		}
		fillRoundedHex(dc, ox, y, nameW, cellH, r, nameBg)
		textLeft(dc, fitString(dc, row[0], fontMedium, 16, nameW-32), ox, y, nameW, cellH, 16, fontMedium, 16, colText)

		for j := 1; j < cols; j++ {
			x := colX(j)
			val := ""
			if j < len(row) {
				val = strings.ToUpper(strings.TrimSpace(row[j]))
			}

			switch val {
			case "TRUE":
				fillRoundedHex(dc, x, y, cellW, cellH, r, colAccent)
				drawCheck(dc, x+cellW/2, y+cellH/2, 18, colOnDark)
			case "FALSE", "":
				fillRoundedHex(dc, x, y, cellW, cellH, r, "#F1F5F9")
			default:
				fillRoundedHex(dc, x, y, cellW, cellH, r, "#E2E8F0")
				text(dc, val, x, y, cellW, cellH, fontMedium, 14, colMuted)
			}
		}
	}

	return dc.SavePNG(out)
}
