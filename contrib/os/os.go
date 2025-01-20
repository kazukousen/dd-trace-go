// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2024 Datadog, Inc.

// Package os provides integrations into the standard library's `os` package,
// allowing protection against Local File Inclusion (LFI) attacks.
package os

// These imports satisfy injected dependencies for Orchestrion auto instrumentation.
import (
	"context"
	"os"

	"gopkg.in/DataDog/dd-trace-go.v1/appsec/events"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/internal/appsec/dyngo"
	"gopkg.in/DataDog/dd-trace-go.v1/internal/appsec/emitter/ossec"
	"gopkg.in/DataDog/dd-trace-go.v1/internal/telemetry"
)

const componentName = "os"

func init() {
	telemetry.LoadIntegration(componentName)
	tracer.MarkIntegrationImported("os")
}

// OpenFile is a [context.Context]-aware version of [os.OpenFile], that allows
// the use of ASM rules to protect against Local File Inclusion (LFI) attacks.
func OpenFile(ctx context.Context, path string, flag int, perm os.FileMode) (file *os.File, err error) {
	parent, _ := dyngo.FromContext(ctx)
	if parent != nil {
		op := &ossec.OpenOperation{
			Operation: dyngo.NewOperation(parent),
		}

		var block bool
		dyngo.OnData(op, func(*events.BlockingSecurityEvent) {
			block = true
		})

		dyngo.StartOperation(op, ossec.OpenOperationArgs{
			Path:  path,
			Flags: flag,
			Perms: perm,
		})

		defer dyngo.FinishOperation(op, ossec.OpenOperationRes[*os.File]{
			File: &file,
			Err:  &err,
		})

		if block {
			return
		}
	}

	return os.OpenFile(path, flag, perm)
}
