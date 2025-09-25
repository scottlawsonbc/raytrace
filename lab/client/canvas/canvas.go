//go:build js && wasm

package canvas

import (
	"syscall/js"
)

// Context wraps the JavaScript canvas rendering context.
type Context struct {
	Val *js.Value
}

// Width returns the width of the canvas.
func (ctx *Context) Width() float64 {
	return ctx.Val.Get("canvas").Get("width").Float()
}

// Height returns the height of the canvas.
func (ctx *Context) Height() float64 {
	return ctx.Val.Get("canvas").Get("height").Float()
}

// Translate applies a translation transformation to the canvas.
func (ctx *Context) Translate(x, y float64) {
	ctx.Val.Call("translate", x, y)
}

// Rotate applies a rotation transformation to the canvas.
func (ctx *Context) Rotate(angle float64) {
	ctx.Val.Call("rotate", angle)
}

// BeginPath starts a new path on the canvas.
func (ctx *Context) BeginPath() {
	ctx.Val.Call("beginPath")
}

// ClosePath closes the current path on the canvas.
func (ctx *Context) ClosePath() {
	ctx.Val.Call("closePath")
}

// Stroke strokes the current path with the current stroke style.
func (ctx *Context) Stroke() {
	ctx.Val.Call("stroke")
}

// Fill fills the current path with the current fill style.
func (ctx *Context) Fill() {
	ctx.Val.Call("fill")
}

// Rect adds a rectangle to the current path.
func (ctx *Context) Rect(x, y, w, h float64) {
	ctx.Val.Call("rect", x, y, w, h)
}

// FillRect draws a filled rectangle on the canvas.
func (ctx *Context) FillRect(x, y, w, h float64) {
	ctx.Val.Call("fillRect", x, y, w, h)
}

// StrokeRect draws a rectangular outline on the canvas.
func (ctx *Context) StrokeRect(x, y, w, h float64) {
	ctx.Val.Call("strokeRect", x, y, w, h)
}

// MoveTo moves the starting point of a new sub-path to the specified coordinates.
func (ctx *Context) MoveTo(x, y float64) {
	ctx.Val.Call("moveTo", x, y)
}

// LineTo adds a straight line to the current path.
func (ctx *Context) LineTo(x, y float64) {
	ctx.Val.Call("lineTo", x, y)
}

// FillStyle sets the fill style used for drawing shapes.
func (ctx *Context) FillStyle(s string) {
	ctx.Val.Set("fillStyle", s)
}

// StrokeStyle sets the stroke style used for drawing lines.
func (ctx *Context) StrokeStyle(s string) {
	ctx.Val.Set("strokeStyle", s)
}

// Arc adds an arc to the current path.
func (ctx *Context) Arc(x, y, radius, startAngle, endAngle float64, clockwise bool) {
	ctx.Val.Call("arc", x, y, radius, startAngle, endAngle, clockwise)
}

// DrawImage draws an image onto the canvas with specified source and destination parameters.
func (ctx *Context) DrawImage(src js.Value, sx, sy, sw, sh, dx, dy, dw, dh float64) {
	ctx.Val.Call("drawImage", src, sx, sy, sw, sh, dx, dy, dw, dh)
}

// GetImageData retrieves image data for a specified rectangle on the canvas.
func (ctx *Context) GetImageData(x, y, w, h float64) js.Value {
	return ctx.Val.Call("getImageData", x, y, w, h)
}

// PutImageData places image data onto the canvas at the specified coordinates.
func (ctx *Context) PutImageData(img js.Value, x, y float64) {
	ctx.Val.Call("putImageData", img, x, y)
}

// ClearRect clears the specified rectangle area on the canvas.
func (ctx *Context) ClearRect(x, y, w, h float64) {
	ctx.Val.Call("clearRect", x, y, w, h)
}
