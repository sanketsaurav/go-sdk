package sentry

import (
	"context"
	"net/http"
	"runtime"
	"time"

	raven "github.com/getsentry/sentry-go"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
)

var (
	_ Sender = (*Client)(nil)
)

// MustNew returns a new client and panics on error.
func MustNew(cfg Config) *Client {
	c, err := New(cfg)
	if err != nil {
		panic(err)
	}
	return c
}

// New returns a new client.
func New(cfg Config) (*Client, error) {
	return &Client{
		Config: cfg,
		xport:  &http.Transport{},
	}, nil
}

// Client is a wrapper for the sentry-go client.
type Client struct {
	Config Config
	xport  *http.Transport
}

// Notify sends a notification.
func (c Client) Notify(ctx context.Context, ee logger.ErrorEvent) error {
	c.captureEvent(errEvent(ctx, ee))
	return c.flush(ctx)
}

func (c Client) captureEvent(event *Event) {

}

func (c Client) flush(ctx context.Context) error {

	return nil
}

func errEvent(ctx context.Context, ee logger.ErrorEvent) *raven.Event {
	return &Event{
		Timestamp:   logger.GetEventTimestamp(ctx, ee).Unix(),
		Fingerprint: errFingerprint(ctx, ex.ErrClass(ee.Err).Error()),
		Level:       Level(ee.GetFlag()),
		Tags:        errTags(ctx),
		Extra:       errExtra(ctx),
		Platform:    "go",
		Sdk: SdkInfo{
			Name:    SDK,
			Version: raven.Version,
			Packages: []raven.SdkPackage{{
				Name:    SDK,
				Version: raven.Version,
			}},
		},
		Request: errRequest(ee),
		Message: ex.ErrClass(ee.Err).Error(),
		Exception: []Exception{
			{
				Type:       ex.ErrClass(ee.Err).Error(),
				Value:      ex.ErrMessage(ee.Err),
				Stacktrace: errStackTrace(ee.Err),
			},
		},
	}
}

func errFingerprint(ctx context.Context, extra ...string) []string {
	if fingerprint := GetFingerprint(ctx); fingerprint != nil {
		return fingerprint
	}
	return append(logger.GetPath(ctx), extra...)
}

func errTags(ctx context.Context) map[string]string {
	return logger.GetLabels(ctx)
}

func errExtra(ctx context.Context) map[string]interface{} {
	return logger.GetAnnotations(ctx)
}

func errRequest(ee logger.ErrorEvent) (requestMeta Request) {
	if ee.State == nil {
		return
	}
	typed, ok := ee.State.(*http.Request)
	if !ok {
		return
	}
	requestMeta = requestMeta.FromHTTPRequest(typed)
	return
}

func errStackTrace(err error) *Stacktrace {
	if err != nil {
		return &Stacktrace{Frames: errFrames(err)}
	}
	return nil
}

func errFrames(err error) []Frame {
	stacktrace := ex.ErrStackTrace(err)
	if stacktrace == nil {
		return []Frame{}
	}
	pointers, ok := stacktrace.(ex.StackPointers)
	if !ok {
		return []Frame{}
	}

	var output []Frame
	runtimeFrames := runtime.CallersFrames(pointers)

	for {
		callerFrame, more := runtimeFrames.Next()
		output = append([]Frame{
			NewFrame(callerFrame),
		}, output...)
		if !more {
			break
		}
	}

	return output
}
