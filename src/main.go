package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	datastar "github.com/starfederation/datastar/sdk/go"
)

const (
	port = 9001
)

func main() {
	r := gin.Default()

	// Serve static files
	r.Static("/css", "./src/html/css")
	r.Static("/img", "./src/html/img")

	// Serve the main HTML page
	r.GET("/", func(c *gin.Context) {
		c.File("./src/html/index.html")
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
				frag := fmt.Sprintf(`<span id="feed" style="color:#%s;border:1px solid #%s;border-radius:0.25rem;padding:1rem;">%s</span>`, hexString, hexString, hexString)

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
