package qtracer

import (
	"context"
	"net/http"

	"github.com/newrelic/go-agent/v3/newrelic"
)

var (
	app *newrelic.Application
)

type Span struct {
	nrSegment Segment
	ctx       context.Context
	md        http.Header
}

type Segment interface {
	End()
	AddAttribute(key string, val interface{})
}

// GetApp returns app variable
func GetApp() *newrelic.Application {
	return app
}

func (s *Span) Finish() {
	if s.nrSegment != nil {
		s.nrSegment.End()
	}
}

func InitTracer(data *newrelic.Application) {
	app = data
}

func StartGoroutineSpanFromContext(ctx context.Context, name string) (Span, context.Context) {
	var span Span
	span.nrSegment, ctx = StartGoroutineSegment(ctx, "[trace] "+name)
	span.ctx = ctx
	span.md = GetMetadataFromContext(ctx)

	return span, ctx
}

func StartGoroutineSegment(ctx context.Context, name string) (*newrelic.Segment, context.Context) {
	if app == nil {
		return nil, nil
	}

	txn := newrelic.FromContext(ctx)

	newTxn := txn.NewGoroutine()
	seg := newTxn.StartSegment(name)

	ctx = newrelic.NewContext(ctx, newTxn)

	return seg, ctx
}

func StartSpanFromContext(ctx context.Context, name string) (Span, context.Context) {
	var span Span

	seg := StartSegment(ctx, "[trace] "+name)
	span.nrSegment = seg
	span.md = GetMetadataFromContext(ctx)

	return span, ctx
}

func StartSegment(ctx context.Context, name string) *newrelic.Segment {
	if app == nil {
		return nil
	}

	txn := newrelic.FromContext(ctx)
	return txn.StartSegment(name)
}

func GetMetadataFromContext(ctx context.Context) http.Header {
	hdr := http.Header{}
	if app == nil {
		return hdr
	}
	txn := newrelic.FromContext(ctx)
	if txn == nil {
		return hdr
	}
	txn.InsertDistributedTraceHeaders(hdr)
	return hdr
}
