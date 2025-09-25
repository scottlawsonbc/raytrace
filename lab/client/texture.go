//go:build js && wasm

package main

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"syscall/js"
)

type texture struct {
	jsImg js.Value
	goImg image.Image
}

func (t texture) ColorModel() color.Model {
	return t.goImg.ColorModel()
}

func (t texture) At(x, y int) color.Color {
	return t.goImg.At(x, y)
}

func (t texture) Bounds() image.Rectangle {
	return t.goImg.Bounds()
}

func (t texture) Free() {
	t.jsImg.Call("remove")
}

func NewTexture(img image.Image) texture {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		panic(err)
	}
	b64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	dataURL := "data:image/png;base64," + b64
	jsImg := js.Global().Get("Image").New()
	loaded := make(chan struct{})
	closer := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		close(loaded)
		return nil
	})
	jsImg.Call("addEventListener", "load", closer)
	jsImg.Set("src", dataURL)
	<-loaded
	closer.Release()
	return texture{jsImg: jsImg, goImg: img}
}
