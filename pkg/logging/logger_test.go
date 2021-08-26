// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logging_test

import (
	"bytes"
	"context"
	"net/url"
	"testing"

	"github.com/jusongchen/REST-app/pkg/logging"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// MemorySink implements zap.Sink by writing all messages to a buffer.
type MemorySink struct {
	*bytes.Buffer
}

// Implement Close and Sync as no-ops to satisfy the interface. The Write
// method is provided by the embedded buffer.
func (s *MemorySink) Close() error { return nil }
func (s *MemorySink) Sync() error  { return nil }

func TestNewLogger(t *testing.T) {
	t.Parallel()

	logger := logging.NewLogger(true)
	require.NotNil(t, logger)

}

func TestSetLogLevel(t *testing.T) {
	t.Parallel()

	// Create a sink instance, and register it with zap for the "memory" protocol.
	sink := &MemorySink{new(bytes.Buffer)}
	zap.RegisterSink("memory", func(*url.URL) (zap.Sink, error) {
		return sink, nil
	})

	// Test default (INFO) level logging
	defaultLogger := logging.NewLogger(false, "memory://")
	require.Equal(t, logging.GetDefaultLogLevel(), zapcore.InfoLevel)

	defaultLogger.Debug("debug")
	require.Empty(t, sink.String())
	sink.Reset()

	defaultLogger.Info("info")
	require.Containsf(t, sink.String(), `"message":"info"`, "Should log message with contents info")
	sink.Reset()

	defaultLogger.Warn("warn")
	require.Containsf(t, sink.String(), `"message":"warn"`, "Should log message with contents warn")
	sink.Reset()

	defaultLogger.Error("error")
	require.Containsf(t, sink.String(), `"message":"error"`, "Should log message with contents info")
	sink.Reset()

	// Test DEBUG level logging
	logging.SetDefaultLevel("DEBUG")
	require.Equal(t, logging.GetDefaultLogLevel(), zapcore.DebugLevel)
	debugLogger := logging.NewLogger(false, "memory://")

	debugLogger.Debug("debug")
	require.Containsf(t, sink.String(), `"message":"debug"`, "Should log message with contents 'debug'")
	sink.Reset()

	debugLogger.Info("info")
	require.Containsf(t, sink.String(), `"message":"info"`, "Should log message with contents 'info'")
	sink.Reset()

	debugLogger.Warn("warn")
	require.Containsf(t, sink.String(), `"message":"warn"`, "Should log message with contents 'warn'")
	sink.Reset()

	debugLogger.Error("error")
	require.Containsf(t, sink.String(), `"message":"error"`, "Should log message with contents 'error'")
	sink.Reset()

	// Test INFO level logging
	logging.SetDefaultLevel("INFO")
	require.Equal(t, logging.GetDefaultLogLevel(), zapcore.InfoLevel)
	infoLogger := logging.NewLogger(false, "memory://")

	infoLogger.Debug("debug")
	require.Empty(t, sink.String())
	sink.Reset()

	infoLogger.Info("info")
	require.Containsf(t, sink.String(), `"message":"info"`, "Should log message with contents 'info'")
	sink.Reset()

	infoLogger.Warn("warn")
	require.Containsf(t, sink.String(), `"message":"warn"`, "Should log message with contents 'warn'")
	sink.Reset()

	infoLogger.Error("error")
	require.Containsf(t, sink.String(), `"message":"error"`, "Should log message with contents 'error'")
	sink.Reset()

	// Test INFO level logging
	logging.SetDefaultLevel("WARN")
	require.Equal(t, logging.GetDefaultLogLevel(), zapcore.WarnLevel)
	warnLogger := logging.NewLogger(false, "memory://")

	warnLogger.Debug("debug")
	require.Empty(t, sink.String())
	sink.Reset()

	warnLogger.Info("info")
	require.Empty(t, sink.String())
	sink.Reset()

	warnLogger.Warn("warn")
	require.Containsf(t, sink.String(), `"message":"warn"`, "Should log message with contents 'warn'")
	sink.Reset()

	warnLogger.Error("error")
	require.Containsf(t, sink.String(), `"message":"error"`, "Should log message with contents 'error'")
	sink.Reset()

	// Test INFO level logging
	logging.SetDefaultLevel("ERROR")
	require.Equal(t, logging.GetDefaultLogLevel(), zapcore.ErrorLevel)
	errorLogger := logging.NewLogger(false, "memory://")

	errorLogger.Debug("debug")
	require.Empty(t, sink.String())
	sink.Reset()

	errorLogger.Info("info")
	require.Empty(t, sink.String())
	sink.Reset()

	errorLogger.Warn("warn")
	require.Empty(t, sink.String())
	sink.Reset()

	errorLogger.Error("error")
	require.Containsf(t, sink.String(), `"message":"error"`, "Should log message with contents 'error'")
	sink.Reset()
}

func TestDefaultLogger(t *testing.T) {
	t.Parallel()

	logger1 := logging.DefaultLogger()
	if logger1 == nil {
		t.Fatal("expected logger to never be nil")
	}

	logger2 := logging.DefaultLogger()
	if logger2 == nil {
		t.Fatal("expected logger to never be nil")
	}

	// Intentionally comparing identities here
	if logger1 != logger2 {
		t.Errorf("expected %#v to be %#v", logger1, logger2)
	}
}

func TestContext(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger1 := logging.FromContext(ctx)
	if logger1 == nil {
		t.Fatal("expected logger to never be nil")
	}

	ctx = logging.WithLogger(ctx, logger1)

	logger2 := logging.FromContext(ctx)
	if logger1 != logger2 {
		t.Errorf("expected %#v to be %#v", logger1, logger2)
	}
}
