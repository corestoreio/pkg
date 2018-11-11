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
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/go-sql-driver/mysql"
)

func writeMySQLDefaults(cfg *mysql.Config, o *SQLDumpOptions) (string, error) {

	var myCnf bytes.Buffer
	myCnf.WriteString("[mysql]\n")

	if cfg.Addr != "" {

		host, port, err := net.SplitHostPort(cfg.Addr)
		if err != nil {
			return "", errors.WithStack(err)
		}
		fmt.Fprintf(&myCnf, "host = %s\n", host)

		if port != "" {
			fmt.Fprintf(&myCnf, "port = %s\n", port)
		}
	}
	if cfg.User != "" {
		fmt.Fprintf(&myCnf, "user = %s\n", cfg.User)
	}
	if cfg.Passwd != "" {
		fmt.Fprintf(&myCnf, "password = %s\n", cfg.Passwd)
	}
	for _, ma := range o.MySQLArgs {
		myCnf.WriteString(ma)
		myCnf.WriteByte('\n')
	}
	fmt.Fprintf(&myCnf, "database = %s\n", cfg.DBName)

	df, err := ioutil.TempFile("", "mydefaults-")
	if err != nil {
		return "", errors.WithStack(err)
	}
	defer df.Close()

	if _, err := df.Write(myCnf.Bytes()); err != nil {
		return "", errors.WithStack(err)
	}
	return df.Name(), nil
}

func stdInExec(ctx context.Context, f io.ReadCloser, name string, args ...string) (err error) {
	defer func() {
		if err2 := f.Close(); err == nil && err2 != nil {
			err = err2
		}
	}()
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdin = f // os.File, so special case, no need for a pipe to support large files

	if out, err2 := cmd.CombinedOutput(); err2 != nil {
		err = errors.Wrapf(err2, "[dmltest] SQLDumpLoad\n%s", string(out))
		return
	}
	return
}

// SQLDumpOptions can set a different MySQL/MariaDB binary path and adds more
// arguments.
type SQLDumpOptions struct {
	MySQLPath string
	MySQLArgs []string
	DSN       string
	// mocked out for testing.
	execCommandContext func(ctx context.Context, file io.ReadCloser, cmd string, arg ...string) error
}

// SQLDumpLoad reads all files recognized by `globPattern` argument into
// MySQL/MariaDB. The password will NOT be visible via process manager but gets
// temporarily written into the TMP dir of the OS. This function does even work
// when the server and the client runs on different machines. For now it only
// works when the program `bash` has been installed. This function supports any
// file size of a `.sql` file. Bonus: if file names contain the string
// "cleanup", they will be run in the defer function. The returned function must
// be run in the defer part of a test. This function skips a test, if the DSN
// environment variable cannot be found.
func SQLDumpLoad(t testing.TB, globPattern string, o *SQLDumpOptions) func() {
	if o == nil {
		o = &SQLDumpOptions{}
	}

	if o.DSN == "" {
		o.DSN = MustGetDSN(t)
	}

	matches, err := filepath.Glob(globPattern)
	FatalIfError(t, err)
	if len(matches) == 0 {
		FatalIfError(t, errors.NotFound.Newf("No files found for glob pattern: %q", globPattern))
	}

	cfg, err := mysql.ParseDSN(o.DSN)
	FatalIfError(t, err)

	execCmd := o.execCommandContext
	if execCmd == nil {
		execCmd = stdInExec
	}

	dfFile, err := writeMySQLDefaults(cfg, o)
	FatalIfError(t, err)

	myPath := o.MySQLPath
	if myPath == "" {
		myPath = "mysql"
	}
	ctx := context.TODO()

	runExec := func(file string) {
		f, err := os.Open(file)
		FatalIfError(t, err)

		err = execCmd(ctx, f, "bash", "-c", fmt.Sprintf("%s --defaults-file=%s", myPath, dfFile))
		FatalIfError(t, err)
	}

	var cleanUpFiles []string
	for _, file := range matches {
		if strings.Contains(file, "cleanup") {
			cleanUpFiles = append(cleanUpFiles, file)
		} else {
			runExec(file)
		}
	}

	return func() {
		defer os.Remove(dfFile)
		for _, file := range cleanUpFiles {
			runExec(file)
		}
	}
}
