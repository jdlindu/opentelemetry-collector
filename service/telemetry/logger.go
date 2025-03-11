// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package telemetry // import "go.opentelemetry.io/collector/service/telemetry"

import (
	"github.com/google/uuid"
	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// newLogger creates a Logger and a LoggerProvider from Config.
func newLogger(set Settings, cfg Config) (*zap.Logger, log.LoggerProvider, error) {
	// Copied from NewProductionConfig.
	ec := zap.NewProductionEncoderConfig()
	ec.EncodeTime = zapcore.ISO8601TimeEncoder
	zapCfg := &zap.Config{
		Level:             zap.NewAtomicLevelAt(cfg.Logs.Level),
		Development:       cfg.Logs.Development,
		Encoding:          cfg.Logs.Encoding,
		EncoderConfig:     ec,
		OutputPaths:       cfg.Logs.OutputPaths,
		ErrorOutputPaths:  cfg.Logs.ErrorOutputPaths,
		DisableCaller:     cfg.Logs.DisableCaller,
		DisableStacktrace: cfg.Logs.DisableStacktrace,
		InitialFields:     cfg.Logs.InitialFields,
	}

	if zapCfg.Encoding == "console" {
		// Human-readable timestamps for console format of logs.
		zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	rotationSchema := "rotation-" + uuid.NewString()
	err := zap.RegisterSink(rotationSchema, getRotationSinkFactory(NewDefaultRotateConfig()))
	if err != nil {
		return nil, nil, err
	}
	zapCfg.OutputPaths, err = setRotationURL(zapCfg.OutputPaths, rotationSchema)
	if err != nil {
		return nil, nil, err
	}
	zapCfg.ErrorOutputPaths, err = setRotationURL(zapCfg.ErrorOutputPaths, rotationSchema)
	if err != nil {
		return nil, nil, err
	}

	logger, err := zapCfg.Build(set.ZapOptions...)
	if err != nil {
		return nil, nil, err
	}

	var lp log.LoggerProvider

	if len(cfg.Logs.Processors) > 0 && set.SDK != nil {
		lp = set.SDK.LoggerProvider()

		logger = logger.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
			core, err := zapcore.NewIncreaseLevelCore(zapcore.NewTee(
				c,
				otelzap.NewCore("go.opentelemetry.io/collector/service/telemetry",
					otelzap.WithLoggerProvider(lp),
				),
			), zap.NewAtomicLevelAt(cfg.Logs.Level))
			if err != nil {
				panic(err)
			}
			return core
		}))
	}

	if cfg.Logs.Sampling != nil && cfg.Logs.Sampling.Enabled {
		logger = newSampledLogger(logger, cfg.Logs.Sampling)
	}

	return logger, lp, nil
}

func newSampledLogger(logger *zap.Logger, sc *LogsSamplingConfig) *zap.Logger {
	// Create a logger that samples every Nth message after the first M messages every S seconds
	// where N = sc.Thereafter, M = sc.Initial, S = sc.Tick.
	opts := zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewSamplerWithOptions(
			core,
			sc.Tick,
			sc.Initial,
			sc.Thereafter,
		)
	})
	return logger.WithOptions(opts)
}
