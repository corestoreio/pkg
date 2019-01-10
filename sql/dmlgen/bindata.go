// Code generated by go-bindata. DO NOT EDIT.
// sources:
// _tpl/10_tables.go.tpl
// _tpl/20_entity.go.tpl
// _tpl/30_collection_methods.go.tpl
// _tpl/40_binary.go.tpl
// _tpl/90_test.go.tpl
// _tpl/fbs_10_header.go.tpl
// _tpl/fbs_20_table.go.tpl
// _tpl/protobuf_10_header.go.tpl
// _tpl/protobuf_20_table.go.tpl

package dmlgen


import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}


type asset struct {
	bytes []byte
	info  fileInfoEx
}

type fileInfoEx interface {
	os.FileInfo
	MD5Checksum() string
}

type bindataFileInfo struct {
	name        string
	size        int64
	mode        os.FileMode
	modTime     time.Time
	md5checksum string
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) MD5Checksum() string {
	return fi.md5checksum
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _bindataTpl10tablesgotpl = []byte(
	"\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x51\xc1\x6e\xdb\x3a\x10\x3c\x93\x5f\x31\xcf\x78\x07\xa9\x60\xe9\xbb" +
	"\x0b\x1d\x1a\xa7\xe8\xa5\x49\x1b\x24\x40\xcf\xb4\xb4\xb2\x89\x4a\xa4\xb1\xa4\x5b\x1b\x02\xff\xbd\x58\xc9\x48\x7c" +
	"\x6b\xd1\x93\xa8\xd9\xd9\xd9\xd9\x9d\x36\x86\x94\x51\xe9\x69\x7a\x0f\x76\x61\x4f\xf8\x3f\xbb\xdd\x40\xd8\x34\xb0" +
	"\x2f\xf2\x4a\x28\x45\xab\xf9\xf9\xe8\x46\x9a\xa6\x97\xf8\x39\x6e\xdd\x48\xc3\xd6\x25\xba\x92\xa4\x52\x0a\x1a\xac" +
	"\xa6\xe9\x16\x59\xcd\xc2\x14\xba\x52\x74\xad\xf5\x7a\x8d\x47\xfa\x75\x95\x65\xca\x27\x0e\x09\x0e\xfb\xe8\x7b\x4f" +
	"\x1d\x7e\x12\x27\x1f\x03\x62\x8f\x7c\x20\x3c\x5c\x9e\x9f\xbe\xac\x1f\x1c\x7b\x77\x7f\x87\xc5\x56\x6a\x0f\x34\x3a" +
	"\xf4\x91\x85\x22\x82\x33\x9e\x36\xb8\x1d\x9c\x4a\x91\xd2\xc7\x53\x8e\xd8\x53\x20\x76\x99\x3a\xec\x2e\xe8\xc6\x61" +
	"\x4f\xc1\xea\xfe\x14\xda\x37\x2b\x55\x9b\xcf\x68\x63\xc8\x74\xce\x76\xbb\x7c\x0d\xba\x9d\xd0\xed\xd3\x89\xf8\xf2" +
	"\xe9\x4c\xed\x37\xa6\xa3\x63\x62\x83\x78\xcc\x09\xd6\xda\xae\x1b\x96\x99\x5f\x8f\xd9\xc7\x50\xa3\xca\x23\xde\xbd" +
	"\xa2\xc9\x10\x33\x88\x39\x72\x8d\x49\x2b\xdf\x23\x8f\x46\x00\x34\x10\xd6\x9b\x01\xad\x94\x00\xdf\x7d\x3e\x6c\x99" +
	"\x5c\xa6\x19\x17\x5f\xe2\xc3\xe0\x4f\xf1\x28\xfc\x55\x40\x06\xab\x95\xd1\x4a\xbd\x86\x52\x9b\x9b\xc9\xf7\x77\x55" +
	"\xb7\x13\xa4\xfe\x30\x9b\xfc\xaf\x41\xf0\x83\x38\x57\x4b\x56\xf2\x6b\x96\x85\xd2\xdc\xf1\x9c\x5d\xfb\xa3\x22\xe6" +
	"\x5a\xab\x32\x2f\x28\x7d\x9b\x06\x79\xb4\xcb\x4d\x52\x25\xc7\xb2\xd6\xfe\xab\xe6\x95\x25\x87\x0b\x7e\xd0\x45\xff" +
	"\x0e\x00\x00\xff\xff\xea\x74\xb5\xf9\xb2\x02\x00\x00")

func bindataTpl10tablesgotplBytes() ([]byte, error) {
	return bindataRead(
		_bindataTpl10tablesgotpl,
		"_tpl/10_tables.go.tpl",
	)
}



func bindataTpl10tablesgotpl() (*asset, error) {
	bytes, err := bindataTpl10tablesgotplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{
		name: "_tpl/10_tables.go.tpl",
		size: 690,
		md5checksum: "",
		mode: os.FileMode(420),
		modTime: time.Unix(1543562479, 0),
	}

	a := &asset{bytes: bytes, info: info}

	return a, nil
}

var _bindataTpl20entitygotpl = []byte(
	"\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x57\x5f\x6f\xe3\xb8\x11\x7f\x96\x3e\xc5\x9c\xd1\x3b\x48\x81\x43\xef" +
	"\x43\xdb\x87\x2c\xfc\x90\x4b\xb2\xdb\xb4\x5d\xdf\x62\x93\x43\x1f\x82\xa0\x61\xa8\x91\xcd\x0b\x45\x6a\x49\xea\xbc" +
	"\x3e\x41\xdf\xbd\x18\x4a\xfe\x23\xd9\x4e\xb2\xbb\x28\x70\x01\xe2\xc4\xa2\x38\x7f\x7e\x33\xf3\x9b\x99\xc9\x04\xea" +
	"\x9a\x5d\x69\x2f\xfd\xaa\x69\xc0\x62\x69\xd1\xa1\xf6\x0e\x38\x38\xa9\xe7\x0a\xc1\x9a\x25\xe4\xc6\xc2\xe5\xcf\xe0" +
	"\xf9\xa3\x42\x78\xa8\x6b\x76\x4b\xff\xcd\x78\x81\x4d\xf3\xc0\xe2\xc9\x04\xce\x2b\x6f\x60\x8e\x1a\x2d\xf7\x98\xb1" +
	"\xba\x5e\x4a\xbf\x00\x76\x61\x8a\x02\xb5\x6f\x9a\xb8\xae\x19\x9c\x36\x4d\x5d\xa3\xce\xe8\xcf\x29\xc8\x1c\xd8\x3f" +
	"\xb8\xbb\xe2\x6e\xf5\x4f\x67\xf4\x07\x6e\xdd\x82\x2b\xb4\xd0\x34\xf1\x64\x82\xdc\xad\x7e\x73\x46\x9f\xd1\x47\x77" +
	"\x2b\xf6\xab\x12\x7b\x06\x3b\x6f\x2b\xe1\xa1\x8e\xeb\xda\x72\x3d\x47\xd2\xa8\xaa\x42\x3b\x52\x71\x6b\xde\x9b\x0b" +
	"\x5e\xa0\xba\xe0\x0e\x81\xbd\x93\xa8\xb2\xa6\x81\xba\x7e\x6f\x6e\x57\x25\xce\x2a\xa5\x80\x35\x4d\x1c\x45\x9d\x39" +
	"\x1a\x81\xdd\x04\x89\xb7\x7c\x0e\xa3\x11\x19\x4c\xde\x6e\x9e\x85\xaf\xa7\x10\x8c\x21\x3b\xde\x9b\x5d\x07\xdb\xc7" +
	"\x4d\x1c\xf0\x70\x4e\xce\xf5\xbf\xb9\xf3\xd7\xda\xa1\xf5\xd7\x97\x50\x95\x19\xf7\xe8\xc0\x2f\x10\xa4\x16\x16\xe9" +
	"\x22\x5c\x5f\x42\x4e\x86\x41\x00\x8c\xce\x14\x77\x1e\x64\xb8\x85\x19\x5c\x5f\x92\xb8\xdc\x9a\x02\xb8\x86\xeb\xd9" +
	"\xcd\xd5\xa7\x5b\x30\x25\xe1\x2c\x8d\x66\x70\x5d\x94\x2a\x48\x72\x90\x15\x8a\xad\xb5\xb5\xfa\xd1\xb2\x61\x64\xe2" +
	"\xbc\xd2\x02\x12\x84\x93\x1d\x1c\xd3\x03\xf6\x26\x32\x03\xa9\xfd\xdf\xff\x9a\x42\x1d\x47\x87\xf0\x95\x39\x70\x9d" +
	"\x01\xbb\x76\x1f\xff\x45\x9f\xa4\xe9\x7a\xed\x58\xd3\x00\xb2\xa3\x31\x98\x6e\xa2\x40\x11\x48\x64\x96\x92\x8a\x2e" +
	"\x35\x50\x67\x84\x7c\xdc\x22\xf9\x81\x97\x9d\x4e\x90\x5b\x5f\xa5\xf6\x68\x73\x2e\x10\xda\xc3\x0f\xbc\x2c\xd1\x82" +
	"\xd1\x6a\x05\x25\xb7\x5e\x72\xa5\x56\xaf\x75\x7e\xab\x22\x11\x05\x9c\x10\x8e\x1b\xa9\x29\xa0\xb5\xc6\x12\x06\x32" +
	"\x07\x51\xb0\x0f\x26\xc3\x24\x85\xe9\x14\x7a\xef\xb5\xd2\x3e\x21\xcf\xce\x95\xa2\xb7\x23\x8b\xbe\xb2\x1a\x44\xb1" +
	"\x8f\x1d\x23\xef\xdf\x55\x5a\xac\x73\x30\xf9\xe9\x38\x58\x69\x07\x0c\xbb\xb2\x36\x49\xe3\xa8\x89\x23\xaa\x47\x51" +
	"\xb0\x19\x7e\xf1\x49\x08\x4f\xe4\x96\xd2\x8b\x05\x08\x38\x9b\xd2\x49\xab\x29\x49\xdf\x82\x80\x1a\xf6\x0d\x88\xa3" +
	"\x28\x12\xa4\x64\x54\xd7\x6b\x3d\xa3\xcd\x6b\xe7\x4a\x72\x87\xae\x69\xc6\x74\x1c\x4e\x82\x05\x67\x74\x2d\x12\xc5" +
	"\x37\x99\x4f\x77\x33\xcc\x79\xa5\x7c\x2b\xa7\xc3\x27\xc0\xeb\xd8\xcc\xf8\x77\xa6\xd2\x19\x9b\xe1\x32\x4f\x46\x77" +
	"\x75\xcd\x3e\x72\xf1\xc4\xe7\xd8\x34\xf7\xbd\x9a\x6f\x9d\x80\x1f\x3f\x83\x36\x1e\x72\xba\x34\x1a\x83\x48\xe3\x88" +
	"\xa0\x69\xe2\x81\xe0\xff\x48\xbf\xb8\xf1\x5c\x3c\x25\xa2\x68\x21\x4c\xbb\xcc\xba\x2a\x4a\xbf\x02\x2c\x4a\x2f\xd1" +
	"\x01\x57\x2a\x94\x5f\xa8\x46\x07\x26\x0f\xdf\x44\x65\x2d\x95\xa9\x79\xfc\x0d\x85\x67\x70\xae\x9c\x81\x27\x6d\x96" +
	"\x1a\xb8\x83\x4f\xe8\xd0\x1f\xcb\xaa\x20\x3e\x49\x7b\x4f\xa1\x86\x13\x0c\xd9\xbf\x79\x54\x37\x6f\x61\x6d\x71\x47" +
	"\x1e\x75\x4d\x91\x52\x28\xa8\xc4\x87\xb4\x2c\x36\x27\x10\xe8\xb0\xc7\xcd\x7d\x6a\x26\x59\x33\xe3\xc1\x2f\x2c\xf2" +
	"\x0c\x1c\xcf\x71\xaf\x24\xfe\xcf\x4c\xdd\xf3\x63\xc3\xd6\xd1\x25\xf7\x1c\xb6\x3f\x51\x74\x77\xdf\x83\xe9\x21\x08" +
	"\x1c\x65\xdc\xf3\xb1\x29\xa4\xa7\x20\xad\x46\x0f\x71\xf4\x33\xe6\xc6\xe2\xb6\x60\x23\xc2\x3e\xa9\x02\x51\x8d\x07" +
	"\x01\x68\x0b\xb7\x13\x75\x4a\xb7\xcf\x73\x8f\x76\x87\x50\xbe\xea\x76\x1b\x9b\x19\x2e\x87\x6e\x09\x8b\x81\xd8\x39" +
	"\x68\x5c\x82\xd4\x92\xd8\x47\xfe\x81\xd9\x4e\xac\x8e\x50\xd1\xbe\xb4\x2e\x61\x7a\x0a\xea\x4d\x4e\xff\x34\x38\xa2" +
	"\xd2\x27\x2c\xcf\xa0\xe0\x4f\x98\xf4\x51\x1c\xc3\x9b\x31\xfc\x2d\x1d\x53\x51\x34\x71\x97\xa5\x42\xec\xc9\x4f\xc1" +
	"\x09\xae\x8f\x32\xe0\xb8\x9f\xd8\x63\x90\xd9\x17\xa8\xba\xd6\xb0\x4b\x8e\x68\x6d\xa0\x1f\xc1\x86\x51\x4a\x64\xf6" +
	"\x65\x0c\x98\xbe\x0d\xef\xfc\x30\x05\x2d\x7b\x1c\xb9\x57\xaa\x68\x6d\xcb\x74\x5b\xb1\xc8\x7a\x34\xfd\x5d\xb2\x84" +
	"\x60\x83\x54\xf8\x76\x0b\xbb\x17\xb4\x54\xcf\x36\xac\x1e\xa4\xd4\xaa\x36\x1d\xec\x58\x97\x3a\x18\xa9\x57\xb6\xaa" +
	"\xae\x1d\x14\x5d\x3b\x68\x7b\xd6\x5b\x28\xe8\x2c\xf0\xfe\xf1\xde\x35\x3e\x76\x76\x83\x81\xb8\x89\x6f\xe4\x18\x90" +
	"\x44\xb7\xfd\x42\x08\x16\xea\x99\xe0\xea\x83\xdc\xcf\xab\x31\xe0\xb8\x4b\x9c\x44\xa6\xfb\x50\xbf\x04\x76\xe0\x77" +
	"\xfa\xdd\xf7\xe0\x46\x70\x4d\xc6\xb5\x2d\xfa\xc2\x54\xda\x53\x87\x7e\xd3\xca\x5d\x1b\x38\x5d\x9b\x7a\x77\xf6\xe6" +
	"\xbe\x15\x15\x05\x3f\x34\x2e\x93\xdd\xe2\x8f\x5f\xf6\x63\xad\xe6\x80\x1b\x2f\x78\x41\x5a\xb7\x16\x51\x32\xe8\x2c" +
	"\xe9\x1e\x50\x0a\x1e\xf2\x6f\x9b\x04\x83\x48\xf4\xa7\x80\x67\xc7\x00\x3a\x27\x36\xef\xcf\x02\xec\x57\x2d\x3f\x57" +
	"\xb8\xce\xd8\xd3\xef\x9d\x0d\x20\x28\x1e\x0e\x08\x8e\x1c\x3c\x36\x21\xb8\x24\x65\x8c\xa5\x6b\xfb\x36\xd3\x42\x5d" +
	"\x4f\x4e\xe0\x98\xc5\x32\x97\x98\x6d\xe6\x99\x28\x9a\x4c\x40\x1b\x90\x19\x72\x6a\x57\x7e\xc1\x3d\x48\x07\x1a\x31" +
	"\xc3\xec\x4f\xe1\x11\x9c\x4c\xbe\x7f\x08\xea\xb5\x85\x67\x06\xa1\x75\xa5\xec\xe8\xda\xd3\x74\x53\x95\xa5\xa1\x1d" +
	"\xe3\xb0\xb6\x5f\x75\x3b\xe6\x10\x77\x9c\xc1\x8f\x9f\x47\x63\xea\xdf\x52\xcf\x93\x22\xed\x31\xdf\x7a\xac\x8a\x9b" +
	"\xbd\x2d\xac\x9f\x5c\xed\x40\x72\x14\xb3\x6e\x0a\x0a\xab\xa7\x92\x02\xc1\xd8\xae\x3a\x1c\x78\xb3\x79\x4a\x73\xda" +
	"\xef\x5c\x55\xe8\x0e\x2d\x9e\x3d\xee\xfc\xcb\x80\x3c\x8f\x87\xcb\xa2\x07\xc6\xd8\x70\x35\x4c\xe1\xee\x7e\xf8\xac" +
	"\x6b\x76\x74\x63\xda\x6b\x13\x30\x5d\x37\xe1\xe1\x95\xd0\x88\x15\xea\x75\x95\xa7\xdb\x29\xfe\xbf\x47\x99\xb4\x95" +
	"\xd8\xd1\x83\x45\x3f\x7e\x66\x9f\xea\xc5\xc3\xa2\x8f\x69\x3b\x6d\xeb\x28\x7e\x5d\x01\x3d\x1f\x99\x47\x54\x46\xcf" +
	"\x43\x18\xc2\x50\xdc\x26\xde\xc3\xb6\x9a\x1e\xe8\x3e\x2d\x83\xaf\x8b\x61\x58\xd3\xaa\x90\x1b\x5d\x2c\xdb\x79\x9b" +
	"\xfb\x4e\x76\x08\xed\xed\x62\x73\xba\x94\x4a\xc1\x23\x0d\xe7\xca\xa3\xc5\xac\xed\x9e\x9a\x96\x3c\x90\x1a\x38\xbc" +
	"\x37\x50\xf0\x92\xc1\xcc\xd0\x30\xfc\xb9\x42\xbb\x82\x39\x7a\x47\x72\xf0\x0b\x8a\x8a\x92\xe3\xab\x92\xa5\x4d\xdd" +
	"\x57\xa7\xcc\x30\x5d\xbe\x32\x55\x9e\x49\x13\xa2\xc1\x38\x8a\x6e\x7f\xb9\xfc\xe5\x0c\x38\x58\xac\x5c\x98\xf5\x0b" +
	"\x5e\x06\xc8\x2b\x6a\x16\x32\xcf\x31\x2c\x2a\x5c\xcd\x8d\x95\x7e\x51\x38\xc8\x90\x80\x97\x7a\x0e\xb4\x2b\x2c\x10" +
	"\x9c\xfc\x03\xe3\x28\x5a\xaf\x36\x5d\xb2\x85\x90\x30\xb8\x31\x05\x7a\x59\x84\x39\x36\x37\x76\x42\xd9\xa9\x8c\x29" +
	"\xc1\x56\xda\x41\xce\x9d\x47\x4b\x21\x22\xb4\x09\xea\x38\x6a\x09\x2d\xab\xca\x8b\x05\x8a\x27\xca\xe2\xe0\x54\xc1" +
	"\xcb\xbb\x5d\xb7\xee\x1f\x8d\x51\x7b\x8e\x3d\x9f\xfc\x32\x87\x1f\xd6\x82\xef\x30\x80\xfc\xd1\xca\x42\x7a\xf9\xfb" +
	"\xa6\xaa\xee\x37\xdd\xf6\x40\x9d\x1c\x78\x3f\xf0\xe2\xcb\x42\xa7\xe0\x6d\x85\xc3\x65\xb2\x5f\x54\xff\x0b\x00\x00" +
	"\xff\xff\xa6\xc2\x23\xd9\x37\x13\x00\x00")

func bindataTpl20entitygotplBytes() ([]byte, error) {
	return bindataRead(
		_bindataTpl20entitygotpl,
		"_tpl/20_entity.go.tpl",
	)
}



func bindataTpl20entitygotpl() (*asset, error) {
	bytes, err := bindataTpl20entitygotplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{
		name: "_tpl/20_entity.go.tpl",
		size: 4919,
		md5checksum: "",
		mode: os.FileMode(420),
		modTime: time.Unix(1544711384, 0),
	}

	a := &asset{bytes: bytes, info: info}

	return a, nil
}

var _bindataTpl30collectionmethodsgotpl = []byte(
	"\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xa4\x95\x41\x8f\xdb\x36\x10\x85\xcf\xd2\xaf\x78\x87\xa2\x90\xba\xb2\xdc" +
	"\xbd\x3a\xf1\x21\x75\x92\x22\xe7\x4d\x4f\x86\x51\xd0\xd4\xc8\xa2\x4d\x91\x2e\x45\xad\x61\x09\xfa\xef\x05\x29\xc9" +
	"\x96\x77\x1d\x04\x81\x81\x3d\x88\x4b\x72\xe6\xcd\xf7\x66\xe8\xf9\x1c\x5f\x85\xb4\x64\xbe\x17\xa2\x42\xee\x3f\x2b" +
	"\xd8\x82\xc0\x6b\x63\x48\x59\x54\x52\x70\xc2\xf6\x8c\xa3\xa1\x4c\x70\x66\x09\x39\x4e\xc2\x16\xba\xb6\x28\xa9\xd4" +
	"\xe6\x0c\x26\xa5\xe6\xcc\x0a\xad\xd2\x70\x3e\xc7\xa7\xda\x6a\xec\x48\x91\x61\x96\x32\xbc\x0a\x86\xac\x94\x3b\x52" +
	"\x69\x98\xd7\x8a\x23\xe2\x1c\x7f\xb4\x6d\xba\xd2\x52\x12\x77\xd7\xba\x2e\x1e\x74\x44\x39\xdc\x99\xc8\xed\x7f\x51" +
	"\x56\xd8\xb3\xdb\xdb\x6a\x2d\xe3\x77\x77\xd0\x86\xc1\x16\x8b\x25\x38\x4f\x3f\x33\xcb\xd6\x8b\x3f\x37\x61\x90\x6b" +
	"\x83\x7f\x13\x90\xdb\x30\x4c\xed\x68\xdc\x76\xc7\x03\x91\x23\x8f\x28\xf6\xdf\xc1\x16\x4b\xb0\xe3\x91\x54\x16\x6d" +
	"\x13\x50\x1c\x06\x41\x17\xba\xbf\xf1\xc6\x12\xdb\x30\x30\x64\x6b\xa3\xc0\x79\xd8\x85\xae\xbc\x2f\x8c\x17\x38\x09" +
	"\x29\x61\x6a\xe5\xd5\x3a\x39\xc8\xa1\x95\x23\x01\x61\xa9\xac\x20\x14\xd6\x9b\x69\x15\x0f\xa0\x71\x19\xef\x82\xb9" +
	"\xcf\xc4\x11\x10\xf7\xcb\xcf\xa3\x91\x95\xd8\xc4\xbe\xd4\x9b\xea\xda\x76\xb8\xe3\x62\xd6\xa5\xaa\xd2\x7f\x94\xf8" +
	"\xaf\xa6\x61\x85\x59\xd7\xdd\xe8\xfc\xed\x8d\xd0\x17\x6d\xec\x5f\xe7\xb6\xfd\xae\xff\xd6\x2b\x56\x92\x5c\xb1\x8a" +
	"\x90\x7e\x15\x24\xb3\xae\x8b\x3c\xf6\x4a\x1b\x9b\xbe\xb8\xae\x1a\xb5\x24\xbe\x30\x91\x60\x0f\xa1\x6c\xef\xb6\x57" +
	"\x7b\xd1\x36\x4a\x4e\x7f\x14\x1a\x1f\x2f\xa7\xf6\x3f\x3e\x15\x06\x5d\x1c\x76\x68\xdb\x19\x48\xb9\xb5\xb3\x64\x55" +
	"\xdb\xc1\x4d\x2a\xf5\x2b\x8d\xf6\xc1\x16\x46\xd7\xbb\x02\xfb\xd9\xf3\x03\xd6\xad\x6a\x3b\xa9\xec\x9e\x5b\xcd\xa4" +
	"\x83\x31\x9f\x83\xeb\xe3\x79\x98\xba\x82\x58\x46\x26\x0c\xdc\xbf\xa2\x66\x2d\x16\x9b\x04\xcd\x7a\xbf\x70\xd6\x39" +
	"\x93\x0f\x09\x94\xbb\x2d\x49\x45\x4d\x3c\xdb\x3f\x89\x64\xf8\xfe\x80\x03\x3e\x42\x7d\xc0\xe1\xe9\xc9\xa3\x6c\xd6" +
	"\x87\x0d\x96\x50\x42\xba\x1c\xd6\x4d\x7b\x55\xe8\x5a\x66\x60\xaf\x5a\x64\x7e\xe4\x87\x69\x96\xc4\x0e\xbe\x35\x1a" +
	"\x2c\xd1\xac\x17\xd7\xe8\x9b\xe9\x64\x34\xef\x27\xe3\xe5\xc4\x8e\x3d\xcb\x8a\x59\x51\xe5\x67\x1f\xd6\x1b\xfe\x4d" +
	"\x59\x32\x39\xe3\xf4\x00\x4b\x17\x7e\x02\xb3\x9d\x34\x46\x32\xb1\x1f\xcb\xc9\x22\x99\x1c\x42\xaf\xf2\x33\x49\xb2" +
	"\x74\xe3\x39\x53\xde\x76\xe4\x46\x97\xbd\x66\xc7\xff\x01\xa9\x7d\x8e\x48\xfc\xa2\xed\x97\xd4\x17\xeb\x49\x65\x57" +
	"\x87\x31\xc3\xb3\xb7\x60\x24\x41\x2a\x8b\xdf\xb6\x87\x78\x7a\xf6\x0d\xd2\xac\x49\x65\xbf\xe6\xf9\x60\xb8\xbb\xf7" +
	"\x13\xa7\xbf\xa9\x8a\xcc\x30\x37\x47\xc9\x38\x81\x41\xd1\xa9\xa7\xc8\x2c\x8e\xba\x12\xfe\x55\x14\x0f\x40\xec\x93" +
	"\x44\x0a\xd3\x37\x2f\xc1\xe3\x50\x9b\xeb\xbb\xdf\x24\xf8\x7d\x12\xbd\xed\x26\x3c\x1d\x48\x4f\x74\xe0\x29\x3c\xcd" +
	"\x9f\x90\xf9\xe4\xe3\xf6\x64\x58\x96\xbd\xe1\xe2\xb4\xb8\x6d\x9d\xbf\xd3\xff\x00\xa8\x3e\x67\xa4\x90\xa6\xe9\xed" +
	"\x0f\xe7\x3d\x48\x57\xf9\x03\x83\xf1\x19\x86\x4a\xd3\x34\xbe\x2d\xe9\xff\x00\x00\x00\xff\xff\xed\xda\x0a\x7c\x21" +
	"\x08\x00\x00")

func bindataTpl30collectionmethodsgotplBytes() ([]byte, error) {
	return bindataRead(
		_bindataTpl30collectionmethodsgotpl,
		"_tpl/30_collection_methods.go.tpl",
	)
}



func bindataTpl30collectionmethodsgotpl() (*asset, error) {
	bytes, err := bindataTpl30collectionmethodsgotplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{
		name: "_tpl/30_collection_methods.go.tpl",
		size: 2081,
		md5checksum: "",
		mode: os.FileMode(420),
		modTime: time.Unix(1543131929, 0),
	}

	a := &asset{bytes: bytes, info: info}

	return a, nil
}

var _bindataTpl40binarygotpl = []byte(
	"\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\xcf\xb1\x4e\xc3\x30\x10\xc6\xf1\x19\x3f\xc5\x37\x30\xc4\x08\xf9\xde" +
	"\x01\x26\x86\x8e\x4c\x88\xc1\xb9\x5c\x53\x4b\xc9\x5d\x75\xb5\x91\xaa\x28\xef\x8e\xa0\x50\x51\xa6\x76\xf6\xdf\x9f" +
	"\xee\x47\x84\x57\x9d\xb3\x1f\x76\x79\x7a\x2a\x9a\xfd\x88\x32\xef\x27\x99\x45\xeb\x01\xa2\x6c\x43\xd1\x31\x9d\x9e" +
	"\xce\xa5\x78\x0a\xdb\xa6\x8c\x8e\x19\x0f\xcb\x72\x9f\x9e\x6d\x9a\x84\x6b\x31\x5d\xd7\xf8\x7f\xb2\x1b\x72\xcd\x78" +
	"\x7b\xef\x8f\x55\x22\xc4\xdd\x1c\x4b\xb8\x73\xa9\xcd\x15\xcc\xe9\xdc\x7f\x97\x11\x44\x78\xf9\xbd\x42\x06\x7c\x94" +
	"\x8c\xb1\xd4\x5d\xeb\x13\xdb\x4c\xa3\x8d\x46\x7b\xb7\x6a\x7d\xdb\x86\x35\x04\x22\x6c\xae\x24\x6c\xae\x03\x5c\xcc" +
	"\x75\x11\x7f\x05\x8f\x5f\x82\x93\x22\x5e\x32\x7e\x7e\x75\x11\x37\x0a\x3e\x03\x00\x00\xff\xff\x22\xed\xb8\x7c\x86" +
	"\x01\x00\x00")

func bindataTpl40binarygotplBytes() ([]byte, error) {
	return bindataRead(
		_bindataTpl40binarygotpl,
		"_tpl/40_binary.go.tpl",
	)
}



func bindataTpl40binarygotpl() (*asset, error) {
	bytes, err := bindataTpl40binarygotplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{
		name: "_tpl/40_binary.go.tpl",
		size: 390,
		md5checksum: "",
		mode: os.FileMode(420),
		modTime: time.Unix(1543131929, 0),
	}

	a := &asset{bytes: bytes, info: info}

	return a, nil
}

var _bindataTpl90testgotpl = []byte(
	"\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xbc\x56\xdd\x6f\xdb\x36\x10\x7f\xa6\xfe\x8a\x9b\xb0\x16\x52\xab\x32\x2d" +
	"\x30\x14\x98\xbb\x3e\x34\x52\xd2\x19\x4b\x93\xac\x31\xba\x01\x45\x51\xd0\xe2\xd9\x26\x4a\x93\x02\x79\xaa\x1d\x08" +
	"\xfa\xdf\x07\x52\xf2\x47\xbf\xb2\x3c\x0c\x7b\xb3\xc9\xbb\xe3\xfd\x3e\x78\xd4\xa2\x35\x35\xcc\xd0\xd3\x25\x6e\x66" +
	"\x62\xae\xd1\x67\x04\x8f\x08\x3d\x29\xb3\xe4\xb3\x1c\xba\x84\xc9\x39\x4c\x5e\x82\x5c\xeb\xb0\xcc\xdf\xb4\x9e\x4a" +
	"\x6b\x0c\xd6\x54\x9d\x66\x94\x27\x4c\xe2\x02\xdd\x7e\xbf\xd4\xd6\x63\x46\x05\xc8\x79\x9e\x24\xac\xeb\x36\x8a\x56" +
	"\xc0\xc3\x19\x37\x7f\x5e\x54\xed\xba\x79\xad\xed\xfc\x5a\xd0\xaa\xef\xbf\xcc\x1c\xb7\x2f\xac\x90\x21\x3f\xed\x3a" +
	"\xde\xf7\x69\x01\x0f\xbf\x0a\xb8\x6a\x48\x59\xe3\xbb\x84\xb1\x9b\x4f\xaa\xa9\x4e\x4b\x8d\xc2\xb4\xcd\x04\xc8\xb5" +
	"\x58\x24\xac\xcf\xb3\xbc\xeb\xd0\xc8\xbe\x4f\x12\x56\xd3\x36\xf4\x5f\x5b\x43\xb8\x25\x3e\xbb\xaa\xae\xb2\x3c\x61" +
	"\x34\xd7\xbe\x00\x74\x2e\x6c\x1e\xd0\xd7\xb4\x0d\xad\xf3\xea\x34\x4f\x98\xf0\x1e\x1d\xf1\x4b\x7b\xe6\x9c\x75\xa1" +
	"\x29\x74\x2e\xa0\xa2\xb9\xbe\x14\x6b\xf4\x21\x37\x14\xe2\x63\x76\x9e\x30\x6f\x1d\xf1\x1b\x72\xca\x2c\x7d\xb6\x8b" +
	"\x3b\xd4\x3a\xdb\x8a\x9a\xf4\x6d\xa8\xf5\xfe\x83\x8f\x61\x1d\x74\xdd\x13\x70\xc2\x2c\x11\x7e\xa6\x50\x28\x94\x1d" +
	"\x4b\x42\xdf\xa7\x5d\x37\xfe\x0b\xb5\x02\x25\x21\x3e\xc2\xeb\x0b\x38\x1c\x91\xb0\x80\x66\x6c\xe8\x9d\xd0\x4a\x0a" +
	"\xc2\x00\xe8\xc7\x48\xd8\x67\xe1\xa0\xf1\xf0\xa8\xf1\xd8\x4a\xcb\x6f\xd0\x7d\x56\x35\x26\xac\xf1\xf0\x12\xc6\xc5" +
	"\xa0\xf8\x25\x6e\xc6\xbd\xec\x69\x01\x0f\xc7\x9d\x9d\x12\x17\xc2\x2c\x27\x90\x4a\x4c\x8b\x73\x6d\x05\xbd\x11\xdb" +
	"\x0a\x6b\xb5\x16\xda\x4f\x9e\xf7\x45\xc2\xd8\x98\xf0\x97\xa2\xd5\x4c\x2c\xcf\xc5\x27\x3c\x6f\x4d\x9d\xa5\x1b\x9c" +
	"\x7b\x45\xf8\x51\xc9\xb4\x80\x60\xc6\x6c\x2d\xb6\x17\x68\x40\x19\xca\x21\x53\x86\xd0\x2d\x44\x8d\x5d\x1f\x3b\xb6" +
	"\x2e\x3a\x92\x31\x87\xd4\x3a\x03\xcf\x0a\x30\x4a\x27\x8c\xf5\xf9\x5d\xa7\x78\xb2\xee\xbf\x3a\x23\x70\x5f\xb6\x9e" +
	"\xec\xba\xb4\x12\x21\xfd\x2e\x49\x23\x33\x29\x3c\xe9\xfb\x84\x05\x69\x4e\x4e\x20\x38\x0f\x5c\x6b\x80\x56\xd6\x23" +
	"\x04\x4b\x7b\x50\x06\x1a\xe1\x84\xd6\xa8\x13\x76\xb7\x0f\x12\x46\xfc\x6d\x6b\xb2\xb4\xeb\x66\xf6\xb5\x2d\xc5\x1a" +
	"\x75\x29\x3c\x7e\x61\x8e\x8f\x67\x86\x14\xdd\xee\x90\x7e\x73\x99\x59\x5d\xcb\xbd\x6b\x43\xd3\x31\x37\xdb\x57\xb8" +
	"\xab\x76\xc0\xc1\x94\xb9\xa1\x35\xed\x6f\x4e\x5d\x4b\x3e\x35\xc1\x5c\x59\xce\x4f\x5b\xa5\xe5\x3b\xa1\xdb\x70\x15" +
	"\xf8\xb5\xc3\x46\xb8\xc1\x81\x70\x72\x02\x95\x05\x63\x09\x5a\x8f\x30\x5d\x1a\xeb\x30\xcb\x81\x2c\xf8\xb6\x69\x1c" +
	"\x7a\x0f\xd5\xe9\x40\xbf\xe7\x09\xfb\xbe\x61\x0b\x48\x1f\x3c\xfe\x9c\xee\xcc\xcb\x94\xf1\xaf\x1c\x29\x2f\x4c\x68" +
	"\x65\xe8\x2c\xaa\xff\xca\x2d\xe3\x65\xfc\xd1\x6c\x1a\x42\x23\x1e\x8f\xfa\xa8\x46\x80\x73\x83\x1a\x6b\x3a\xbd\xbd" +
	"\xfe\x23\xcb\x8f\xaa\xf1\xb3\x6d\x23\x8c\xbc\xd6\xa2\xc6\xdf\xad\x96\xe8\xc2\x11\x09\x63\x0b\xeb\x40\x85\xe4\xa7" +
	"\x2f\x40\xc1\x6f\xf0\xeb\x0b\x50\x8f\x1f\x0f\x1e\xc2\xa8\xc6\x34\xd6\x36\xb8\xc9\xee\xa6\x97\x31\xa6\x16\x3b\x66" +
	"\x1b\xcf\x83\x83\x2b\x41\x22\xdb\x95\xc9\x5f\xc4\xdd\x9f\x5e\x06\x4b\x0e\x27\x30\xe2\x91\xa2\x45\x96\x4e\xab\xbf" +
	"\xdf\x3f\x90\x1f\x26\x30\x90\xa4\xf6\x3c\xed\xbc\x1c\x7e\x86\x91\xc8\x98\x9e\x56\xc7\x43\xbd\x5c\x61\xfd\xe9\x42" +
	"\x78\x1a\xa4\x9c\x56\x71\x00\xc7\xba\x93\x2f\x9f\x07\x7e\x1f\xf3\xe5\xd9\x41\x19\xfe\x16\x6b\xeb\x64\x96\x06\xd9" +
	"\x76\x30\xf8\xd9\x16\xeb\x72\x98\xc8\xd1\x1f\x03\xf6\xe3\x24\x8f\x34\xd0\x3b\x72\x78\xd5\xd2\xbd\x49\x74\x76\x53" +
	"\xda\xd6\x1c\x5c\x7a\x10\x99\x4f\x0d\x3d\xff\xc5\x67\x7a\x5a\xe5\x3c\x3e\x35\x71\xe2\xef\xcf\x88\xf9\xf7\x32\xdf" +
	"\x77\x26\x7a\xab\x42\xf1\xec\x59\x5e\xc0\xa1\x85\x20\xcb\x03\x39\x81\xb7\xe3\x0a\x48\x25\xe3\x45\x58\x0b\xaa\x57" +
	"\x41\xa7\x01\xe6\xd1\xe5\xaf\xad\x0e\x5d\x0f\x43\x80\x97\x56\xb7\x6b\x33\x4c\x80\x18\xa7\x16\x31\x84\x4f\xfd\xf0" +
	"\xce\x0c\x53\x86\x7d\xd3\xd2\x05\x9a\x25\xad\x42\x63\x5d\x17\x13\xca\x95\x70\x6f\xe2\xec\x5b\xd2\x6a\x60\x22\x3c" +
	"\x21\x0f\x77\xba\x7c\x23\x6e\xcc\x3a\x57\xa8\xe5\x51\xdc\x55\x4b\xff\x12\xb8\x83\x7c\x57\x14\xf8\x95\x6d\xb5\xdc" +
	"\xb3\x10\x04\x19\xf1\x3d\x01\xd4\x1e\x41\x2d\x22\x4d\x3b\xa8\xb7\x9e\x70\xfd\x0e\x9d\x57\xd6\xa0\xfc\x11\xe6\x28" +
	"\xd6\xfd\xe0\xfc\x7f\x68\xe2\x97\xc8\x57\x7f\xfa\xf0\x9d\x92\x1c\x2d\xf5\xc9\x3f\x01\x00\x00\xff\xff\xa3\x26\xfd" +
	"\x2b\x8d\x09\x00\x00")

func bindataTpl90testgotplBytes() ([]byte, error) {
	return bindataRead(
		_bindataTpl90testgotpl,
		"_tpl/90_test.go.tpl",
	)
}



func bindataTpl90testgotpl() (*asset, error) {
	bytes, err := bindataTpl90testgotplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{
		name: "_tpl/90_test.go.tpl",
		size: 2445,
		md5checksum: "",
		mode: os.FileMode(420),
		modTime: time.Unix(1547103572, 0),
	}

	a := &asset{bytes: bytes, info: info}

	return a, nil
}

var _bindataTplFbs10headergotpl = []byte(
	"\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x8d\x31\x4e\x43\x31\x10\x05\x7b\x9f\x62\x15\xd1\xf2\xdd\x13\x51\xd0" +
	"\xd1\x81\xc4\x09\x36\xf6\xc3\x58\xf1\xdf\xfd\xac\x6d\x0a\x56\xbe\x3b\x22\x07\x48\x33\xd5\x68\x26\x46\x7a\x99\x43" +
	"\xa9\x40\x60\x3c\x90\xe9\xa7\x32\x95\x3a\xbe\xe6\x65\x4b\xba\xc7\xa4\x86\x3e\xd4\x50\x35\x1e\xd7\x12\xfb\x77\x8b" +
	"\x79\x6f\x05\x12\x82\xf0\x8e\x7e\x70\x02\xb9\x6f\xef\x9c\xae\x5c\xb0\xd6\x39\x84\x2a\xa9\xcd\x0c\x3a\xdd\x0b\x0d" +
	"\x35\x2e\x88\x32\x5b\xbb\x61\xfb\xbc\xf4\xd3\x39\x04\x77\x63\x29\xa0\x07\x3d\x46\xa7\xa7\x67\xda\x3e\x60\x95\x5b" +
	"\xfd\x85\xbd\x82\x33\xec\xed\x18\x55\xa5\xd3\xe3\x5a\xc1\xfd\xe6\xfd\x6f\xdd\x21\x79\xad\xf0\x17\x00\x00\xff\xff" +
	"\x85\x4a\x09\x56\xd4\x00\x00\x00")

func bindataTplFbs10headergotplBytes() ([]byte, error) {
	return bindataRead(
		_bindataTplFbs10headergotpl,
		"_tpl/fbs_10_header.go.tpl",
	)
}



func bindataTplFbs10headergotpl() (*asset, error) {
	bytes, err := bindataTplFbs10headergotplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{
		name: "_tpl/fbs_10_header.go.tpl",
		size: 212,
		md5checksum: "",
		mode: os.FileMode(420),
		modTime: time.Unix(1544556941, 0),
	}

	a := &asset{bytes: bytes, info: info}

	return a, nil
}

var _bindataTplFbs20tablegotpl = []byte(
	"\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x8f\xb1\x4a\x03\x41\x10\x86\xeb\xdb\xa7\x98\x17\x70\xd3\x27\x95\x5e" +
	"\xd4\xce\xc6\x74\x22\x64\x34\xbf\xc7\xc2\xdc\xee\xb1\x3b\x87\xc4\x61\xde\x5d\xc2\x5a\x24\xad\xdd\x0f\xf3\xf1\x31" +
	"\xdf\x66\x43\x66\xf1\x31\x6b\xd2\xb3\x3b\x55\x2c\x15\x0d\x59\x1b\x31\xb5\x94\x27\x01\xd5\xf2\x4d\x5f\xa5\xd2\xfe" +
	"\x81\x94\x3f\x04\x74\x34\x8b\x87\xcb\x7a\xe1\x19\xee\xc7\x48\xf7\xab\x16\x9a\x90\x51\x59\x71\x8a\xa1\x73\xd7\x62" +
	"\x0b\x83\xd9\x1d\x55\xce\x13\x28\x8e\x45\xd6\x39\x37\xf7\x30\x0c\x66\x87\xf2\x5c\x46\x9e\x21\x23\x37\x50\x7c\x4a" +
	"\x90\x93\xfb\xd6\xec\x15\x35\xb1\xa4\x1f\xd4\xc3\x79\x01\x45\xf7\x1d\xf5\x87\xff\x98\x2e\x45\xbe\x4c\x0f\xa1\xdf" +
	"\xc6\x22\x82\x4f\x4d\x25\xdf\x06\xcd\xab\x68\x5a\x7a\x50\xfb\x77\xd1\x8d\xdd\xc2\xb0\x67\xe5\xed\xdb\x55\xea\xfb" +
	"\x2e\x78\xf8\x0d\x00\x00\xff\xff\xc3\x7c\xdc\xd6\x57\x01\x00\x00")

func bindataTplFbs20tablegotplBytes() ([]byte, error) {
	return bindataRead(
		_bindataTplFbs20tablegotpl,
		"_tpl/fbs_20_table.go.tpl",
	)
}



func bindataTplFbs20tablegotpl() (*asset, error) {
	bytes, err := bindataTplFbs20tablegotplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{
		name: "_tpl/fbs_20_table.go.tpl",
		size: 343,
		md5checksum: "",
		mode: os.FileMode(420),
		modTime: time.Unix(1544557066, 0),
	}

	a := &asset{bytes: bytes, info: info}

	return a, nil
}

var _bindataTplProtobuf10headergotpl = []byte(
	"\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x90\x31\x4b\x05\x31\x10\x84\xfb\xfb\x15\x4b\xb0\xf5\x52\xd8\xf9\xb8" +
	"\xc2\xce\x4e\xc1\x1f\x20\xfb\xee\xd6\x35\xbc\x5c\x36\x6e\xf6\x44\x0d\xf9\xef\x62\xe4\xe9\x81\x60\x13\x66\x60\xe6" +
	"\x63\xb2\xde\xc3\xcd\x66\x02\x4c\x89\x14\x8d\x16\x78\x0d\x08\x1c\xec\x79\x3b\x8e\xb3\xac\x7e\x16\xa5\x62\xa2\x14" +
	"\xc4\xe7\x13\xfb\xf2\x12\xfd\xb2\x46\xa6\x34\x94\xf7\x64\xf8\x06\x13\xb8\xac\x62\x72\xe5\x0e\x43\xc6\xf9\x84\x4c" +
	"\x50\xeb\x78\xff\x2d\x5b\x3b\x0c\x61\xcd\xa2\x06\x6e\x87\x65\x61\xf1\xbd\x76\xdc\x9e\xba\xeb\xa6\xab\xb1\x4b\xb7" +
	"\xeb\x89\x70\xa4\xdf\xb8\x85\x95\x8a\xe1\x9a\xff\x26\xff\x19\x6e\xa2\xc8\xe4\xd3\x16\x63\x7f\x7e\xca\x92\x2d\x48" +
	"\x02\x96\xc7\xf3\xfc\x09\xdc\xfe\x07\xee\x30\xd4\xaa\x98\x98\xe0\x42\xb2\x15\xb8\x9e\x60\x7c\x20\x0d\x18\xc3\x07" +
	"\xe9\x2d\xe1\x42\x7a\xd7\x29\x05\x2e\x5b\x3b\x13\x6b\xed\xf1\xaf\x13\xd4\x4a\x69\x69\x6d\xf8\x0c\x00\x00\xff\xff" +
	"\x81\xcf\xcc\x5d\x70\x01\x00\x00")

func bindataTplProtobuf10headergotplBytes() ([]byte, error) {
	return bindataRead(
		_bindataTplProtobuf10headergotpl,
		"_tpl/protobuf_10_header.go.tpl",
	)
}



func bindataTplProtobuf10headergotpl() (*asset, error) {
	bytes, err := bindataTplProtobuf10headergotplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{
		name: "_tpl/protobuf_10_header.go.tpl",
		size: 368,
		md5checksum: "",
		mode: os.FileMode(420),
		modTime: time.Unix(1544460649, 0),
	}

	a := &asset{bytes: bytes, info: info}

	return a, nil
}

var _bindataTplProtobuf20tablegotpl = []byte(
	"\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xa4\x90\xcf\x4a\x03\x31\x10\xc6\xcf\xcd\x53\x0c\x3d\xe9\xc1\x14\xcf\xa5" +
	"\x07\xdd\xaa\x37\x11\xec\x4d\x84\x8e\xed\x67\x08\x24\x99\x25\xc9\x22\x35\xcc\xbb\x4b\x2c\x45\xf6\xec\x6d\x66\xf8" +
	"\xfe\x0c\xbf\xd5\x8a\x5a\xb3\x0f\xa9\xfa\x7a\x52\xa5\x8c\x31\xa3\x20\xd5\x42\x4c\xc5\x27\x17\x40\x59\xbe\xe8\x53" +
	"\x32\x6d\xef\xa9\xf2\x47\x00\xed\x5b\xb3\xbb\x3e\x3d\x73\x84\xea\xde\xd2\xdd\x54\x85\x1c\x12\x32\x57\x1c\xad\x89" +
	"\x28\x85\x1d\x66\xd1\xcd\x2c\x5a\xbb\xa1\xcc\xc9\x81\xec\x20\x61\x8a\xa9\xa8\xf6\xeb\x2b\xb2\xe7\xe0\xbf\x91\x77" +
	"\xa7\x11\x64\xbb\xba\xd9\x47\x8f\x70\x54\xa5\x4d\x5f\x5e\xa4\xa8\xd2\xdb\x95\x13\x27\x63\x96\x2a\xf6\x30\x95\x2a" +
	"\x31\x71\xc4\xf5\x66\xd9\xda\x4e\x9e\x64\xe0\x88\x30\x70\x01\x5d\xcc\x4b\xea\x9d\x7f\xf9\xc3\xaf\xe9\xd2\xf2\xbe" +
	"\x3e\xff\x84\x74\x54\x35\x6a\xcc\x99\xc6\x20\x21\xe0\x50\xbd\xa4\x39\x91\x38\x85\xea\xc7\x33\x91\xf2\x0f\x24\xb3" +
	"\xfc\x66\x16\x19\x23\xba\x68\x86\x6b\xcb\x95\x69\x43\xb7\x6b\xa3\xe6\x27\x00\x00\xff\xff\x31\x69\x70\x22\xa5\x01" +
	"\x00\x00")

func bindataTplProtobuf20tablegotplBytes() ([]byte, error) {
	return bindataRead(
		_bindataTplProtobuf20tablegotpl,
		"_tpl/protobuf_20_table.go.tpl",
	)
}



func bindataTplProtobuf20tablegotpl() (*asset, error) {
	bytes, err := bindataTplProtobuf20tablegotplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{
		name: "_tpl/protobuf_20_table.go.tpl",
		size: 421,
		md5checksum: "",
		mode: os.FileMode(420),
		modTime: time.Unix(1544555797, 0),
	}

	a := &asset{bytes: bytes, info: info}

	return a, nil
}


//
// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
//
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
}

//
// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
// nolint: deadcode
//
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

//
// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or could not be loaded.
//
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
}

//
// AssetNames returns the names of the assets.
// nolint: deadcode
//
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

//
// _bindata is a table, holding each asset generator, mapped to its name.
//
var _bindata = map[string]func() (*asset, error){
	"_tpl/10_tables.go.tpl":             bindataTpl10tablesgotpl,
	"_tpl/20_entity.go.tpl":             bindataTpl20entitygotpl,
	"_tpl/30_collection_methods.go.tpl": bindataTpl30collectionmethodsgotpl,
	"_tpl/40_binary.go.tpl":             bindataTpl40binarygotpl,
	"_tpl/90_test.go.tpl":               bindataTpl90testgotpl,
	"_tpl/fbs_10_header.go.tpl":         bindataTplFbs10headergotpl,
	"_tpl/fbs_20_table.go.tpl":          bindataTplFbs20tablegotpl,
	"_tpl/protobuf_10_header.go.tpl":    bindataTplProtobuf10headergotpl,
	"_tpl/protobuf_20_table.go.tpl":     bindataTplProtobuf20tablegotpl,
}

//
// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
//
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, &os.PathError{
					Op: "open",
					Path: name,
					Err: os.ErrNotExist,
				}
			}
		}
	}
	if node.Func != nil {
		return nil, &os.PathError{
			Op: "open",
			Path: name,
			Err: os.ErrNotExist,
		}
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}


type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{Func: nil, Children: map[string]*bintree{
	"_tpl": {Func: nil, Children: map[string]*bintree{
		"10_tables.go.tpl": {Func: bindataTpl10tablesgotpl, Children: map[string]*bintree{}},
		"20_entity.go.tpl": {Func: bindataTpl20entitygotpl, Children: map[string]*bintree{}},
		"30_collection_methods.go.tpl": {Func: bindataTpl30collectionmethodsgotpl, Children: map[string]*bintree{}},
		"40_binary.go.tpl": {Func: bindataTpl40binarygotpl, Children: map[string]*bintree{}},
		"90_test.go.tpl": {Func: bindataTpl90testgotpl, Children: map[string]*bintree{}},
		"fbs_10_header.go.tpl": {Func: bindataTplFbs10headergotpl, Children: map[string]*bintree{}},
		"fbs_20_table.go.tpl": {Func: bindataTplFbs20tablegotpl, Children: map[string]*bintree{}},
		"protobuf_10_header.go.tpl": {Func: bindataTplProtobuf10headergotpl, Children: map[string]*bintree{}},
		"protobuf_20_table.go.tpl": {Func: bindataTplProtobuf20tablegotpl, Children: map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
