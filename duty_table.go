package main

import (
	"strings"
	"time"
)

func drawDutyTable(rows [][]string) {
	if len(rows) < 4 {
		return
	}

	months := []string{}
	monthCols := [][2]int{}
	for i, val := range rows[1] {
		if strings.TrimSpace(val) != "" {
			months = append(months, val)
			monthCols = append(monthCols, [2]int{i, i + 1})
		}
	}
	if len(months) == 0 {
		return
	}

	monthNames := map[string]int{
		"Январь": 1, "Февраль": 2, "Март": 3,
		"Апрель": 4, "Май": 5, "Июнь": 6,
		"Июль": 7, "Август": 8, "Сентябрь": 9,
		"Октябрь": 10, "Ноябрь": 11, "Декабрь": 12,
	}

	dataRows := rows[3:]
	today := time.Now().Day()
	currentMonth := int(time.Now().Month())

	dateW := 92.0
	nameW := 214.0
	monthW := dateW + nameW
	gap := 14.0
	rowH := 50.0
	monthH := 54.0
	subH := 42.0
	r := 12.0

	headerH := monthH + subH
	contentW := monthW*float64(len(months)) + gap*float64(len(months)-1)
	contentH := headerH + rowH*float64(len(dataRows))

	dc, ox, oy := newCard(contentW, contentH, 36, 56, "График дежурств")

	dayFromString := func(s string) int {
		clean := strings.TrimRight(strings.TrimSpace(s), "*")
		n := 0
		for _, c := range clean {
			if c >= '0' && c <= '9' {
				n = n*10 + int(c-'0')
			}
		}
		return n
	}

	for mIdx, month := range months {
		mx := ox + float64(mIdx)*(monthW+gap)
		monthNum := monthNames[month]

		fillTopRounded(dc, mx, oy, monthW, monthH, r, colHeader)
		text(dc, month, mx, oy, monthW, monthH, fontBold, 22, colOnDark)

		sy := oy + monthH
		setHex(dc, colAccentBg)
		dc.DrawRectangle(mx, sy, monthW, subH)
		dc.Fill()
		text(dc, "Дата", mx, sy, dateW, subH, fontMedium, 17, colAccentTx)
		text(dc, "Дежурный", mx+dateW, sy, nameW, subH, fontMedium, 17, colAccentTx)

		for rIdx, row := range dataRows {
			ry := oy + headerH + float64(rIdx)*rowH
			last := rIdx == len(dataRows)-1

			dateVal, nameVal := "", ""
			if monthCols[mIdx][0] < len(row) {
				dateVal = strings.TrimSpace(row[monthCols[mIdx][0]])
			}
			if monthCols[mIdx][1] < len(row) {
				nameVal = strings.TrimSpace(row[monthCols[mIdx][1]])
			}

			empty := dateVal == "" || dateVal == "-"
			isToday := !empty && dayFromString(dateVal) == today && monthNum == currentMonth

			bg := colCard
			if rIdx%2 == 1 {
				bg = colRowAlt
			}
			if isToday {
				bg = colAccentBg
			}
			if last {
				fillBottomRounded(dc, mx, ry, monthW, rowH, r, bg)
			} else {
				setHex(dc, bg)
				dc.DrawRectangle(mx, ry, monthW, rowH)
				dc.Fill()
			}

			if isToday {
				setHex(dc, colAccent)
				dc.DrawRectangle(mx, ry, 4, rowH)
				dc.Fill()
			}

			dtHex, nmHex := colText, colText
			if empty {
				dtHex, nmHex = colMuted, colMuted
			}
			if isToday {
				dtHex, nmHex = colAccentTx, colAccentTx
			}
			dateText := dateVal
			if empty {
				dateText = "—"
			}
			dateFont := fontRegular
			if isToday {
				dateFont = fontBold
			}
			text(dc, dateText, mx, ry, dateW, rowH, dateFont, 18, dtHex)
			textLeft(dc, nameVal, mx+dateW, ry, nameW, rowH, 14, fontRegular, 18, nmHex)
		}

		setHex(dc, colBorder)
		dc.SetLineWidth(1)
		dc.DrawLine(mx+dateW, oy+headerH, mx+dateW, oy+headerH+rowH*float64(len(dataRows)))
		dc.Stroke()
	}

	dc.SavePNG("duty.png")
}
