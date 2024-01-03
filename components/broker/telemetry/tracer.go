// Copyright 2023 XigXog
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package telemetry

import (
	"context"
	"fmt"
	"time"

	"github.com/xigxog/kubefox/api"
	"github.com/xigxog/kubefox/core"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	otelTracer = otel.Tracer("")
)

type Span interface {
	End(*core.Event)
}

type span struct {
	cancel   context.CancelFunc
	otelSpan trace.Span
	req      *core.Event
}

func NewSpan(ctx context.Context, timeout time.Duration, req *core.Event) (context.Context, Span) {
	ctx, cancel := context.WithTimeout(ctx, timeout)

	typ := api.EventTypeUnknown
	if req.Type != "" {
		typ = api.EventType(req.Type)
	}

	trcId, _ := trace.TraceIDFromHex(req.TraceId())
	spnId, _ := trace.SpanIDFromHex(req.SpanId())
	trcFlags := trace.TraceFlags(req.TraceFlags())

	ctx = trace.ContextWithRemoteSpanContext(ctx, trace.NewSpanContext(
		trace.SpanContextConfig{
			TraceID:    trcId,
			SpanID:     spnId,
			TraceFlags: trcFlags,
		}))

	ctx, otelSpan := otelTracer.Start(ctx, fmt.Sprintf("%s event", typ),
		trace.WithAttributes(traceAttrs(req)...),
	)
	req.SetTraceId(otelSpan.SpanContext().TraceID().String())
	req.SetSpanId(otelSpan.SpanContext().SpanID().String())
	req.SetTraceFlags(byte(otelSpan.SpanContext().TraceFlags()))

	return ctx, &span{
		cancel:   cancel,
		otelSpan: otelSpan,
		req:      req,
	}
}

func (sp *span) End(resp *core.Event) {
	if resp != nil {
		sp.otelSpan.SetAttributes(traceAttrs(resp)...)
		resp.SetTraceId(sp.otelSpan.SpanContext().TraceID().String())
		resp.SetSpanId(sp.otelSpan.SpanContext().SpanID().String())
		resp.SetTraceFlags(byte(sp.otelSpan.SpanContext().TraceFlags()))

		// TODO decide on how to deal with errors
		// if err := resp.GetError(); err != nil {
		// 	sp.otelSpan.RecordError(err)
		// 	sp.otelSpan.SetStatus(codes.Error, err.Error())
		// }
	}

	sp.otelSpan.End()
	sp.cancel()
}

func traceAttrs(req *core.Event) []attribute.KeyValue {
	attrs := []attribute.KeyValue{}

	if req != nil && req.Type != "" {
		attrs = append(attrs, attribute.Key("core.event.type").String(req.Type))
	}
	if req != nil && req.Id != "" {
		attrs = append(attrs, attribute.Key("core.event.id").String(req.Id))
	}
	if req != nil && req.ParentId != "" {
		attrs = append(attrs, attribute.Key("core.event.parent-id").String(req.ParentId))
	}
	if req != nil && req.Source != nil {
		attrs = append(attrs, attribute.Key("core.event.source").String(req.Source.Key()))
	}
	if req != nil && req.Target != nil {
		attrs = append(attrs, attribute.Key("core.event.target").String(req.Source.Key()))
	}

	return attrs
}
