package tracer

import (
	"bytes"
	"context"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func TestStart_CreatesSpan(t *testing.T) {
	ctx := context.Background()
	ctx, span := Start(ctx, "test-span")

	if span == nil {
		t.Fatal("expected span, got nil")
	}

	if span.name != "test-span" {
		t.Errorf("expected name test-span, got %v", span.name)
	}

	if span.traceID == "" {
		t.Error("expected traceID to be set")
	}

	// verify it's stored in context
	got := ctx.Value(spanKey)
	if got == nil {
		t.Fatal("expect span in context")
	}

	if got != span {
		t.Error("context span mismatch")
	}
}

func TestStart_ReusesParentTraceID(t *testing.T) {
	ctx, parent := Start(context.Background(), "parent")
	_, child := Start(ctx, "child")

	if child.traceID != parent.traceID {
		t.Errorf("expected same traceID, got %v and %v", parent.traceID, child.traceID)
	}
}

func TestSpan_End_Logs(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	span := &Span{
		traceID: "abc",
		name: "test",
		start: time.Now().Add(-50 * time.Millisecond),
	}

	span.End()

	out := buf.String()
	if !strings.Contains(out, "trace=abc") {
		t.Error("missing traceID in log")
	}

	if !strings.Contains(out, "span=test") {
		t.Error("missing span name in log")
	}

	if !strings.Contains(out, "duration=") {
		t.Error("missing duration in log")
	}

	if !strings.Contains(out, "alloc_bytes") {
		t.Error("missing alloc bytes in log")
	}
}
