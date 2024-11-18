// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package sizeprocessor // import "github.com/multiplayer-app/opentelemetry-collector-contrib/processor/sizeprocessor"

import (
	"context"
	"encoding/json"

	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

type spanAttributesProcessor struct {
	logger *zap.Logger
}

// newTracesProcessor returns a processor that modifies attributes of a span.
// To construct the attributes processors, the use of the factory methods are required
// in order to validate the inputs.
func newSpanAttributesProcessor(logger *zap.Logger) *spanAttributesProcessor {
	return &spanAttributesProcessor{
		logger: logger,
	}
}

func calculateSpanSize(span ptrace.Span) (int, error) {
	data, err := json.Marshal(span)
	if err != nil {
		return 0, err
	}
	return len(data), nil
}

func (a *spanAttributesProcessor) processTraces(ctx context.Context, td ptrace.Traces) (ptrace.Traces, error) {
	rss := td.ResourceSpans()
	for i := 0; i < rss.Len(); i++ {
		rs := rss.At(i)
		ilss := rs.ScopeSpans()
		for j := 0; j < ilss.Len(); j++ {
			ils := ilss.At(j)
			spans := ils.Spans()
			// scope := ils.Scope()
			for k := 0; k < spans.Len(); k++ {
				span := spans.At(k)
				size, err := calculateSpanSize(span)
				if err != nil {
					continue
				}

				span.Attributes().PutInt("span.size", int64(size))
				// if a.skipExpr != nil {
				// 	skip, err := a.skipExpr.Eval(ctx, ottlspan.NewTransformContext(span, scope, resource, ils, rs))
				// 	if err != nil {
				// 		return td, err
				// 	}
				// 	if skip {
				// 		continue
				// 	}
				// }
				// a.attrProc.Process(ctx, a.logger, span.Attributes())
			}
		}
	}
	return td, nil
}
