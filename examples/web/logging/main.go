package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"
	"github.com/blend/go-sdk/webutil"
)

func main() {
	bigPayload := make([]byte, 1<<20) // 1mb
	rand.Read(bigPayload)

	log := logger.Prod()
	app := web.MustNew(web.OptLog(log))
	app.GET("/", func(r *web.Ctx) web.Result {
		r.Response.Header().Set(webutil.HeaderContentType, webutil.ContentTypeText)

		for x := 0; x < 100; x++ {
			select {
			case <-r.Context().Done():
				return nil
			case <-time.After(100 * time.Millisecond):
				if _, err := io.Copy(r.Response, bytes.NewReader(bigPayload)); err != nil {
					return web.Text.InternalError(err)
				}
				r.Response.(http.Flusher).Flush()
			}
		}
		return nil
	})
	log.Listen(webutil.HTTPRequest, logger.DefaultListenerName, webutil.NewHTTPRequestEventListener(func(_ context.Context, wre webutil.HTTPRequestEvent) {
		log.Infof("got a new request at route: %s", wre.Route)
	}))

	errorFilters := ex.Filter(
		ex.FilterIs(syscall.EPIPE),
		ex.FilterIs(os.ErrNotExist),
		webutil.IsNetOpError(),
	)
	filteredErrorListener := logger.NewErrorEventListener(func(_ context.Context, ee logger.ErrorEvent) {
		if errorFilters.Any(ee.Err) {
			fmt.Fprintf(os.Stdout, "filtered err: %+v", ee.Err)
			return
		}
		fmt.Fprintf(os.Stderr, "important err: %+v %#v", ee.Err, ee.Err)
	})
	log.Listen(logger.Error, logger.DefaultListenerName, filteredErrorListener)
	log.Listen(logger.Fatal, logger.DefaultListenerName, filteredErrorListener)

	graceful.Shutdown(app)
}
