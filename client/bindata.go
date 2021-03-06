// Code generated by go-bindata.
// sources:
// cert/tunnel.crt
// DO NOT EDIT!

package main

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
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
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
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _certTunnelCrt = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x64\x95\xc7\xd2\xab\x3a\x16\x85\xe7\x3c\x45\xcf\x5d\x5d\x04\x83\x31\x43\x89\x9c\x73\x9c\x91\x31\x60\x4c\x30\xf1\xe9\xfb\xbf\xe7\x74\xdd\xae\xd3\x57\x23\xd5\x27\xa9\xb4\x6b\xed\x55\x7b\xfd\xfb\xaf\x05\x79\x51\x36\xfe\xc5\xf2\x8e\x27\x0b\x32\x0b\x3c\xfe\x17\x45\x74\x59\x16\x14\x8e\x65\xc1\xbe\xb3\xac\xcd\x3a\xc6\x37\x2a\xb8\xe4\x2b\x45\x1e\x30\x60\xdd\x4d\x4d\xf7\x12\x99\x1d\x83\xc0\xf6\x05\xc0\x41\x5f\xb7\x97\x9d\xb5\x63\x2e\xb0\x6d\x91\xdf\x15\x0e\xf1\x5a\xde\xd6\x01\x29\x02\xdc\xe7\x59\xa8\x4b\xf6\x3b\x18\xd3\x77\xbf\x26\xd7\xff\x38\x3c\xfe\xe0\x8d\x0e\x9f\xbf\xef\x37\x7a\x8c\xb8\x21\x85\x25\x91\xb2\x26\x91\x5d\x07\x44\xdf\x25\x44\x8f\xe5\x27\xb4\x0b\xa9\xab\x3d\xc9\xe9\x74\x48\x46\x9c\xc7\x13\x3a\x27\xef\xa6\xc7\x9f\x06\x27\xe3\x86\xf0\xf9\x61\xf2\x2f\x86\x18\xed\xdf\x70\x0f\x38\x5e\xd3\x41\xf7\xfb\xd7\x46\x67\x6d\x8c\x3c\x78\x0e\x98\xb0\x36\x02\x08\x6a\x0f\x62\x4a\x9f\x86\xe3\x98\xbd\xf3\xbf\x39\x02\x41\xfe\xc7\x81\xec\x81\xea\xf7\x83\x8f\x27\xf2\xfd\x5a\x88\xc1\x99\xbd\x03\x4c\x16\x8a\x31\x11\x8b\xb1\x90\xf4\xda\x97\x1c\x4a\xe6\x0f\x2c\xf9\xd1\x0e\x01\x2f\x79\xe7\xec\x58\x51\x3f\x89\xdc\x6c\xb9\x01\x6c\x1e\x42\x1b\x70\x75\x2d\x5b\xe0\x2f\x71\xeb\x0f\xfb\xb3\x87\x40\x6f\x86\x8e\x56\x17\x83\xd5\xea\xed\x98\xe6\x19\xc4\x58\x37\x84\x2b\x32\xba\x26\xa7\x5e\x4a\xa8\x3e\xd7\xef\xe7\x69\x3c\xc7\x35\x4c\xe3\x34\xb4\xd3\xc0\xd8\xc5\x47\x6b\x26\x9f\x14\xba\xb6\x94\xf2\xea\x37\x3f\xd6\xf1\x50\x3e\x58\x8d\x2f\xc6\x91\x0f\x1b\xcf\xb4\x3d\x8e\xb0\xed\x1e\x5c\x94\x45\x9f\xd7\x5c\x12\x0b\x2f\x0f\x6c\xf6\xba\xf6\xe6\x6d\x17\x91\x37\x3d\x6f\x0f\x17\x77\x93\xb7\x3a\x94\xde\xfa\x95\x2f\xd2\x68\xa6\x46\xc5\xeb\x74\x82\xd0\x1a\x37\x3b\xc5\x28\xc4\xa0\x07\x01\xd4\x15\xbe\x00\x71\x97\x1d\x2f\xa1\x31\x2f\x27\xc2\xf9\xe1\x1f\x65\xa3\x2a\x3a\xcb\xb7\xa3\x13\x48\xd2\xdb\x67\xc6\x9b\xa3\x1a\x33\x1a\x2a\xc6\x62\xbc\x47\x48\x85\xfb\xcd\x9e\x24\x44\xb4\x65\x91\x47\x39\xda\x7d\x10\x49\x93\x4c\x29\xbd\x11\x8d\xe8\x7f\x3b\x50\xc5\x1a\xb0\xa3\xea\xd9\xca\x32\x1a\x02\x56\x54\xd1\x0c\xb4\x6f\x5a\x87\x8b\x54\xe7\xee\xeb\x0a\x15\x45\x5a\xfc\x10\xe1\xd9\x12\x77\x50\xab\xe1\x8c\xcc\x99\xcf\x6d\x7a\x9a\x4a\xea\x69\x0f\x4d\x25\xc3\xdc\xbb\x6d\x98\x91\x2a\x8a\x6f\x78\xcb\x32\x8d\x23\xd9\x07\x8b\x31\x0a\x2d\x15\x93\xbb\x1d\x17\x07\x99\xf8\x05\xb2\x0c\xc0\x79\xdd\x59\xc6\xab\xd1\xe1\xdc\xb3\x00\xb7\xd5\x82\xb9\x11\x65\xd9\x10\x04\x6b\x58\x2f\xec\x05\x84\xeb\x34\x36\x2b\xa9\x52\x0a\x48\xb6\x4f\xb7\xae\xab\xb9\xd4\xd1\xf7\xb7\x07\x93\x3d\x91\x36\x28\xbb\x8d\xeb\x97\x36\x5f\x43\xd5\x8e\x73\x33\x0b\x1e\xe4\xbe\xfa\xfb\x99\xd3\x36\xff\x99\xb6\x2e\xac\xe6\xc0\x26\xe1\xf6\xe3\x32\x16\x5a\xd4\xed\x79\x9b\xda\xd4\x6e\x50\x0d\xf3\x83\xc5\x42\xd4\x9b\xd3\x31\xab\x26\x25\xe5\x8d\x88\x93\x1a\x5d\x68\xdb\x93\xae\x3e\x22\x2d\x9f\x9c\xd9\x7d\xc0\x27\xf4\x24\xf8\xe1\xf3\x80\xd6\xe6\x84\x06\x8a\xe9\xb8\x35\xf3\xa6\x0d\x8b\xe7\x95\x60\xae\x86\x50\x45\x75\xf5\x4d\x50\xad\x73\x17\x0c\xaf\x32\x68\x4e\x35\x72\xfa\xcb\xc4\x8e\x9d\xa5\x63\x96\x0b\x62\x8c\x1c\x9d\x81\xed\x69\x97\x57\xbe\x74\x12\x08\xf7\x70\x89\x5e\x2a\x55\x33\x80\xcd\x12\x0a\x91\x6e\xa1\x76\x66\x54\xb8\xb9\x50\x1e\x85\x9c\xa4\x37\xda\x16\x44\xa1\x49\x00\x4c\xb8\x62\xf2\x13\x85\x8c\x7f\xaa\x7e\x15\xe7\x31\x47\x1f\x72\x15\xf4\xe8\x3b\xec\x96\xf3\xf0\x47\x34\xf4\xf6\x09\x61\xd1\x3d\xef\x50\x40\x2a\x53\x39\x7e\x84\xec\x74\xcc\xad\x00\xb5\x0e\x01\xe0\xff\xdf\xe2\xc2\x7f\x2d\x0e\x81\x54\x6b\x94\x8a\x45\xaf\xf7\x19\x21\x41\x58\x76\xaa\xea\x37\xbd\x78\x8b\xad\x4b\xc5\xb3\xb4\x52\xf8\xee\xc0\x0a\x21\x2c\x08\x34\xc2\x14\xc9\x42\x3d\xae\xfd\xda\x2e\xbc\x71\x68\x55\x0f\x77\x70\xe5\x9f\xbe\x31\x55\x76\x83\x96\xd6\x21\xaf\xd6\x8f\xb4\x11\x17\x05\x6a\x7f\x8e\xba\x39\xd2\x54\xfb\xe0\x5a\xa1\xba\x09\x2e\x46\x12\x7e\xfd\x3d\xc5\x61\xf4\x8e\xef\xc5\xde\x5b\x03\x95\x98\xf6\x2c\x6c\xac\xf9\x08\xe1\x55\xa2\xb3\xea\x7a\x08\xc7\xa0\x5c\x35\x1b\x59\x94\x6c\x5b\x21\x99\xb2\x0d\xe7\xaf\xa7\x3e\x04\xce\x8c\x9a\x92\x50\x30\x4e\x01\x15\x6d\x8b\xdd\x14\x68\xf9\x3c\x4b\x49\xd7\x15\x0f\x0f\x1f\x9e\xb7\x3e\xbd\x54\x5e\xb3\x11\x2a\x5f\xe9\xcf\x74\xd2\xa5\xc9\x8f\x44\x42\xd1\xc9\xa4\xcb\x45\x01\xd0\x3a\xbb\xef\xa9\xf9\xd9\x06\x7c\xf8\xba\xea\xb5\x50\x0d\xdc\xd8\x17\xf3\x33\x1d\x2e\xa6\x37\x63\x3e\x3a\x53\xf1\x38\x54\x05\x71\x89\x5d\x6f\xa4\x75\x57\x64\x91\xd6\x4f\x45\xaf\x48\xd6\x25\xe6\x02\x7d\x8b\x3a\x7f\x11\xcc\x39\xd9\x97\xfb\xd9\xc6\xdb\x91\x28\x37\x9e\x68\xcf\xb1\x34\x7c\xa3\x02\xd7\x52\xbc\x77\x93\x2f\xdf\x48\x70\x90\x56\x92\x88\xf7\x6b\x79\x99\xbb\x1b\xef\x31\x11\x31\x78\x19\x04\x97\x66\xec\x4c\xb5\xb1\x12\x61\x69\xe6\xd6\x66\xee\xbd\xf0\x3c\xd9\xd4\xe7\x53\x39\x63\x68\x44\x98\x41\x3e\xa8\x8b\xde\x91\x79\x8c\xd1\xb2\x33\x3b\xed\xfd\xd3\x9e\x44\x5e\xb8\xd3\x76\x5f\x47\x78\x97\xcf\x66\x38\xef\x24\x4b\x35\x06\xa5\x15\xdc\x92\x29\xdc\xac\xcf\x2e\x91\xdf\xee\xb8\xcb\x1f\xb0\xbd\xc2\x39\xbc\xe0\x1d\x21\xc4\x67\xe6\xb4\xf3\xc2\x04\x65\xf9\xad\xcc\x6c\xb3\x27\xe2\x49\x6a\x14\xd1\x99\x4c\x42\x77\x83\x5c\x76\x90\xc8\x0c\x3c\x03\x33\x3a\x57\x9d\x73\x6d\x5f\xc8\xa3\xdb\x09\xfb\xdd\x5d\xcc\x71\x40\x18\xfa\x95\xa2\x3a\x95\xa7\xab\xa5\x98\x38\x5e\x54\x9d\xec\xfb\x03\x77\xdc\xf3\xdb\x43\x6f\xc5\xde\xfa\x3e\x68\x8f\x21\xa9\x71\x12\x8e\x41\xa1\xd6\x8f\x3b\x8e\x8f\x51\xf5\x4a\x22\x67\x27\x37\xe8\x91\x63\x5d\xd9\xed\xf5\xf8\x08\x81\xd7\xa5\x4c\x5a\xdc\x66\xa6\x7f\x93\x55\xc1\x94\xad\x23\xe4\x43\x3f\x69\x1c\xb4\x88\x04\x96\xb5\xec\xbb\x9c\xc6\x87\x8e\xc1\xcd\x77\xd2\x29\x9a\xe6\x43\x71\xc3\x1d\x19\xf5\x69\x1a\xb2\x74\x8f\x93\xe8\xac\xef\x53\x00\x1c\xbd\x5f\x77\xe3\x6e\x08\x9c\x12\x48\xcf\x0d\xf9\x95\x80\xbc\xc1\xfd\x33\x15\xff\x13\x00\x00\xff\xff\x97\xb1\xa3\xc1\x32\x07\x00\x00")

func certTunnelCrtBytes() ([]byte, error) {
	return bindataRead(
		_certTunnelCrt,
		"cert/tunnel.crt",
	)
}

func certTunnelCrt() (*asset, error) {
	bytes, err := certTunnelCrtBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "cert/tunnel.crt", size: 1842, mode: os.FileMode(420), modTime: time.Unix(1460617929, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"cert/tunnel.crt": certTunnelCrt,
}

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
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
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

var _bintree = &bintree{nil, map[string]*bintree{
	"cert": &bintree{nil, map[string]*bintree{
		"tunnel.crt": &bintree{certTunnelCrt, map[string]*bintree{}},
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
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
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
