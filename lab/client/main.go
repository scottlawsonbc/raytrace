//go:build js && wasm

package main

import (
	"fmt"
	"io/fs"
	"syscall/js"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/lab/client/canvas"
	"github.com/scottlawsonbc/slam/code/photon/raytrace/obj"
	"github.com/scottlawsonbc/slam/code/photon/raytrace/phys"
)

// logf logs a formatted message to the browser console with the specified color.
func logf(color string, format string, args ...interface{}) {
	console := js.Global().Get("console")
	if console.IsUndefined() {
		return
	}
	message := fmt.Sprintf(format, args...)
	styledMessage := fmt.Sprintf("%%c%s", message)
	css := fmt.Sprintf("color: %s;", color)
	console.Call("log", styledMessage, css)
}

func loadNodes(fsys fs.FS, objPath string) ([]phys.Node, error) {
	parsedObj, err := obj.ParseFS(fsys, objPath)
	if err != nil {
		return nil, err
	}
	nodes, err := phys.ConvertObjectToNodes(parsedObj, fsys)
	if err != nil {
		return nil, err
	}
	fmt.Printf("got %d nodes\n", len(nodes))
	if len(nodes) == 0 {
		return nil, fmt.Errorf("no nodes found in OBJ file")
	}
	for i, n := range nodes {
		fmt.Printf("Node %d: %v\n", i, n.Shape.Bounds())
	}
	return nodes, nil
}

func renderContext() canvas.Context {
	doc := js.Global().Get("document")
	node := doc.Call("getElementById", "canvas")
	ctx := node.Call("getContext", "2d")
	return canvas.Context{Val: &ctx}
}

func main() {
	done := make(chan struct{})
	doc := js.Global().Get("document")

	var labApp = &app{
		eventOut:   make(chan any),
		eventIn:    make(chan any),
		glctx:      renderContext(),
	}
	labApp.initWorker()

	labApp.eventIn = pump(labApp.eventOut)
	labApp.RegisterFilter(logEventFilter)
	go autoResize(labApp.eventIn)

	handleKeyEvent := redirectKeyboardEvent(labApp.eventIn)
	defer handleKeyEvent.Release()
	doc.Call("addEventListener", "keydown", handleKeyEvent)
	doc.Call("addEventListener", "keyup", handleKeyEvent)

	handleMouseEvent := redirectMouseEvent(labApp.eventIn)
	defer handleMouseEvent.Release()
	doc.Call("addEventListener", "mousedown", handleMouseEvent)
	doc.Call("addEventListener", "mouseup", handleMouseEvent)
	doc.Call("addEventListener", "mousemove", handleMouseEvent)

	handleWheelEvent := redirectWheelEvent(labApp.eventIn)
	defer handleWheelEvent.Release()
	doc.Call("addEventListener", "wheel", handleWheelEvent)

	handleContextMenu := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		args[0].Call("preventDefault")
		return nil
	})
	doc.Call("addEventListener", "contextmenu", handleContextMenu)
	doc.Call("getElementById", "loading").Get("style").Set("visibility", "hidden")
	labApp.main()
	<-done
}
