package render

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/fogleman/gg"
)

var (
	fontDir     = resolveFontDir()
	fontRegular = filepath.Join(fontDir, "Roboto-Regular.ttf")
	fontMedium  = filepath.Join(fontDir, "Roboto-Medium.ttf")
	fontBold    = filepath.Join(fontDir, "Roboto-Bold.ttf")
)

func resolveFontDir() string {
	if d := os.Getenv("ASSETS_DIR"); d != "" {
		return filepath.Join(d, "fonts")
	}
	if exe, err := os.Executable(); err == nil {
		d := filepath.Join(filepath.Dir(exe), "assets", "fonts")
		if _, err := os.Stat(d); err == nil {
			return d
		}
	}
	return filepath.Join("assets", "fonts")
}

const (
	colBackdrop = "#ECFDF5"
	colCard     = "#FFFFFF"
	colHeader   = "#059669"
	colOnDark   = "#FFFFFF"
	colText     = "#1F2937"
	colMuted    = "#6B7280"
	colBorder   = "#D1FAE5"
	colRowAlt   = "#F0FDF4"
	colTile     = "#F1F5F9"
	colTileEmp  = "#E5E7EB"
	colAccent   = "#10B981"
	colAccentBg = "#D1FAE5"
	colAccentTx = "#047857"
)

func hexRGB(hex string) (float64, float64, float64) {
	var r, g, b int
	if len(hex) == 7 && hex[0] == '#' {
		fmt.Sscanf(hex[1:], "%02x%02x%02x", &r, &g, &b)
	}
	return float64(r) / 255, float64(g) / 255, float64(b) / 255
}

func hexRGBA(hex string, a float64) (float64, float64, float64, float64) {
	r, g, b := hexRGB(hex)
	return r, g, b, a
}

func setHex(dc *gg.Context, hex string) {
	dc.SetRGB(hexRGB(hex))
}

var fontWarnOnce sync.Once

func useFont(dc *gg.Context, path string, size float64) {
	if err := dc.LoadFontFace(path, size); err != nil {
		fontWarnOnce.Do(func() {
			log.Printf("render: cannot load font %q: %v — text will be blank; check assets/fonts or set ASSETS_DIR", path, err)
		})
		dc.LoadFontFace(fontRegular, size)
	}
}

func fillRoundedHex(dc *gg.Context, x, y, w, h, r float64, hex string) {
	setHex(dc, hex)
	dc.DrawRoundedRectangle(x, y, w, h, r)
	dc.Fill()
}

func text(dc *gg.Context, s string, x, y, w, h float64, font string, size float64, hex string) {
	useFont(dc, font, size)
	setHex(dc, hex)
	dc.DrawStringAnchored(s, x+w/2, y+h/2, 0.5, 0.5)
}

func textLeft(dc *gg.Context, s string, x, y, w, h, pad float64, font string, size float64, hex string) {
	useFont(dc, font, size)
	setHex(dc, hex)
	dc.DrawStringAnchored(s, x+pad, y+h/2, 0, 0.5)
}

func drawCheck(dc *gg.Context, cx, cy, s float64, hex string) {
	setHex(dc, hex)
	dc.SetLineWidth(s * 0.18)
	dc.SetLineCapRound()
	dc.MoveTo(cx-s*0.45, cy+s*0.02)
	dc.LineTo(cx-s*0.10, cy+s*0.38)
	dc.LineTo(cx+s*0.50, cy-s*0.40)
	dc.Stroke()
}

func fitString(dc *gg.Context, s, font string, size, maxW float64) string {
	useFont(dc, font, size)
	if w, _ := dc.MeasureString(s); w <= maxW {
		return s
	}
	r := []rune(s)
	for len(r) > 1 {
		r = r[:len(r)-1]
		cand := string(r) + "…"
		if w, _ := dc.MeasureString(cand); w <= maxW {
			return cand
		}
	}
	return string(r)
}

func newCard(contentW, contentH, pad, titleH float64, title string) (*gg.Context, float64, float64) {
	margin := 28.0
	W := contentW + pad*2 + margin*2
	H := contentH + pad*2 + margin*2 + titleH

	dc := gg.NewContext(int(W), int(H))
	setHex(dc, colBackdrop)
	dc.Clear()

	cardX, cardY := margin, margin
	cardW := W - margin*2
	cardH := H - margin*2

	dc.SetRGBA(hexRGBA(colText, 0.08))
	dc.DrawRoundedRectangle(cardX+3, cardY+6, cardW, cardH, 22)
	dc.Fill()

	fillRoundedHex(dc, cardX, cardY, cardW, cardH, 22, colCard)

	if title != "" {
		text(dc, title, cardX, cardY, cardW, titleH+pad, fontBold, 30, colText)
		setHex(dc, colBorder)
		dc.SetLineWidth(1)
		dc.DrawLine(cardX+pad, cardY+titleH+pad/2, cardX+cardW-pad, cardY+titleH+pad/2)
		dc.Stroke()
	}

	return dc, cardX + pad, cardY + titleH + pad
}
