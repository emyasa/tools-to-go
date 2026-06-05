// Package tracer handles observability related code
package tracer

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/google/uuid"
)

type Span struct {
	traceID string
	name string
	start time.Time
	startAlloc uint64
}

func Start(ctx context.Context, name string) (context.Context, *Span) {
	parent := extractParentSpan(ctx)

	traceID := uuid.NewString()
	if parent != nil {
		traceID = parent.traceID
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	span := &Span{
		traceID: traceID,
		name: name,
		start: time.Now(),
		startAlloc: m.TotalAlloc,
	}

	ctx = context.WithValue(ctx, spanKey, span)

	return ctx, span
}

func (s *Span) End() {
	elapsed := time.Since(s.start)

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	log.Printf("trace=%s span=%s duration=%s alloc_bytes=%d",
		s.traceID,
		s.name,
		elapsed,
		m.TotalAlloc - s.startAlloc)
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

func formatBytes(b uint64) string {
	const unit = 1024

	if b < unit {
		return fmt.Sprintf("%d B", b)
	}

	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %ciB",
		float64(b)/float64(div),
		"KMGTPE"[exp],
	)
}
