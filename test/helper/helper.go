package helper

import (
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
)

func init() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("error getting runtime information")
	}

	RepoBaseDir = func(fpath string) string {
		for fpath != "." {
			fpath = filepath.Dir(fpath)
			if _, err := os.Stat(filepath.Join(fpath, "go.mod")); err == nil {
				return fpath
			}
		}
		panic("could not find go.mod")
	}(file)

	RepoBasePkg = func(fpath string) string {
		type foo struct{}
		pkgPath := reflect.TypeOf(foo{}).PkgPath()
		fpath = filepath.Dir(fpath)
		for path.Base(pkgPath) == filepath.Base(fpath) && filepath.Base(fpath) != filepath.Base(RepoBaseDir) {
			pkgPath = path.Dir(pkgPath)
			fpath = filepath.Dir(fpath)
		}
		return pkgPath
	}(file)
}

var RepoBaseDir string
var RepoBasePkg string
