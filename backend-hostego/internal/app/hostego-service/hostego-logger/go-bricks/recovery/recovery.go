// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package Recovery

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net"
	"net/http/httputil"
	"os"
	"strings"
	"time"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// RecoveryFunc defines the function passable to CustomRecovery.
type RecoveryFunc func(c *gin.Context, err interface{})

//CustomRecovery returns a middleware that recovers from any panics and calls the provided handle func to handle it.
func CustomRecovery(handle RecoveryFunc) gin.HandlerFunc {
	return recoveryWithWriter(gin.DefaultErrorWriter, handle)
}

// RecoveryWithWriter returns a middleware for a given writer that recovers from any panics and writes a 500 if there was one.
func recoveryWithWriter(out io.Writer, recovery ...RecoveryFunc) gin.HandlerFunc {
	return customRecoveryWithWriter(out, recovery[0])
}

// CustomRecoveryWithWriter returns a middleware for a given writer that recovers from any panics and calls the provided handle func to handle it.
func customRecoveryWithWriter(out io.Writer, handle RecoveryFunc) gin.HandlerFunc {
	var logger *log.Logger
	if out != nil {
		logger = log.New(out, "\n\n\x1b[31m", log.LstdFlags)
	}
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				if logger != nil {
					stack := stack(3)
					httpRequest, _ := httputil.DumpRequest(c.Request, false)
					headers := strings.Split(string(httpRequest), "\r\n")
					for idx, header := range headers {
						current := strings.Split(header, ":")
						if current[0] == "Authorization" {
							headers[idx] = current[0] + ": *"
						}
					}
					headersToStr := strings.Join(headers, "\r\n")
					if brokenPipe {
						logger.Printf("%s\n%s%s", err, headersToStr)
					} else if gin.IsDebugging() {
						logger.Printf("[Recovery] %s panic recovered:\n%s\n%s\n%s%s",
							timeFormat(time.Now()), headersToStr, err, stack)
					} else {
						logger.Printf("[Recovery] %s panic recovered:\n%s\n%s",
							timeFormat(time.Now()), err, stack)
					}
				}
				if brokenPipe {
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
				} else {
					handle(c, err)
				}
			}
		}()
		c.Next()
	}
}
