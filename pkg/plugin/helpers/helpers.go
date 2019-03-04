// VulcanizeDB
// Copyright © 2018 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package helpers

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/mitchellh/go-homedir"
)

func CleanPath(str string) (string, error) {
	path, err := homedir.Expand(filepath.Clean(str))
	if err != nil {
		return "", err
	}
	if strings.Contains(path, "$GOPATH") {
		env := os.Getenv("GOPATH")
		spl := strings.Split(path, "$GOPATH")[1]
		path = filepath.Join(env, spl)
	}

	return path, nil
}

func ClearFiles(files ...string) error {
	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			err = os.Remove(file)
			if err != nil {
				return err
			}
		} else if os.IsNotExist(err) {
			continue
		} else {
			return err
		}
	}

	return nil
}

func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	out, err := os.OpenFile(dst, syscall.O_CREAT|syscall.O_EXCL|os.O_WRONLY, os.FileMode(0666)) // Doesn't overwrite files
	if err != nil {
		in.Close()
		return err
	}

	_, err = io.Copy(out, in)
	in.Close()
	out.Close()
	return err
}

func CopyDir(src string, dst string, excludeRecursiveDir string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := filepath.Join(src, fd.Name())
		dstfp := filepath.Join(dst, fd.Name())

		if fd.IsDir() {
			if fd.Name() != excludeRecursiveDir {
				err = CopyDir(srcfp, dstfp, "")
				if err != nil {
					return err
				}
			}
		} else {
			err = CopyFile(srcfp, dstfp)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
