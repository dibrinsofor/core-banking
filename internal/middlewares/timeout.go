package middlewares

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tw timeoutWriter

		tw.ResponseWriter = c.Writer
		tw.h = make(http.Header)
		tw.body = bytes.NewBufferString("")

		c.Writer = &tw

		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		finished := make(chan struct{})
		panicChan := make(chan interface{}, 1)

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
					fmt.Println(p)
				}
			}()

			c.Next()
			finished <- struct{}{}
		}()

		select {
		case <-panicChan:
			tw.ResponseWriter.Header().Add("Content-Type", "application/json")
			tw.ResponseWriter.WriteHeader(tw.ResponseWriter.Status())
			tw.ResponseWriter.Write(tw.body.Bytes())
		case <-finished:
			tw.mu.Lock()
			defer tw.mu.Unlock()

			// map Headers from tw.Header()
			dst := tw.ResponseWriter.Header()
			for k, vv := range tw.Header() {
				dst[k] = vv
			}
			tw.ResponseWriter.Header().Set("Content-Type", "application/json")
			tw.ResponseWriter.Write(tw.body.Bytes())
		case <-ctx.Done():
			tw.mu.Lock()
			defer tw.mu.Unlock()

			tw.ResponseWriter.Header().Set("Content-Type", "application/json")
			tw.ResponseWriter.WriteHeader(http.StatusRequestTimeout)
			eResp, _ := json.Marshal(gin.H{
				"message": "Request timed out.",
			})
			tw.ResponseWriter.Write(eResp)
			c.Abort()
			tw.SetTimedOut()
		}
	}
}

type timeoutWriter struct {
	gin.ResponseWriter
	h           http.Header
	body        *bytes.Buffer
	mu          sync.Mutex
	timedOut    bool
	wroteHeader bool
	code        int
}

func (tw *timeoutWriter) Write(b []byte) (int, error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	return tw.body.Write(b)
}

func (tw *timeoutWriter) WriteHeader(code int) {
	checkWriteHeaderCode(code)
	tw.mu.Lock()
	defer tw.mu.Unlock()

	if tw.timedOut || tw.wroteHeader {
		return
	}
	tw.writeHeader(code)
}

func (tw *timeoutWriter) writeHeader(code int) {
	tw.wroteHeader = true
	tw.code = code
}

func (tw *timeoutWriter) Header() http.Header {
	return tw.h
}

func (tw *timeoutWriter) SetTimedOut() {
	tw.timedOut = true
}

func checkWriteHeaderCode(code int) {
	if code < 100 || code > 999 {
		panic(fmt.Sprintf("invalid WriteHeader code %v", code))
	}
}
