package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// hudRow is one row of the overlay panel. An empty key puts the text in the
// first column (used by the headings and the spacer row).
type hudRow struct {
	key   string
	text  string
	color rl.Color
}

// hudRows builds the overlay for the current state.
func hudRows(sceneName string, bodies int, paused, vehicle bool) []hudRow {
	rows := []hudRow{
		{text: sceneName, color: rl.Black},
		{text: fmt.Sprintf("bodies %d    fps %d", bodies, rl.GetFPS()), color: rl.Gray},
	}
	if paused {
		rows = append(rows, hudRow{text: "PAUSED — Space resumes", color: rl.Red})
	}
	rows = append(rows, hudRow{}, // spacer
		hudRow{key: "N", text: "next scene", color: rl.DarkGray},
		hudRow{key: "R", text: "reset scene", color: rl.DarkGray},
		hudRow{key: "Space", text: "pause / resume", color: rl.DarkGray},
		hudRow{key: "Right", text: "step one frame", color: rl.DarkGray},
		hudRow{key: "F", text: "wireframe / solid", color: rl.DarkGray},
		hudRow{key: "P", text: "save a screenshot", color: rl.DarkGray},
		hudRow{key: "Esc", text: "quit", color: rl.DarkGray},
		hudRow{key: "Mouse", text: "drag to orbit, wheel to zoom", color: rl.Gray},
	)
	if vehicle {
		rows = append(rows,
			hudRow{key: "Arrows", text: "drive", color: rl.DarkGray},
			hudRow{key: "B", text: "brake", color: rl.DarkGray},
		)
	}
	return rows
}

// drawHUD draws the overlay panel. The layout is derived from the rows, so the
// panel always fits its content and rows cannot overlap — the original had a
// fixed-size panel and drew the orbit and vehicle hints at the same y.
func drawHUD(rows []hudRow) {
	const (
		fontSize = 20
		leading  = 28
		pad      = 14
		keyWidth = 130
		margin   = 12
	)

	width := int32(0)
	for _, r := range rows {
		w := rl.MeasureText(r.text, fontSize)
		if r.key != "" {
			w += keyWidth
		}
		if w > width {
			width = w
		}
	}
	width += 2 * pad
	height := 2*pad + leading*int32(len(rows))

	rl.DrawRectangle(margin, margin, width, height, rl.NewColor(255, 255, 255, 230))
	rl.DrawRectangleLines(margin, margin, width, height, rl.NewColor(0, 0, 0, 60))

	y := int32(margin + pad)
	for _, r := range rows {
		x := int32(margin + pad)
		if r.key != "" {
			rl.DrawText(r.key, x, y, fontSize, rl.Black)
			x += keyWidth
		}
		rl.DrawText(r.text, x, y, fontSize, r.color)
		y += leading
	}
}
