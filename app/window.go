// SPDX-License-Identifier: Unlicense OR MIT

package app

import (
	"errors"
	"image"
	"image/color"
	"runtime"
	"sync"
	"time"
	"unicode/utf8"

	"gio.mleku.dev/f32"
	"gio.mleku.dev/gpu"
	"gio.mleku.dev/internal/debug"
	"gio.mleku.dev/internal/ops"
	"gio.mleku.dev/io/event"
	"gio.mleku.dev/io/input"
	"gio.mleku.dev/io/key"
	"gio.mleku.dev/io/pointer"
	"gio.mleku.dev/io/system"
	"gio.mleku.dev/op"
	"gio.mleku.dev/unit"
	"lol.mleku.dev/log"
)

// Option configures a window.
type Option func(unit.Metric, *Config)

// Window represents an operating system window.
//
// The zero-value Window is useful; the GUI window is created and shown the first
// time the [Event] method is called. On iOS or Android, the first Window represents
// the window previously created by the platform.
//
// More than one Window is not supported on iOS, Android, WebAssembly.
type Window struct {
	initialOpts    []Option
	initialActions []system.Action

	ctx context
	gpu gpu.GPU
	// timer tracks the delayed invalidate goroutine.
	timer struct {
		// quit is shuts down the goroutine.
		quit chan struct{}
		// update the invalidate time.
		update chan time.Time
	}

	animating    bool
	hasNextFrame bool
	nextFrame    time.Time
	// viewport is the latest frame size with insets applied.
	viewport image.Rectangle
	// metric is the metric from the most recent frame.
	metric      unit.Metric
	queue       input.Router
	cursor      pointer.Cursor
	decorations struct {
		op.Ops
		// enabled tracks the Decorated option as
		// given to the Option method. It may differ
		// from Config.Decorated depending on platform
		// capability.
		enabled bool
		Config
		height        unit.Dp
		currentHeight int
	}
	nocontext bool
	// semantic data, lazily evaluated if requested by a backend to speed up
	// the cases where semantic data is not needed.
	semantic struct {
		// uptodate tracks whether the fields below are up to date.
		uptodate bool
		root     input.SemanticID
		prevTree []input.SemanticNode
		tree     []input.SemanticNode
		ids      map[input.SemanticID]input.SemanticNode
	}
	imeState editorState
	driver   driver
	// gpuErr tracks the GPU error that is to be reported when
	// the window is closed.
	gpuErr error

	// invMu protects mayInvalidate.
	invMu         sync.Mutex
	mayInvalidate bool

	// coalesced tracks the most recent events waiting to be delivered
	// to the client.
	coalesced eventSummary
	// frame tracks the most recent frame event.
	lastFrame struct {
		sync bool
		size image.Point
		off  image.Point
		deco op.CallOp
	}
}

type eventSummary struct {
	wakeup       bool
	cfg          *ConfigEvent
	view         *ViewEvent
	frame        *frameEvent
	framePending bool
	destroy      *DestroyEvent
	windowMouse  *WindowMouseEvent
}

type callbacks struct {
	w *Window
}

func (w *Window) validateAndProcess(size image.Point, sync bool, frame *op.Ops, sigChan chan<- struct{}) error {
	signal := func() {
		if sigChan != nil {
			// We're done with frame, let the client continue.
			sigChan <- struct{}{}
			// Signal at most once.
			sigChan = nil
		}
	}
	defer signal()
	for {
		if w.gpu == nil && !w.nocontext {
			var err error
			if w.ctx == nil {
				if w.ctx, err = w.driver.NewContext(); err != nil {
					return err
				}
				if err = w.ctx.Lock(); err != nil {
					w.destroyGPU()
					return err
				}
				sync = true
			}
		}
		if sync && w.ctx != nil {
			if err := w.ctx.Refresh(); err != nil {
				if errors.Is(err, errOutOfDate) {
					// Surface couldn't be created for transient reasons. Skip
					// this frame and wait for the next.
					return nil
				}
				w.destroyGPU()
				if errors.Is(err, gpu.ErrDeviceLost) {
					continue
				}
				return err
			}
		}
		if w.gpu == nil && !w.nocontext {
			gpu, err := gpu.New(w.ctx.API())
			if err != nil {
				w.ctx.Unlock()
				w.destroyGPU()
				return err
			}
			w.gpu = gpu
		}
		if w.gpu != nil {
			if err := w.frame(frame, size); err != nil {
				w.ctx.Unlock()
				if errors.Is(err, errOutOfDate) {
					// GPU surface needs refreshing.
					sync = true
					continue
				}
				w.destroyGPU()
				if errors.Is(err, gpu.ErrDeviceLost) {
					continue
				}
				return err
			}
		}
		w.queue.Frame(frame)
		// Let the client continue as soon as possible, in particular before
		// a potentially blocking Present.
		signal()
		var err error
		if w.gpu != nil {
			err = w.ctx.Present()
		}
		return err
	}
}

func (w *Window) frame(frame *op.Ops, viewport image.Point) error {
	if runtime.GOOS == "js" {
		// Use transparent black when Gio is embedded, to allow mixing of Gio and
		// foreign content below.
		w.gpu.Clear(color.NRGBA{A: 0x00, R: 0x00, G: 0x00, B: 0x00})
	} else {
		w.gpu.Clear(color.NRGBA{A: 0xff, R: 0xff, G: 0xff, B: 0xff})
	}
	target, err := w.ctx.RenderTarget()
	if err != nil {
		return err
	}
	return w.gpu.Frame(frame, target, viewport)
}

func (w *Window) processFrame(frame *op.Ops, ack chan<- struct{}) {
	w.coalesced.framePending = false
	wrapper := &w.decorations.Ops
	off := op.Offset(w.lastFrame.off).Push(wrapper)
	ops.AddCall(&wrapper.Internal, &frame.Internal, ops.PC{}, ops.PCFor(&frame.Internal))
	off.Pop()
	w.lastFrame.deco.Add(wrapper)
	if err := w.validateAndProcess(w.lastFrame.size, w.lastFrame.sync, wrapper, ack); err != nil {
		w.destroyGPU()
		w.gpuErr = err
		w.driver.Perform(system.ActionClose)
		return
	}
	w.updateState()
	w.updateCursor()
}

func (w *Window) updateState() {
	for k := range w.semantic.ids {
		delete(w.semantic.ids, k)
	}
	w.semantic.uptodate = false
	q := &w.queue
	switch q.TextInputState() {
	case input.TextInputOpen:
		w.driver.ShowTextInput(true)
	case input.TextInputClose:
		w.driver.ShowTextInput(false)
	}
	if hint, ok := q.TextInputHint(); ok {
		w.driver.SetInputHint(hint)
	}
	if mime, txt, ok := q.WriteClipboard(); ok {
		w.driver.WriteClipboard(mime, txt)
	}
	if txt, ok := q.WritePrimaryClipboard(); ok {
		log.I.F("📤 APP: Writing to primary clipboard: %q\n", txt)
		w.driver.WritePrimaryClipboard(txt)
	}
	if q.ClipboardRequested() {
		w.driver.ReadClipboard()
	}
	if q.PrimaryClipboardRequested() {
		log.I.Ln("📥 APP: Reading from primary clipboard")
		w.driver.ReadPrimaryClipboard()
	}
	oldState := w.imeState
	newState := oldState
	newState.EditorState = q.EditorState()
	if newState != oldState {
		w.imeState = newState
		w.driver.EditorStateChanged(oldState, newState)
	}
	if t, ok := q.WakeupTime(); ok {
		w.setNextFrame(t)
	}
	w.updateAnimation()
}

// Invalidate the window such that a [FrameEvent] will be generated immediately.
// If the window is inactive, an unspecified event is sent instead.
//
// Note that Invalidate is intended for externally triggered updates, such as a
// response from a network request. The [op.InvalidateCmd] command is more efficient
// for animation.
//
// Invalidate is safe for concurrent use.
func (w *Window) Invalidate() {
	w.invMu.Lock()
	defer w.invMu.Unlock()
	if w.mayInvalidate {
		w.mayInvalidate = false
		w.driver.Invalidate()
	}
}

// Option applies the options to the window. The options are hints; the platform is
// free to ignore or adjust them.
func (w *Window) Option(opts ...Option) {
	if len(opts) == 0 {
		return
	}
	if w.driver == nil {
		w.initialOpts = append(w.initialOpts, opts...)
		return
	}
	w.Run(func() {
		cnf := Config{Decorated: w.decorations.enabled}
		for _, opt := range opts {
			opt(w.metric, &cnf)
		}
		w.decorations.enabled = cnf.Decorated
		// Using server-side decorations, no height option needed
		w.driver.Configure(opts)
		w.setNextFrame(time.Time{})
		w.updateAnimation()
	})
}

// Run f in the same thread as the native window event loop, and wait for f to
// return or the window to close. If the window has not yet been created,
// Run calls f directly.
//
// Note that most programs should not call Run; configuring a Window with
// [CustomRenderer] is a notable exception.
func (w *Window) Run(f func()) {
	if w.driver == nil {
		f()
		return
	}
	done := make(chan struct{})
	w.driver.Run(func() {
		defer close(done)
		f()
	})
	<-done
}

func (w *Window) updateAnimation() {
	if w.driver == nil {
		return
	}
	animate := false
	if w.hasNextFrame {
		if dt := time.Until(w.nextFrame); dt <= 0 {
			animate = true
		} else {
			// Schedule redraw.
			w.scheduleInvalidate(w.nextFrame)
		}
	}
	if animate != w.animating {
		w.animating = animate
		w.driver.SetAnimating(animate)
	}
}

func (w *Window) scheduleInvalidate(t time.Time) {
	if w.timer.quit == nil {
		w.timer.quit = make(chan struct{})
		w.timer.update = make(chan time.Time)
		go func() {
			var timer *time.Timer
			for {
				var timeC <-chan time.Time
				if timer != nil {
					timeC = timer.C
				}
				select {
				case <-w.timer.quit:
					w.timer.quit <- struct{}{}
					return
				case t := <-w.timer.update:
					if timer != nil {
						timer.Stop()
					}
					timer = time.NewTimer(time.Until(t))
				case <-timeC:
					w.Invalidate()
				}
			}
		}()
	}
	w.timer.update <- t
}

func (w *Window) setNextFrame(at time.Time) {
	if !w.hasNextFrame || at.Before(w.nextFrame) {
		w.hasNextFrame = true
		w.nextFrame = at
	}
}

func (c *callbacks) SetDriver(d driver) {
	if d == nil {
		panic("nil driver")
	}
	c.w.invMu.Lock()
	defer c.w.invMu.Unlock()
	c.w.driver = d
}

func (c *callbacks) ProcessFrame(frame *op.Ops, ack chan<- struct{}) {
	c.w.processFrame(frame, ack)
}

func (c *callbacks) ProcessEvent(e event.Event) bool {
	return c.w.processEvent(e)
}

// SemanticRoot returns the ID of the semantic root.
func (c *callbacks) SemanticRoot() input.SemanticID {
	c.w.updateSemantics()
	return c.w.semantic.root
}

// LookupSemantic looks up a semantic node from an ID. The zero ID denotes the root.
func (c *callbacks) LookupSemantic(semID input.SemanticID) (input.SemanticNode, bool) {
	c.w.updateSemantics()
	n, found := c.w.semantic.ids[semID]
	return n, found
}

func (c *callbacks) AppendSemanticDiffs(diffs []input.SemanticID) []input.SemanticID {
	c.w.updateSemantics()
	if tree := c.w.semantic.prevTree; len(tree) > 0 {
		c.w.collectSemanticDiffs(&diffs, c.w.semantic.prevTree[0])
	}
	return diffs
}

func (c *callbacks) SemanticAt(pos f32.Point) (input.SemanticID, bool) {
	c.w.updateSemantics()
	return c.w.queue.SemanticAt(pos)
}

func (c *callbacks) EditorState() editorState {
	return c.w.imeState
}

func (c *callbacks) SetComposingRegion(r key.Range) {
	c.w.imeState.compose = r
}

func (c *callbacks) EditorInsert(text string) {
	sel := c.w.imeState.Selection.Range
	c.EditorReplace(sel, text)
	start := min(sel.End, sel.Start)
	sel.Start = start + utf8.RuneCountInString(text)
	sel.End = sel.Start
	c.SetEditorSelection(sel)
}

func (c *callbacks) EditorReplace(r key.Range, text string) {
	c.w.imeState.Replace(r, text)
	c.w.driver.ProcessEvent(key.EditEvent{Range: r, Text: text})
	c.w.driver.ProcessEvent(key.SnippetEvent(c.w.imeState.Snippet.Range))
}

func (c *callbacks) SetEditorSelection(r key.Range) {
	c.w.imeState.Selection.Range = r
	c.w.driver.ProcessEvent(key.SelectionEvent(r))
}

func (c *callbacks) SetEditorSnippet(r key.Range) {
	if sn := c.EditorState().Snippet.Range; sn == r {
		// No need to expand.
		return
	}
	c.w.driver.ProcessEvent(key.SnippetEvent(r))
}

func (w *Window) moveFocus(dir key.FocusDirection) {
	w.queue.MoveFocus(dir)
	if _, handled := w.queue.WakeupTime(); handled {
		w.queue.RevealFocus(w.viewport)
	} else {
		var v image.Point
		switch dir {
		case key.FocusRight:
			v = image.Pt(+1, 0)
		case key.FocusLeft:
			v = image.Pt(-1, 0)
		case key.FocusDown:
			v = image.Pt(0, +1)
		case key.FocusUp:
			v = image.Pt(0, -1)
		default:
			return
		}
		const scrollABit = unit.Dp(50)
		dist := v.Mul(int(w.metric.Dp(scrollABit)))
		w.queue.ScrollFocus(dist)
	}
}

func (c *callbacks) ClickFocus() {
	c.w.queue.ClickFocus()
	c.w.setNextFrame(time.Time{})
	c.w.updateAnimation()
}

func (c *callbacks) ActionAt(p f32.Point) (system.Action, bool) {
	return c.w.queue.ActionAt(p)
}

func (w *Window) destroyGPU() {
	if w.gpu != nil {
		w.ctx.Lock()
		w.gpu.Release()
		w.ctx.Unlock()
		w.gpu = nil
	}
	if w.ctx != nil {
		w.ctx.Release()
		w.ctx = nil
	}
}

// updateSemantics refreshes the semantics tree, the id to node map and the ids of
// updated nodes.
func (w *Window) updateSemantics() {
	if w.semantic.uptodate {
		return
	}
	w.semantic.uptodate = true
	w.semantic.prevTree, w.semantic.tree = w.semantic.tree, w.semantic.prevTree
	w.semantic.tree = w.queue.AppendSemantics(w.semantic.tree[:0])
	w.semantic.root = w.semantic.tree[0].ID
	for _, n := range w.semantic.tree {
		w.semantic.ids[n.ID] = n
	}
}

// collectSemanticDiffs traverses the previous semantic tree, noting changed nodes.
func (w *Window) collectSemanticDiffs(diffs *[]input.SemanticID, n input.SemanticNode) {
	newNode, exists := w.semantic.ids[n.ID]
	// Ignore deleted nodes, as their disappearance will be reported through an
	// ancestor node.
	if !exists {
		return
	}
	diff := newNode.Desc != n.Desc || len(n.Children) != len(newNode.Children)
	for i, ch := range n.Children {
		if !diff {
			newCh := newNode.Children[i]
			diff = ch.ID != newCh.ID
		}
		w.collectSemanticDiffs(diffs, ch)
	}
	if diff {
		*diffs = append(*diffs, n.ID)
	}
}

func (c *callbacks) Invalidate() {
	c.w.setNextFrame(time.Time{})
	c.w.updateAnimation()
	// Guarantee a wakeup, even when not animating.
	c.w.processEvent(wakeupEvent{})
}

func (c *callbacks) nextEvent() (event.Event, bool) {
	return c.w.nextEvent()
}

func (w *Window) nextEvent() (event.Event, bool) {
	s := &w.coalesced
	defer func() {
		// Every event counts as a wakeup.
		s.wakeup = false
	}()
	switch {
	case s.framePending:
		// If the user didn't call FrameEvent.Event, process
		// an empty frame.
		w.processFrame(new(op.Ops), nil)
	case s.view != nil:
		e := *s.view
		s.view = nil
		return e, true
	case s.destroy != nil:
		e := *s.destroy
		// Clear pending events after DestroyEvent is delivered.
		*s = eventSummary{}
		return e, true
	case s.cfg != nil:
		e := *s.cfg
		s.cfg = nil
		return e, true
	case s.windowMouse != nil:
		e := *s.windowMouse
		s.windowMouse = nil
		return e, true
	case s.frame != nil:
		e := *s.frame
		s.frame = nil
		s.framePending = true
		return e.FrameEvent, true
	case s.wakeup:
		return wakeupEvent{}, true
	}
	w.invMu.Lock()
	defer w.invMu.Unlock()
	w.mayInvalidate = w.driver != nil
	return nil, false
}

func (w *Window) processEvent(e event.Event) bool {
	switch e2 := e.(type) {
	case wakeupEvent:
		w.coalesced.wakeup = true
	case frameEvent:
		if e2.Size == (image.Point{}) {
			panic(errors.New("internal error: zero-sized Draw"))
		}
		w.metric = e2.Metric
		w.hasNextFrame = false
		e2.Frame = w.driver.Frame
		e2.Source = w.queue.Source()
		// Prepare the decorations and update the frame insets.
		viewport := image.Rectangle{
			Min: image.Point{
				X: e2.Metric.Dp(e2.Insets.Left),
				Y: e2.Metric.Dp(e2.Insets.Top),
			},
			Max: image.Point{
				X: e2.Size.X - e2.Metric.Dp(e2.Insets.Right),
				Y: e2.Size.Y - e2.Metric.Dp(e2.Insets.Bottom),
			},
		}
		// Scroll to focus if viewport is shrinking in any dimension.
		if old, new := w.viewport.Size(), viewport.Size(); new.X < old.X || new.Y < old.Y {
			w.queue.RevealFocus(viewport)
		}
		w.viewport = viewport
		wrapper := &w.decorations.Ops
		wrapper.Reset()
		m := op.Record(wrapper)
		offset := w.decorate(e2.FrameEvent, wrapper)
		w.lastFrame.deco = m.Stop()
		w.lastFrame.size = e2.Size
		w.lastFrame.sync = e2.Sync
		w.lastFrame.off = offset
		e2.Size = e2.Size.Sub(offset)
		w.coalesced.frame = &e2
	case DestroyEvent:
		if w.gpuErr != nil {
			e2.Err = w.gpuErr
		}
		w.destroyGPU()
		w.invMu.Lock()
		w.mayInvalidate = false
		w.driver = nil
		w.invMu.Unlock()
		if q := w.timer.quit; q != nil {
			q <- struct{}{}
			<-q
		}
		w.coalesced.destroy = &e2
	case ViewEvent:
		if !e2.Valid() && w.gpu != nil {
			w.ctx.Lock()
			w.gpu.Release()
			w.gpu = nil
			w.ctx.Unlock()
		}
		w.coalesced.view = &e2
	case ConfigEvent:
		// Using server-side decorations, no client-side decoration state needed
		wasFocused := w.decorations.Config.Focused
		w.decorations.Config = e2.Config
		e2.Config = w.effectiveConfig()
		w.coalesced.cfg = &e2
		if f := w.decorations.Config.Focused; f != wasFocused {
			w.queue.Queue(key.FocusEvent{Focus: f})
		}
		t, handled := w.queue.WakeupTime()
		if handled {
			w.setNextFrame(t)
			w.updateAnimation()
		}
		return handled
	case WindowMouseEvent:
		// Window-level mouse enter/leave events are handled directly by the application
		// They don't need to go through the input queue
		w.coalesced.windowMouse = &e2
	case event.Event:
		// Skip focus navigation processing - just pass through raw events
		w.queue.Queue(e2)
		t, handled := w.queue.WakeupTime()
		w.updateCursor()
		if handled {
			w.setNextFrame(t)
			w.updateAnimation()
		}
		return handled
	}
	return true
}

// Event blocks until an event is received from the window, such as
// [FrameEvent], or until [Invalidate] is called. The window is created
// and shown the first time Event is called.
func (w *Window) Event() event.Event {
	if w.driver == nil {
		w.init()
	}
	if w.driver == nil {
		e, ok := w.nextEvent()
		if !ok {
			panic("window initialization failed without a DestroyEvent")
		}
		return e
	}
	return w.driver.Event()
}

func (w *Window) init() {
	debug.Parse()
	// Using server-side decorations, no height measurement needed
	defaultOptions := []Option{
		Size(800, 600),
		Title("Gio"),
		Decorated(true),
	}
	options := append(defaultOptions, w.initialOpts...)
	w.initialOpts = nil
	var cnf Config
	cnf.apply(unit.Metric{}, options)

	w.nocontext = cnf.CustomRenderer
	// Using server-side decorations, no client-side decoration setup needed
	w.decorations.enabled = cnf.Decorated
	w.imeState.compose = key.Range{Start: -1, End: -1}
	w.semantic.ids = make(map[input.SemanticID]input.SemanticNode)
	newWindow(&callbacks{w}, options)
	for _, acts := range w.initialActions {
		w.Perform(acts)
	}
	w.initialActions = nil
}

func (w *Window) updateCursor() {
	if c := w.queue.Cursor(); c != w.cursor {
		w.cursor = c
		w.driver.SetCursor(c)
	}
}

func (w *Window) fallbackDecorate() bool {
	cnf := w.decorations.Config
	return w.decorations.enabled && !cnf.Decorated && cnf.Mode != Fullscreen && !w.nocontext
}

// decorate the window if enabled and returns the corresponding Insets.
func (w *Window) decorate(e FrameEvent, o *op.Ops) image.Point {
	// Using server-side decorations, no client-side decoration needed
	return image.Pt(0, 0)
}

func (w *Window) effectiveConfig() Config {
	cnf := w.decorations.Config
	cnf.Size.Y -= w.decorations.currentHeight
	cnf.Decorated = w.decorations.enabled || cnf.Decorated
	return cnf
}

// splitActions splits options from actions and return them and the remaining
// actions.
func splitActions(actions system.Action) ([]Option, system.Action) {
	var opts []Option
	walkActions(actions, func(action system.Action) {
		switch action {
		case system.ActionMinimize:
			opts = append(opts, Minimized.Option())
		case system.ActionMaximize:
			opts = append(opts, Maximized.Option())
		case system.ActionUnmaximize:
			opts = append(opts, Windowed.Option())
		case system.ActionFullscreen:
			opts = append(opts, Fullscreen.Option())
		default:
			return
		}
		actions &^= action
	})
	return opts, actions
}

// Perform the actions on the window.
func (w *Window) Perform(actions system.Action) {
	opts, acts := splitActions(actions)
	w.Option(opts...)
	if acts == 0 {
		return
	}
	if w.driver == nil {
		w.initialActions = append(w.initialActions, acts)
		return
	}
	w.Run(func() {
		w.driver.Perform(actions)
	})
}

// Title sets the title of the window.
func Title(t string) Option {
	return func(_ unit.Metric, cnf *Config) {
		cnf.Title = t
	}
}

// Size sets the size of the window. The mode will be changed to Windowed.
func Size(w, h unit.Dp) Option {
	if w <= 0 {
		panic("width must be larger than or equal to 0")
	}
	if h <= 0 {
		panic("height must be larger than or equal to 0")
	}
	return func(m unit.Metric, cnf *Config) {
		cnf.Mode = Windowed
		cnf.Size = image.Point{
			X: m.Dp(w),
			Y: m.Dp(h),
		}
	}
}

// MaxSize sets the maximum size of the window.
func MaxSize(w, h unit.Dp) Option {
	if w <= 0 {
		panic("width must be larger than or equal to 0")
	}
	if h <= 0 {
		panic("height must be larger than or equal to 0")
	}
	return func(m unit.Metric, cnf *Config) {
		cnf.MaxSize = image.Point{
			X: m.Dp(w),
			Y: m.Dp(h),
		}
	}
}

// MinSize sets the minimum size of the window.
func MinSize(w, h unit.Dp) Option {
	if w <= 0 {
		panic("width must be larger than or equal to 0")
	}
	if h <= 0 {
		panic("height must be larger than or equal to 0")
	}
	return func(m unit.Metric, cnf *Config) {
		cnf.MinSize = image.Point{
			X: m.Dp(w),
			Y: m.Dp(h),
		}
	}
}

// StatusColor sets the color of the Android status bar.
func StatusColor(color color.NRGBA) Option {
	return func(_ unit.Metric, cnf *Config) {
		cnf.StatusColor = color
	}
}

// NavigationColor sets the color of the navigation bar on Android, or the address bar in browsers.
func NavigationColor(color color.NRGBA) Option {
	return func(_ unit.Metric, cnf *Config) {
		cnf.NavigationColor = color
	}
}

// CustomRenderer controls whether the window contents is
// rendered by the client. If true, no GPU context is created.
//
// Caller must assume responsibility for rendering which includes
// initializing the render backend, swapping the framebuffer and
// handling frame pacing.
func CustomRenderer(custom bool) Option {
	return func(_ unit.Metric, cnf *Config) {
		cnf.CustomRenderer = custom
	}
}

// Decorated controls whether Gio and/or the platform are responsible
// for drawing window decorations. Providing false indicates that
// the application will either be undecorated or will draw its own decorations.
func Decorated(enabled bool) Option {
	return func(_ unit.Metric, cnf *Config) {
		cnf.Decorated = enabled
	}
}

// flushEvent is sent to detect when the user program
// has completed processing of all prior events. Its an
// [io/event.Event] but only for internal use.
type flushEvent struct{}

func (t flushEvent) ImplementsEvent() {}

// WindowMouseEvent represents a window-level mouse enter/leave event.
type WindowMouseEvent struct {
	Kind      WindowMouseKind
	Position  f32.Point
	Time      time.Duration
	Modifiers key.Modifiers
}

// WindowMouseKind represents the type of window-level mouse event.
type WindowMouseKind int

const (
	WindowMouseEnter WindowMouseKind = iota
	WindowMouseLeave
)

func (WindowMouseEvent) ImplementsEvent() {}

// theFlushEvent avoids allocating garbage when sending
// flushEvents.
var theFlushEvent flushEvent
