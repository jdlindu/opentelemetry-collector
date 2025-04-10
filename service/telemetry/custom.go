// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package telemetry // import "go.opentelemetry.io/collector/service/telemetry"

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"runtime"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
)

// RotateConfig TODO
type RotateConfig struct {
	// Enabled controls whether or not rotate logs
	Enabled bool `mapstructure:"enabled"`
	// MaxMegabytes is the maximum size in megabytes of the file before it gets
	// rotated. It defaults to 100 megabytes.
	MaxMegabytes int `mapstructure:"max_megabytes"`
	// MaxDays is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxDays int `mapstructure:"max_days"`
	// MaxBackups is the maximum number of old log files to retain. The default
	// is to 100 files.
	MaxBackups int `mapstructure:"max_backups"`
	// LocalTime determines if the time used for formatting the timestamps in
	// backup files is the computer's local time.  The default is to use UTC
	// time.
	LocalTime bool `mapstructure:"localtime"`
}

// NewDefaultRotateConfig TODO
func NewDefaultRotateConfig() *RotateConfig {
	return &RotateConfig{
		Enabled:      true,
		MaxMegabytes: 100,
		MaxDays:      0,
		MaxBackups:   10,
		LocalTime:    false,
	}
}

// NewWriter TODO
func (cfg *RotateConfig) NewWriter(filename string) (io.WriteCloser, error) {
	if cfg.MaxMegabytes < 0 {
		return nil, fmt.Errorf("invalid MaxMegabytes %d", cfg.MaxMegabytes)
	}
	if cfg.MaxDays < 0 {
		return nil, fmt.Errorf("invalid MaxDays %d", cfg.MaxDays)
	}
	// #nosec G302 G304 -- filename is a trusted safe path, and should allow to be read by other users
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	if !cfg.Enabled {
		return file, nil
	}
	return cfg.newLumberjackWriter(filename), nil
}
func (cfg *RotateConfig) newLumberjackWriter(filename string) io.WriteCloser {
	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    cfg.MaxMegabytes,
		MaxAge:     cfg.MaxDays,
		MaxBackups: cfg.MaxBackups,
		LocalTime:  cfg.LocalTime,
	}
}

func getRotationSinkFactory(cfg *RotateConfig) func(u *url.URL) (zap.Sink, error) {
	return func(u *url.URL) (zap.Sink, error) {
		p := u.Query().Get("path")
		writer, err := cfg.NewWriter(p)
		if err != nil {
			return nil, err
		}
		return nopSyncSink{writer}, nil
	}
}

// lumberjack.Logger does not provide a Sync() method, which is required by zap.Sink
// explanation: https://github.com/natefinch/lumberjack/pull/47#issuecomment-322502210
type nopSyncSink struct {
	io.WriteCloser
}

func (w nopSyncSink) Sync() error {
	return nil
}
func setRotationURL(paths []string, rotationSchema string) ([]string, error) {
	res := make([]string, 0, len(paths))
	for _, p := range paths {
		if runtime.GOOS == "windows" && filepath.IsAbs(p) {
			res = append(res, rotationSchema+":?path="+url.QueryEscape(p))
			continue
		}
		u, err := url.Parse(p)
		if err != nil {
			return nil, err
		}
		if (u.Scheme == "" || u.Scheme == "file") &&
			u.Path != "stdout" && u.Path != "stderr" {
			// Copied from zap. Only clean URLs are allowed
			if u.User != nil {
				return nil, fmt.Errorf("user and password not allowed with file URLs: got %v", u)
			}
			if u.Fragment != "" {
				return nil, fmt.Errorf("fragments not allowed with file URLs: got %v", u)
			}
			if u.RawQuery != "" {
				return nil, fmt.Errorf("query parameters not allowed with file URLs: got %v", u)
			}
			// Error messages are better if we check hostname and port separately.
			if u.Port() != "" {
				return nil, fmt.Errorf("ports not allowed with file URLs: got %v", u)
			}
			if hn := u.Hostname(); hn != "" && hn != "localhost" {
				return nil, fmt.Errorf("file URLs must leave host empty or use localhost: got %v", u)
			}
			res = append(res, rotationSchema+":?path="+url.QueryEscape(u.Path))
			continue
		}
		res = append(res, p)
	}
	return res, nil
}
