// Package tracer handles observability related code
package tracer

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
)

type Span struct {
	traceID string
	name string
	start time.Time
}

func Start(ctx context.Context, name string) (context.Context, *Span) {
	parent := extractParentSpan(ctx)

	traceID := uuid.NewString()
	if parent != nil {
		traceID = parent.traceID
	}

	span := &Span{
		traceID: traceID,
		name: name,
		start: time.Now(),
	}

	ctx = context.WithValue(ctx, spanKey, span)

	return ctx, span
}

func (s *Span) End() {
	elapsed := time.Since(s.start)

	log.Printf("trace=%s span=%s duration=%s",
		s.traceID,
		s.name,
		elapsed)
}

type spanKeyType struct {}
var spanKey = spanKeyType{}

func extractParentSpan(ctx context.Context) *Span {
	val := ctx.Value(spanKey)
	if val == nil {
		return nil
	}

	return val.(*Span)
}
