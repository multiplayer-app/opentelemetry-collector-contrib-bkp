// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package sizeprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/sizeprocessor"

import (
	"context"
	"encoding/json"

	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

type logAttributesProcessor struct {
	logger *zap.Logger
}

// newLogAttributesProcessor returns a processor that modifies attributes of a
// log record. To construct the attributes processors, the use of the factory
// methods are required in order to validate the inputs.
func newLogAttributesProcessor(logger *zap.Logger) *logAttributesProcessor {
	return &logAttributesProcessor{
		logger: logger,
	}
}

func calculateLogSize(lr plog.LogRecord) (int, error) {
	data, err := json.Marshal(lr)
	if err != nil {
		return 0, err
	}
	return len(data), nil
}

func (a *logAttributesProcessor) processLogs(ctx context.Context, ld plog.Logs) (plog.Logs, error) {
	rls := ld.ResourceLogs()
	for i := 0; i < rls.Len(); i++ {
		rs := rls.At(i)
		ilss := rs.ScopeLogs()
		// resource := rs.Resource()
		for j := 0; j < ilss.Len(); j++ {
			ils := ilss.At(j)
			logs := ils.LogRecords()
			// library := ils.Scope()
			for k := 0; k < logs.Len(); k++ {
				lr := logs.At(k)

				size, err := calculateLogSize(lr)
				if err != nil {
					continue
				}

				lr.Attributes().PutInt("log.size", int64(size))
			}
		}
	}

	return ld, nil
}
