// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dmltest

import (
	"bytes"
	"context"
	"os/exec"
	"path/filepath"

	"github.com/corestoreio/errors"
	"github.com/go-sql-driver/mysql"
)

type SQLDumpOptions struct {
	MySQLPath string
	MySQLArgs []string
	// mocked out for testing.
	execCommandContext func(ctx context.Context, name string, arg ...string) error
}

// SQLDumpLoad reads all files regonized by `globPattern` argument into
// MySQL/MariaDB. The password will be visible via a process manager.
func SQLDumpLoad(ctx context.Context, dsn string, globPattern string, o SQLDumpOptions) error {

	matches, err := filepath.Glob(globPattern)
	if err != nil {
		return errors.WithStack(err)
	}

	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		return errors.WithStack(err)
	}

	mybin := "mysql"
	var args []string
	if cfg.Addr != "" {
		args = append(args, "--host", cfg.Addr)
	}
	if cfg.User != "" {
		args = append(args, "--user", cfg.User)
	}
	if cfg.Passwd != "" {
		args = append(args, "--password", cfg.Passwd)
	}
	args = append(args, "--database", cfg.DBName)
	args = append(args, o.MySQLArgs...)
	args = append(args, "<", "<placeholder>")

	execCmd := o.execCommandContext
	if execCmd == nil {
		execCmd = realExecCmd
	}

	for _, file := range matches {
		args[len(args)-1] = file
		if err := execCmd(ctx, mybin, args...); err != nil {
			return errors.Wrapf(err, "[dmltest] Failed to load SQL dump with file %q", file)
		}
	}

	return nil
}

func realExecCmd(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)

	var bufe bytes.Buffer
	var bufo bytes.Buffer
	cmd.Stderr = &bufe
	cmd.Stdout = &bufo

	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "[dmltest] SQLDumpLoad\nStderr: %s\nStdout: %s", bufe.String(), bufo.String())
	}

	return nil
}
