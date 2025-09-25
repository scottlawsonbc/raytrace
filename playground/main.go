// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.
package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/phys"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	mux := http.NewServeMux()
	loggedMux := loggingMiddleware(mux)

	// Serve static files
	mux.Handle("/raytrace/playground/static/", http.StripPrefix("/raytrace/playground/static/", http.FileServer(http.Dir("./static"))))

	// Serve the playground HTML
	mux.HandleFunc("/raytrace/playground", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		http.ServeFile(w, r, "./static/playground.html")
	})

	// Status endpoint
	mux.HandleFunc("/raytrace/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		response := map[string]string{"status": "ok"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Render endpoint
	mux.HandleFunc("/raytrace/render", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
			http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Parse JSON into Scene
		var scene phys.Scene
		err = json.Unmarshal(body, &scene)
		if err != nil {
			log.Printf("JSON Unmarshal error: %v", err)
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{"error": "Invalid JSON: `+err.Error()+`"}`, http.StatusBadRequest)
			return
		}

		// Render the scene with a timeout
		type renderResult struct {
			Image image.Image
			Err   error
		}
		renderCh := make(chan renderResult, 1)

		go func() {
			// Call phys.Render which now returns (reconstruction, error)
			recon, err := phys.Render(context.Background(), &scene)
			if err != nil {
				renderCh <- renderResult{nil, err}
				return
			}
			// Create a montage of rendered images
			renderCh <- renderResult{recon.Image, nil}
		}()

		select {
		case res := <-renderCh:
			if res.Err != nil {
				log.Printf("Raytracer error: %v", res.Err)
				w.Header().Set("Content-Type", "application/json")
				errorMsg := fmt.Sprintf("Render Error: %v", res.Err)
				response := map[string]string{"error": errorMsg}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
				return
			}

			// Encode the image to PNG
			var buf bytes.Buffer
			err = png.Encode(&buf, res.Image)
			if err != nil {
				log.Printf("Error encoding image: %v", err)
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, `{"error": "Failed to encode image"}`, http.StatusInternalServerError)
				return
			}

			// Base64 encode the PNG
			encoded := base64.StdEncoding.EncodeToString(buf.Bytes())

			// Send the response
			response := map[string]string{"image": encoded}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)

		case <-time.After(30 * time.Second):
			log.Printf("Raytracer render timed out")
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{"error": "Raytracer render timed out"}`, http.StatusGatewayTimeout)
		}

		log.Printf("Render request processed in %v", time.Since(start))
	})

	addr := ":8020"
	log.Printf("Starting server at http://localhost%s/raytrace/playground", addr)
	err := http.ListenAndServe(addr, loggedMux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// loggingMiddleware logs each request with method, path, status code, and duration
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lrw, r)
		duration := time.Since(start)
		log.Printf("%s %s %d %v", r.Method, r.URL.Path, lrw.statusCode, duration)
	})
}

// loggingResponseWriter wraps http.ResponseWriter to capture the status code
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
