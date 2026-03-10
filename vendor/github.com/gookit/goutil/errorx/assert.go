package errorx

import (
	"errors"
	"fmt"

	"github.com/gookit/goutil/arrutil"
	"github.com/gookit/goutil/comdef"
	"github.com/gookit/goutil/internal/comfunc"
)

// IsTrue assert result is true, otherwise will return error
func IsTrue(result bool, fmtAndArgs ...any) error {
	if !result {
		return errors.New(formatErrMsg("result should be True", fmtAndArgs))
	}
	return nil
}

// IsFalse assert result is false, otherwise will return error
func IsFalse(result bool, fmtAndArgs ...any) error {
	if result {
		return errors.New(formatErrMsg("result should be False", fmtAndArgs))
	}
	return nil
}

// IsIn value should be in the list, otherwise will return error
func IsIn[T comdef.ScalarType](value T, list []T, fmtAndArgs ...any) error {
	if arrutil.NotIn(value, list) {
		var errMsg string
		if len(fmtAndArgs) > 0 {
			errMsg = comfunc.FormatWithArgs(fmtAndArgs)
		} else {
			errMsg = fmt.Sprintf("value should be in the %v", list)
		}
		return errors.New(errMsg)
	}
	return nil
}

// NotIn value should not be in the list, otherwise will return error
func NotIn[T comdef.ScalarType](value T, list []T, fmtAndArgs ...any) error {
	if arrutil.In(value, list) {
		var errMsg string
		if len(fmtAndArgs) > 0 {
			errMsg = comfunc.FormatWithArgs(fmtAndArgs)
		} else {
			errMsg = fmt.Sprintf("value should not be in the %v", list)
		}
		return errors.New(errMsg)
	}
	return nil
}

func formatErrMsg(defMsg string, fmtAndArgs []any) string {
	if len(fmtAndArgs) > 0 {
		return comfunc.FormatWithArgs(fmtAndArgs)
	}
	return defMsg
}
