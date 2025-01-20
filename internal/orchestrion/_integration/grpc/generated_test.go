// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-present Datadog, Inc.
//
// Code generated by 'go generate'; DO NOT EDIT.

package grpc

import (
	"testing"

	"github.com/DataDog/dd-trace-go/internal/orchestrion/_integration/internal/harness"
)

func Test(t *testing.T) {
	harness.Run(t, new(TestCase))
}
