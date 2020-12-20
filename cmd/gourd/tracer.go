package main

import (
	"context"
	"database/sql/driver"
	"errors"

	"github.com/luna-duclos/instrumentedsql"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
)

type tracer struct {
	traceOrphans bool
}

type span struct {
	tracer
	parent trace.Span
}

// NewTracer returns a tracer that will fetch trace using opentelemetry
// if traceOrphans is set to true, then spans with no parent will be traced anyway, if false, they will not be.
func NewTracer(traceOrphans bool) instrumentedsql.Tracer { return tracer{traceOrphans: traceOrphans} }

func (t tracer) GetSpan(ctx context.Context) instrumentedsql.Span {
	if ctx == nil {
		return span{parent: nil, tracer: t}
	}

	return span{parent: trace.SpanFromContext(ctx), tracer: t}
}

func (s span) NewChild(name string) instrumentedsql.Span {
	if s.parent == nil && !s.traceOrphans {
		return s
	}

	ctx := context.Background()
	if s.parent != nil {
		ctx = trace.ContextWithSpan(context.Background(), s.parent)
	}
	_, parent := otel.Tracer("db").Start(ctx, name)

	return span{parent: parent, tracer: s.tracer}
}

func (s span) SetLabel(k, v string) {
	s.parent.SetAttributes(label.Key(k).String(v))
}

func (s span) SetError(err error) {
	if err == nil || errors.Is(err, driver.ErrSkip) {
		return
	}
	s.parent.RecordError(err)
	s.SetLabel("error", "true")
}

func (s span) Finish() {
	s.parent.End()
}
