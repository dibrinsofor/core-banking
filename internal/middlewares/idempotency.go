package middlewares

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	redisstore "github.com/dibrinsofor/core-banking/internal/redis"
	"github.com/gin-gonic/gin"
)

type header struct {
	IdempotencyKey string `header:"Idempotency-Key" binding:"required"`
}

func Idempotency() gin.HandlerFunc {
	return func(c *gin.Context) {
		var h header

		w := &responseBodyWriter{ResponseWriter: c.Writer, body: bytes.NewBufferString("")}
		c.Writer = w

		if err := c.ShouldBindHeader(&h); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "no idempotency key in header. check documentation: https://github.com/dibrinsofor/core-banking/blob/master/Readme.MD",
			})
			c.Abort()
		}

		// check if key exists in redis
		user, err := redisstore.FindIdempKey(h.IdempotencyKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "uhoh, something went wrong. check documentation: https://github.com/dibrinsofor/core-banking/blob/master/Readme.MD",
			})
			c.Abort()
		}

		user.IdempotencyKey = h.IdempotencyKey
		user.CreatedAt = time.Now()
		user.RespCode = w.ResponseWriter.Status()
		user.RespMessage = w.body.String()
		deadline := user.SetDeadline()

		// can we check request? if nothing has changed return same error. don't process.
		if user.RespCode != 0 && time.Now().Before(deadline) {
			// key has not expired, return response
			// would we ever want to retry 500 errors?
			if user.RespCode == http.StatusOK || user.RespCode == http.StatusCreated || user.RespCode == http.StatusInternalServerError {
				w.ResponseWriter.WriteHeader(w.code)
				w.ResponseWriter.Write(w.body.Bytes())
				c.Abort()
			}
		}

		// retry if request was a 400 (bad request) or if user.RespCode is 0
		finished := make(chan struct{})
		panicChan := make(chan interface{}, 1)

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()

			c.Next()
			finished <- struct{}{}
		}()

		select {
		case <-panicChan:
			// rollover postgres?

			if err := redisstore.AddIdempotencyKey(user); err != nil {
				w.ResponseWriter.WriteHeader(http.StatusInternalServerError)
				eResp, _ := json.Marshal(gin.H{"message": "uhoh, something went wrong. check documentation: https://github.com/dibrinsofor/core-banking/blob/master/Readme.MD"})
				w.ResponseWriter.Write(eResp)
			}
			w.ResponseWriter.WriteHeader(w.code)
			w.ResponseWriter.Write(w.body.Bytes())
		case <-finished:
			// if finished, set headers and write resp
			w.mu.Lock()
			defer w.mu.Unlock()

			if err := redisstore.AddIdempotencyKey(user); err != nil {
				w.ResponseWriter.WriteHeader(http.StatusInternalServerError)
				eResp, _ := json.Marshal(gin.H{"message": "uhoh, something went wrong. check documentation: https://github.com/dibrinsofor/core-banking/blob/master/Readme.MD"})
				w.ResponseWriter.Write(eResp)
			}

			dst := w.ResponseWriter.Header()
			for k, vv := range w.Header() {
				dst[k] = vv
			}
		}
	}
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
	mu   sync.Mutex
	code int
}

func (r *responseBodyWriter) Write(b []byte) (n int, err error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
