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

// +build file csall

package objcache

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/corestoreio/errors"
)

// FileSystemConfig allows to overwrite default values.
type FileSystemConfig struct {
	Path string
	// FileSuffix of a cache file, default ".bin"
	FileSuffix string
	// DirectoryLevel either 1 or 2 levels deep. Default 2.
	DirectoryLevel int
	DirectoryMode  os.FileMode // default 0700
	FileMode       os.FileMode // default 0600
	Location       *time.Location
	CleanOnClose   bool
}

// NewFileSystemClient creates a new file system client. Argument `c` can be nil.
// Warning: do not change manually the files mtime and atime.
func NewFileSystemClient(c *FileSystemConfig) NewStorageFn {
	if c == nil {
		c = &FileSystemConfig{}
	}
	if c.Path == "" {
		c.Path = "testdata/fs"
	}
	if c.FileSuffix == "" {
		c.FileSuffix = ".bin"
	}
	if c.DirectoryLevel == 0 {
		c.DirectoryLevel = 2
	}
	if c.DirectoryMode == 0 {
		c.DirectoryMode = 0700
	}
	if c.FileMode == 0 {
		c.FileMode = 0600
	}
	if c.Location == nil {
		c.Location = time.Local
	}
	return func() (Storager, error) {
		return &fileStorage{
			cfg:           *c,
			keyToFileName: make(map[string]string),
		}, nil
	}
}

type fileStorage struct {
	cfg FileSystemConfig

	mu            sync.Mutex
	keyToFileName map[string]string
}

func (fs *fileStorage) writeFileContent(fileName string, data []byte, expires time.Duration) (err error) {
	var f *os.File
	f, err = os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fs.cfg.FileMode)
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		n := time.Unix(0, 0)
		if expires != 0 {
			n = now().In(fs.cfg.Location).Add(expires)
		}
		if err2 := os.Chtimes(fileName, n, n); err == nil && err2 != nil {
			err = errors.WithStack(err2)
		}
		if err2 := f.Close(); err == nil && err2 != nil {
			err = errors.WithStack(err2)
		}
	}()

	if _, err = f.Write(data); err != nil {
		return errors.WithStack(err)
	}
	err = f.Sync()
	return
}

// Put writes the values into a file by using the keys as their filename. The
// keys will be sha1ed. If no duration gets set, the file atime and mtime
// applies to 1970-01-01 which means that the content never expires. once the
// mtime and atime is greater than the year 1970, the expiration gets checked.
func (fs *fileStorage) Put(_ context.Context, keys []string, values [][]byte, expires []time.Duration) (err error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	for i, key := range keys {
		kfn, err := fs.getCacheFileName(key)
		if err != nil {
			return errors.WithStack(err)
		}
		var e time.Duration
		if len(expires) > 0 {
			e = expires[i]
		}
		if err := fs.writeFileContent(kfn, values[i], e); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (fs *fileStorage) lookupFileContent(fileName string, values [][]byte) ([][]byte, error) {
	if ok, err := fileExists(fileName); !ok {
		return append(values, nil), err
	}

	f, err := os.Open(fileName)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if mt := fi.ModTime(); mt.Year() > 1970 && mt.Before(now()) { // too old, delete the file which cleans up the cache
		err := os.Remove(fileName)
		values = append(values, nil)
		return values, err
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(f); err != nil {
		return nil, errors.WithStack(err)
	}
	values = append(values, buf.Bytes())
	return values, nil
}

func (fs *fileStorage) Get(_ context.Context, keys []string) (values [][]byte, err error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	for _, key := range keys {
		kfn, err := fs.getCacheFileName(key)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		values, err = fs.lookupFileContent(kfn, values)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return values, nil
}

func (fs *fileStorage) Delete(_ context.Context, keys []string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	for _, key := range keys {
		kfn, err := fs.getCacheFileName(key)
		if err != nil {
			return errors.WithStack(err)
		}
		if err = os.Remove(kfn); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (fs *fileStorage) Truncate(ctx context.Context) (err error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	fs.keyToFileName = make(map[string]string)
	return os.RemoveAll(fs.cfg.Path)
}

func (fs *fileStorage) Close() (err error) {
	if fs.cfg.CleanOnClose {
		err = fs.Truncate(context.Background())
	}
	return nil
}

func (fs *fileStorage) getCacheFileName(key string) (string, error) {
	if fileName, ok := fs.keyToFileName[key]; ok {
		return fileName, nil
	}

	m := sha1.New()
	m.Write([]byte(key))
	keyMd5 := hex.EncodeToString(m.Sum(nil))
	cachePath := fs.cfg.Path
	switch fs.cfg.DirectoryLevel {
	case 2:
		cachePath = filepath.Join(cachePath, keyMd5[0:2], keyMd5[2:4])
	case 1:
		cachePath = filepath.Join(cachePath, keyMd5[0:2])
	}

	ok, err := fileExists(cachePath)
	if err != nil {
		return "", errors.Wrapf(err, "[objcache] FileStorage CachePath %q", cachePath)
	}
	if !ok {
		err = os.MkdirAll(cachePath, fs.cfg.DirectoryMode)
		if err != nil {
			return "", errors.Wrapf(err, "[objcache] FileStorage CachePath %q", cachePath)
		}
	}
	fileName := filepath.Join(cachePath, keyMd5+fs.cfg.FileSuffix)
	fs.keyToFileName[key] = fileName
	return fileName, nil
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		err = errors.WithStack(err)
	}
	return false, err
}
