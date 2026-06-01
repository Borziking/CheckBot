package render

import (
	"strings"
	"time"
)

func Duty(rows [][]string, out string) error {
	if len(rows) < 4 {
		return nil
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
		return nil
	}

	monthNames := map[string]int{
		"Январь": 1, "Февраль": 2, "Март": 3,
		"Апрель": 4, "Май": 5, "Июнь": 6,
		"Июль": 7, "Август": 8, "Сентябрь": 9,
		"Октябрь": 10, "Ноябрь": 11, "Декабрь": 12,
	}

	dataRows := rows[3:]

	hasDate := func(row []string) bool {
		for _, c := range monthCols {
			if c[0] < len(row) {
				for _, ch := range row[c[0]] {
					if ch >= '0' && ch <= '9' {
						return true
					}
				}
			}
		}
		return false
	}
	last := -1
	for i, row := range dataRows {
		if hasDate(row) {
			last = i
		}
	}
	if last < 0 {
		return nil
	}
	dataRows = dataRows[:last+1]

	today := time.Now().Day()
	currentMonth := int(time.Now().Month())

	dateW := 66.0
	nameW := 200.0
	cellGap := 6.0
	monthW := dateW + cellGap + nameW
	monthGap := 18.0
	cellH := 44.0
	rowStep := cellH + cellGap
	headerH := 50.0
	subH := 38.0
	r := 9.0

	topBlock := headerH + cellGap + subH + cellGap
	contentW := monthW*float64(len(months)) + monthGap*float64(len(months)-1)
	contentH := topBlock + rowStep*float64(len(dataRows)) - cellGap

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
		mx := ox + float64(mIdx)*(monthW+monthGap)
		nameX := mx + dateW + cellGap
		monthNum := monthNames[month]

		fillRoundedHex(dc, mx, oy, monthW, headerH, r, colHeader)
		text(dc, month, mx, oy, monthW, headerH, fontBold, 21, colOnDark)

		sy := oy + headerH + cellGap
		fillRoundedHex(dc, mx, sy, dateW, subH, r, colAccentBg)
		text(dc, "Дата", mx, sy, dateW, subH, fontMedium, 15, colAccentTx)
		fillRoundedHex(dc, nameX, sy, nameW, subH, r, colAccentBg)
		text(dc, "Дежурный", nameX, sy, nameW, subH, fontMedium, 15, colAccentTx)

		for rIdx, row := range dataRows {
			y := oy + topBlock + float64(rIdx)*rowStep

			dateVal, nameVal := "", ""
			if monthCols[mIdx][0] < len(row) {
				dateVal = strings.TrimSpace(row[monthCols[mIdx][0]])
			}
			if monthCols[mIdx][1] < len(row) {
				nameVal = strings.TrimSpace(row[monthCols[mIdx][1]])
			}

			empty := dateVal == "" || dateVal == "-"
			isToday := !empty && dayFromString(dateVal) == today && monthNum == currentMonth

			tileBg := colTile
			dtTx, nmTx := colText, colText
			if empty {
				tileBg = colTileEmp
				dtTx, nmTx = colMuted, colMuted
			}
			if isToday {
				tileBg = colAccent
				dtTx, nmTx = colOnDark, colOnDark
			}

			dateText := dateVal
			if empty {
				dateText = "—"
			}
			dateFont := fontMedium
			if isToday {
				dateFont = fontBold
			}

			fillRoundedHex(dc, mx, y, dateW, cellH, r, tileBg)
			text(dc, dateText, mx, y, dateW, cellH, dateFont, 17, dtTx)

			fillRoundedHex(dc, nameX, y, nameW, cellH, r, tileBg)
			nmFont := fontRegular
			if isToday {
				nmFont = fontMedium
			}
			textLeft(dc, fitString(dc, nameVal, nmFont, 17, nameW-28), nameX, y, nameW, cellH, 14, nmFont, 17, nmTx)
		}
	}

	return dc.SavePNG(out)
}
