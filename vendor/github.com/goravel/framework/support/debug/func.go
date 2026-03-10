package debug

import (
	"reflect"
	"runtime"
	"strings"
)

type FuncInfo struct {
	File      string
	Name      string
	pkgName   string
	pkgPath   string
	shortName string
	Line      int
}

func (f *FuncInfo) PackageName() string {
	return f.pkgName
}

func (f *FuncInfo) PackagePath() string {
	return f.pkgPath
}

func (f *FuncInfo) ShortName() string {
	return f.shortName
}

func GetFuncInfo(f any) FuncInfo {
	var info FuncInfo
	v := reflect.ValueOf(f)
	if v.Kind() != reflect.Func {
		return info
	}

	// get function info
	ptr := v.Pointer()
	fnObj := runtime.FuncForPC(ptr)
	if fnObj == nil {
		return info
	}

	// get function file and line
	info.Name = fnObj.Name()
	info.File, info.Line = fnObj.FileLine(ptr)

	// split get pkg and func name
	if lastSlash := strings.LastIndex(info.Name, "/"); lastSlash >= 0 {
		info.pkgPath = info.Name[:lastSlash+1]
		info.pkgName, info.shortName, _ = strings.Cut(info.Name[lastSlash+1:], ".")
		info.pkgPath += info.pkgName
	}

	return info
}
