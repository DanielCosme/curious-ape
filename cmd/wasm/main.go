package main

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
	"image"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
)

func main() {
	go func() {
		window := new(app.Window)
		err := run(window)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window) error {
	theme := material.NewTheme()
	var ops op.Ops

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			title := material.H1(theme, "Hello, APE")
			maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
			title.Color = maroon
			title.Alignment = text.Middle
			title.Layout(gtx) // Draw the label to the graphics context.
			sp := SplitVisual{Ratio: -0.3}
			sp.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return FillWithLabel(gtx, theme, "Left", color.NRGBA{R: 255})
			}, func(gtx layout.Context) layout.Dimensions {
				return FillWithLabel(gtx, theme, "Right", color.NRGBA{B: 255})
			})
			e.Frame(gtx.Ops) // Pass the drawing operations to the GPU.
		}
	}
}

type SplitVisual struct {
	// Ratio keeps the current layout.
	// 0 is center, -1 completely to the left, 1 completely to the right.
	Ratio float32
}

func (s SplitVisual) Layout(gtx layout.Context, left, right layout.Widget) layout.Dimensions {
	proportion := (s.Ratio + 1) / 2
	leftSize := int(proportion * float32(gtx.Constraints.Max.X))

	rightOffset := leftSize
	rightSize := gtx.Constraints.Max.X - rightOffset
	{
		gtx := gtx
		gtx.Constraints = layout.Exact(image.Pt(leftSize, gtx.Constraints.Max.Y))
		left(gtx)
	}
	{
		trans := op.Offset(image.Pt(rightOffset, 0)).Push(gtx.Ops)
		gtx := gtx
		gtx.Constraints = layout.Exact(image.Pt(rightSize, gtx.Constraints.Max.Y))
		right(gtx)
		trans.Pop()
	}
	return layout.Dimensions{Size: gtx.Constraints.Max}
}

func FillWithLabel(gtx layout.Context, th *material.Theme, text string, backgroundColor color.NRGBA) layout.Dimensions {
	return layout.Center.Layout(gtx, material.H3(th, text).Layout)
}
