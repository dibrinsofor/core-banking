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

		w := &responseBodyWriter{ResponseWriter: c.Writer, body: &bytes.Buffer{}}
		c.Writer = w

		if err := c.ShouldBindHeader(&h); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "no idempotency key in header. check documentation: https://github.com/dibrinsofor/core-banking/blob/master/Readme.MD",
			})
			return
		}

		// check if key exists in redis
		user, err := redisstore.FindIdempKey(h.IdempotencyKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "uhoh, something went wrong. check documentation: https://github.com/dibrinsofor/core-banking/blob/master/Readme.MD",
			})
			return
		}

		user.IdempotencyKey = h.IdempotencyKey

		// check saved state
		// would we ever want to retry 500 errors?
		if user.RespCode == http.StatusOK || user.RespCode == http.StatusCreated || user.RespCode == http.StatusInternalServerError {
			c.JSON(user.RespCode, user.RespMessage)
			return
		}

		// can we check request? if nothing has changed return same error. don't process.
		// is saved state a thing? how do I monitor failures?
		// if user.RespCode == http.StatusBadRequest {
		// 	return
		// }

		// run if request was a 400 (bad request) or if user.RespCode is 0

		finished := make(chan struct{})        // handler finished
		panicChan := make(chan interface{}, 1) // handle panic if can't recover

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
			user.CreatedAt = time.Now()
			user.RespCode = w.ResponseWriter.Status()
			user.RespMessage = w.body.String()

			// if something fails we still want to store in redis
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

			user.CreatedAt = time.Now()
			user.RespCode = w.ResponseWriter.Status()
			user.RespMessage = w.body.String()

			if err := redisstore.AddIdempotencyKey(user); err != nil {
				w.ResponseWriter.WriteHeader(http.StatusInternalServerError)
				eResp, _ := json.Marshal(gin.H{"message": "uhoh, something went wrong. check documentation: https://github.com/dibrinsofor/core-banking/blob/master/Readme.MD"})
				w.ResponseWriter.Write(eResp)
			}

			dst := w.ResponseWriter.Header()
			for k, vv := range w.Header() {
				dst[k] = vv
			}
			w.ResponseWriter.WriteHeader(w.code)
			w.ResponseWriter.Write(w.body.Bytes())
		}
	}
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
	mu   sync.Mutex
	code int
}

func (r *responseBodyWriter) WriteString(s string) (n int, err error) {
	r.body.WriteString(s)
	return r.ResponseWriter.WriteString(s)
}
