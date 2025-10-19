// SPDX-License-Identifier: Unlicense OR MIT

package rendertest

import (
	"image"
	"image/color"
	"math"
	"testing"

	"gio.mleku.dev/app"
	"gio.mleku.dev/f32"
	"gio.mleku.dev/font/gofont"
	"gio.mleku.dev/gpu/headless"
	"gio.mleku.dev/op"
	"gio.mleku.dev/op/clip"
	"gio.mleku.dev/op/paint"
	"gio.mleku.dev/text"
	"golang.org/x/image/math/fixed"
)

// use some global variables for benchmarking so as to not pollute
// the reported allocs with allocations that we do not want to count.
var (
	c1, c2, c3    = make(chan op.CallOp), make(chan op.CallOp), make(chan op.CallOp)
	op1, op2, op3 op.Ops
)

func setupBenchmark(b *testing.B) (app.Context, *headless.Window, *text.Shaper) {
	sz := image.Point{X: 1024, Y: 1200}
	w := newWindow(b, sz.X, sz.Y)
	ops := new(op.Ops)
	gtx := app.Context{
		Ops:  ops,
		Size: sz,
	}
	shaper := text.NewShaper(text.WithCollection(gofont.Collection()))
	return gtx, w, shaper
}

func resetOps(gtx app.Context) {
	gtx.Ops.Reset()
	op1.Reset()
	op2.Reset()
	op3.Reset()
}

func finishBenchmark(b *testing.B, w *headless.Window) {
	b.StopTimer()
	if *dumpImages {
		img := image.NewRGBA(image.Rectangle{Max: w.Size()})
		err := w.Screenshot(img)
		w.Release()
		if err != nil {
			b.Error(err)
		}
		saveImage(b, b.Name()+".png", img)
	}
}

func BenchmarkDrawUICached(b *testing.B) {
	// As BenchmarkDraw but the same op.Ops every time that is not reset - this
	// should thus allow for maximal cache usage.
	gtx, w, th := setupBenchmark(b)
	defer w.Release()
	drawCore(gtx, th)
	w.Frame(gtx.Ops)

	for b.Loop() {
		w.Frame(gtx.Ops)
	}
	finishBenchmark(b, w)
}

func BenchmarkDrawUI(b *testing.B) {
	// BenchmarkDraw is intended as a reasonable overall benchmark for
	// the drawing performance of the full drawing pipeline, in each iteration
	// resetting the ops and drawing, similar to how a typical UI would function.
	// This will allow font caching across frames.
	gtx, w, th := setupBenchmark(b)
	defer w.Release()
	drawCore(gtx, th)
	w.Frame(gtx.Ops)
	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		resetOps(gtx)

		off := float32(math.Mod(float64(i)/10, 10))
		t := op.Affine(f32.AffineId().Offset(f32.Pt(off, off))).Push(gtx.Ops)

		drawCore(gtx, th)

		t.Pop()
		w.Frame(gtx.Ops)
	}
	finishBenchmark(b, w)
}

func BenchmarkDrawUITransformed(b *testing.B) {
	// Like BenchmarkDraw UI but transformed at every frame
	gtx, w, th := setupBenchmark(b)
	defer w.Release()
	drawCore(gtx, th)
	w.Frame(gtx.Ops)
	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		resetOps(gtx)

		angle := float32(math.Mod(float64(i)/1000, 0.05))
		a := f32.AffineId().Shear(f32.Point{}, angle, angle).Rotate(f32.Point{}, angle)
		t := op.Affine(a).Push(gtx.Ops)

		drawCore(gtx, th)

		t.Pop()
		w.Frame(gtx.Ops)
	}
	finishBenchmark(b, w)
}

func Benchmark1000Circles(b *testing.B) {
	// Benchmark1000Shapes draws 1000 individual shapes such that no caching between
	// shapes will be possible and resets buffers on each operation to prevent caching
	// between frames.
	gtx, w, _ := setupBenchmark(b)
	defer w.Release()
	draw1000Circles(gtx)
	w.Frame(gtx.Ops)
	b.ReportAllocs()

	for b.Loop() {
		resetOps(gtx)
		draw1000Circles(gtx)
		w.Frame(gtx.Ops)
	}
	finishBenchmark(b, w)
}

func Benchmark1000CirclesInstanced(b *testing.B) {
	// Like Benchmark1000Circles but will record them and thus allow for caching between
	// them.
	gtx, w, _ := setupBenchmark(b)
	defer w.Release()
	draw1000CirclesInstanced(gtx)
	w.Frame(gtx.Ops)
	b.ReportAllocs()

	for b.Loop() {
		resetOps(gtx)
		draw1000CirclesInstanced(gtx)
		w.Frame(gtx.Ops)
	}
	finishBenchmark(b, w)
}

func draw1000Circles(gtx app.Context) {
	ops := gtx.Ops
	for x := range 100 {
		op.Offset(image.Pt(x*10, 0)).Add(ops)
		for y := range 10 {
			paint.FillShape(ops,
				color.NRGBA{R: 100 + uint8(x), G: 100 + uint8(y), B: 100, A: 120},
				clip.RRect{Rect: image.Rect(0, 0, 10, 10), NE: 5, SE: 5, SW: 5, NW: 5}.Op(ops),
			)
			op.Offset(image.Pt(0, 100)).Add(ops)
		}
	}
}

func draw1000CirclesInstanced(gtx app.Context) {
	ops := gtx.Ops

	r := op.Record(ops)
	cl := clip.RRect{Rect: image.Rect(0, 0, 10, 10), NE: 5, SE: 5, SW: 5, NW: 5}.Push(ops)
	paint.PaintOp{}.Add(ops)
	cl.Pop()
	c := r.Stop()

	for x := range 100 {
		op.Offset(image.Pt(x*10, 0)).Add(ops)
		for y := range 10 {
			paint.ColorOp{Color: color.NRGBA{R: 100 + uint8(x), G: 100 + uint8(y), B: 100, A: 120}}.Add(ops)
			c.Add(ops)
			op.Offset(image.Pt(0, 100)).Add(ops)
		}
	}
}

func drawCore(gtx app.Context, shaper *text.Shaper) {
	c1 := drawIndividualShapes(gtx, shaper)
	c2 := drawShapeInstances(gtx, shaper)
	c3 := drawText(gtx, shaper)

	(<-c1).Add(gtx.Ops)
	(<-c2).Add(gtx.Ops)
	(<-c3).Add(gtx.Ops)
}

func drawIndividualShapes(gtx app.Context, shaper *text.Shaper) chan op.CallOp {
	// draw 81 rounded rectangles of different solid colors - each one individually
	go func() {
		ops := &op1
		c := op.Record(ops)
		for x := range 9 {
			op.Offset(image.Pt(x*50, 0)).Add(ops)
			for y := range 9 {
				paint.FillShape(ops,
					color.NRGBA{R: 100 + uint8(x), G: 100 + uint8(y), B: 100, A: 120},
					clip.RRect{Rect: image.Rect(0, 0, 25, 25), NE: 10, SE: 10, SW: 10, NW: 10}.Op(ops),
				)
				op.Offset(image.Pt(0, 50)).Add(ops)
			}
		}
		c1 <- c.Stop()
	}()
	return c1
}

func drawShapeInstances(gtx app.Context, shaper *text.Shaper) chan op.CallOp {
	// draw 400 textured circle instances, each with individual transform
	go func() {
		ops := &op2
		co := op.Record(ops)

		r := op.Record(ops)
		cl := clip.RRect{Rect: image.Rect(0, 0, 25, 25), NE: 10, SE: 10, SW: 10, NW: 10}.Push(ops)
		paint.PaintOp{}.Add(ops)
		cl.Pop()
		c := r.Stop()

		squares.Add(ops)
		rad := float32(0)
		for x := range 20 {
			for y := range 20 {
				t := op.Offset(image.Pt(x*50+25, y*50+25)).Push(ops)
				c.Add(ops)
				t.Pop()
				rad += math.Pi * 2 / 400
			}
		}
		c2 <- co.Stop()
	}()
	return c2
}

func drawText(gtx app.Context, shaper *text.Shaper) chan op.CallOp {
	// draw 40 lines of text with different transforms.
	go func() {
		ops := &op3
		c := op.Record(ops)

		for x := range 40 {
			params := text.Parameters{
				Font:     gofont.Collection()[0].Font,
				PxPerEm:  fixed.I(16),
				MaxWidth: 1000,
			}
			shaper.LayoutString(params, textRows[x])

			t := op.Offset(image.Pt(0, 24*x)).Push(ops)
			// Draw the shaped text
			for {
				glyph, ok := shaper.NextGlyph()
				if !ok {
					break
				}
				// Simple text rendering - just draw a rectangle for each glyph
				paint.FillShape(ops, color.NRGBA{A: 255},
					clip.Rect(image.Rect(0, 0, int(glyph.Advance>>6), 16)).Op())
				op.Offset(image.Pt(int(glyph.Advance>>6), 0)).Add(ops)
			}
			t.Pop()
		}
		c3 <- c.Stop()
	}()
	return c3
}

var textRows = []string{
	"1. I learned from my grandfather, Verus, to use good manners, and to",
	"put restraint on anger. 2. In the famous memory of my father I had a",
	"pattern of modesty and manliness. 3. Of my mother I learned to be",
	"pious and generous; to keep myself not only from evil deeds, but even",
	"from evil thoughts; and to live with a simplicity which is far from",
	"customary among the rich. 4. I owe it to my great-grandfather that I",
	"did not attend public lectures and discussions, but had good and able",
	"teachers at home; and I owe him also the knowledge that for things of",
	"this nature a man should count no expense too great.",
	"5. My tutor taught me not to favour either green or blue at the",
	"chariot races, nor, in the contests of gladiators, to be a supporter",
	"either of light or heavy armed. He taught me also to endure labour;",
	"not to need many things; to serve myself without troubling others; not",
	"to intermeddle in the affairs of others, and not easily to listen to",
	"slanders against them.",
	"6. Of Diognetus I had the lesson not to busy myself about vain things;",
	"not to credit the great professions of such as pretend to work",
	"wonders, or of sorcerers about their charms, and their expelling of",
	"Demons and the like; not to keep quails (for fighting or divination),",
	"nor to run after such things; to suffer freedom of speech in others,",
	"and to apply myself heartily to philosophy. Him also I must thank for",
	"my hearing first Bacchius, then Tandasis and Marcianus; that I wrote",
	"dialogues in my youth, and took a liking to the philosopher's pallet",
	"and skins, and to the other things which, by the Grecian discipline,",
	"belong to that profession.",
	"7. To Rusticus I owe my first apprehensions that my nature needed",
	"reform and cure; and that I did not fall into the ambition of the",
	"common Sophists, either by composing speculative writings or by",
	"declaiming harangues of exhortation in public; further, that I never",
	"strove to be admired by ostentation of great patience in an ascetic",
	"life, or by display of activity and application; that I gave over the",
	"study of rhetoric, poetry, and the graces of language; and that I did",
	"not pace my house in my senatorial robes, or practise any similar",
	"affectation. I observed also the simplicity of style in his letters,",
	"particularly in that which he wrote to my mother from Sinuessa. I",
	"learned from him to be easily appeased, and to be readily reconciled",
	"with those who had displeased me or given cause of offence, so soon as",
	"they inclined to make their peace; to read with care; not to rest",
	"satisfied with a slight and superficial knowledge; nor quickly to",
	"assent to great talkers. I have him to thank that I met with the",
}
