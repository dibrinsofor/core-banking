package middlewares

import (
	"encoding/json"
	"log"
	"net/http"

	redisstore "github.com/dibrinsofor/core-banking/internal/redis"
	"github.com/gin-gonic/gin"
)

type header struct {
	IdempotencyKey string `header:"Idempotency-Key" binding:"required"`
}

func Idempotency() gin.HandlerFunc {
	return func(c *gin.Context) {
		var h header

		tw := &timeoutWriter{ResponseWriter: c.Writer, h: make(http.Header)}
		c.Writer = tw

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
		if user.RespCode == http.StatusBadRequest {
			return
		}

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
			// c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			// 	"message": "unable to complete request",
			// })

			tw.ResponseWriter.WriteHeader(http.StatusInternalServerError)
			eResp, _ := json.Marshal(gin.H{"message": "unable to complete request"})
			tw.ResponseWriter.Write(eResp)
		case <-finished:
			// if finished, set headers and write resp
			tw.mu.Lock()
			defer tw.mu.Unlock()

			// todo: persist idempkey to redis and result

			// map Headers from tw.Header() (written to by gin)
			// to tw.ResponseWriter for response
			dst := tw.ResponseWriter.Header()
			for k, vv := range tw.Header() {
				dst[k] = vv
			}
			tw.ResponseWriter.WriteHeader(tw.code)
			// tw.wbuf will have been written to already when gin writes to tw.Write()
			tw.ResponseWriter.Write(tw.wbuf.Bytes())
		}
	}
}
