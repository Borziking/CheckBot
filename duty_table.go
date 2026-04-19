package main

import (
	"strings"
	"time"

	"github.com/fogleman/gg"
)

func drawDutyTable(rows [][]string) {

	colW := 120.0
	cellH := 28.0

	months := []string{}
	monthCols := [][2]int{}
	for i, val := range rows[1] {
		if strings.TrimSpace(val) != "" {
			months = append(months, val)
			monthCols = append(monthCols, [2]int{i, i + 1})
		}
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

	totalCols := len(months) * 2
	width := colW * float64(totalCols)
	height := cellH * float64(len(dataRows)+3)

	dc := gg.NewContext(int(width), int(height))
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.LoadFontFace("Roboto-Regular.ttf", 12)

	darkGreen := [3]float64{0.13, 0.37, 0.17}
	lightGreen := [3]float64{0.72, 0.93, 0.72}
	paleGreen := [3]float64{0.88, 0.97, 0.88}
	todayColor := [3]float64{0.7, 0.75, 0.71}
	grayColor := [3]float64{0.9, 0.9, 0.9}

	drawCell := func(x, y, w, h float64, bg [3]float64, text string, textLight bool) {
		dc.SetRGB(bg[0], bg[1], bg[2])
		dc.DrawRectangle(x, y, w, h)
		dc.Fill()

		dc.SetRGB(0.8, 0.8, 0.8)
		dc.SetLineWidth(0.5)
		dc.DrawRectangle(x, y, w, h)
		dc.Stroke()

		if textLight {
			dc.SetRGB(1, 1, 1)
		} else {
			dc.SetRGB(0, 0, 0)
		}
		dc.DrawStringAnchored(text, x+w/2, y+h/2, 0.5, 0.5)
	}

	drawCell(0, 0, width, cellH, darkGreen, "Месяц", true)

	for i, month := range months {
		x := float64(i*2) * colW
		drawCell(x, cellH, colW*2, cellH, darkGreen, month, true)
	}

	for i := range months {
		x := float64(i*2) * colW
		drawCell(x, cellH*2, colW, cellH, lightGreen, "Дата", false)
		drawCell(x+colW, cellH*2, colW, cellH, lightGreen, "Дежурный", false)
	}

	for rowIdx, row := range dataRows {
		if len(row) < 2 {
			continue
		}

		y := float64(rowIdx+3) * cellH

		for mIdx, cols := range monthCols {
			if cols[1] >= len(row) {
				continue
			}

			dateVal := strings.TrimSpace(row[cols[0]])
			nameVal := strings.TrimSpace(row[cols[1]])

			bg := paleGreen
			if rowIdx%2 == 0 {
				bg = lightGreen
			}

			monthNum := monthNames[months[mIdx]]
			cleanDate := strings.TrimRight(dateVal, "*")
			dayNum := 0
			for _, c := range cleanDate {
				if c >= '0' && c <= '9' {
					dayNum = dayNum*10 + int(c-'0')
				}
			}
			if dayNum == today && monthNum == currentMonth {
				bg = todayColor
			}

			if dateVal == "" || dateVal == "-" {
				bg = grayColor
			}

			x := float64(mIdx*2) * colW
			textLight := false
			drawCell(x, y, colW, cellH, bg, dateVal, textLight)
			drawCell(x+colW, y, colW, cellH, bg, nameVal, textLight)
		}
	}

	dc.SavePNG("duty.png")
}
