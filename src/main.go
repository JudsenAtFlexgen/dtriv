package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	datastar "github.com/starfederation/datastar/sdk/go"

	"github.com/JudsenAtFlexgen/dtriv/trivia"
)

const (
	port = 9001
)

type question struct {
	Quest  string
	Opts   []string
	Answer uint
}

// Needed to convert from templ to gin context
func render(c *gin.Context, status int, template templ.Component) error {
	c.Status(status)
	return template.Render(c.Request.Context(), c.Writer)
}

// Needed to convert from templ to gin context
func renderStr(template templ.Component, buf *bytes.Buffer) error {
	return template.Render(context.Background(), buf)
}

func main() {
	r := gin.Default()

	// Serve static files
	r.Static("/css", "./src/html/css")
	r.Static("/img", "./src/html/img")

	// Serve the main HTML page
	r.GET("/", func(c *gin.Context) {
		c.File("./src/html/index.html")
	})

	r.GET("/test", func(c *gin.Context) {
		w := c.Writer
		r := c.Request

		sse := datastar.NewSSE(w, r)
		color_bytes := make([]byte, 3)
		if _, err := rand.Read(color_bytes); err != nil {
			log.Println("Error generating random bytes:", err)
			return
		}
		hexString := hex.EncodeToString(color_bytes)

		var buf bytes.Buffer
		style := fmt.Sprintf("color:#%s;border:1px solid #%s;border-radius:0.25rem;padding:1rem;", hexString, hexString)
		component := trivia.Question(style, "Wowzers")
		err := renderStr(component, &buf)
		if err != nil {
			panic("ohno")
		}
		sse.MergeFragments(buf.String())
	})

	// SSE stream
	r.GET("/stream", func(c *gin.Context) {
		w := c.Writer
		r := c.Request

		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		sse := datastar.NewSSE(w, r)
		for {
			select {
			case <-r.Context().Done():
				log.Println("Client disconnected")
				return
			case <-ticker.C:
				bytes := make([]byte, 3)
				if _, err := rand.Read(bytes); err != nil {
					log.Println("Error generating random bytes:", err)
					return
				}
				hexString := hex.EncodeToString(bytes)
				frag := fmt.Sprintf(`<button id="test" style="color:#%s;border:1px solid #%s;border-radius:0.25rem;padding:1rem;">%s</button>`, hexString, hexString, hexString)

				sse.MergeFragments(frag)
			}
		}
	})

	// Start server
	addr := fmt.Sprintf("0.0.0.0:%d", port)
	log.Println("Starting server on", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
