package byteutil

import (
	"fmt"
	"strconv"
	"time"
	"unsafe"

	"github.com/gookit/goutil/comdef"
)

// StrOrErr convert to string, return empty string on error.
func StrOrErr(bs []byte, err error) (string, error) {
	if err != nil {
		return "", err
	}
	return string(bs), err
}

// SafeString convert to string, return empty string on error.
func SafeString(bs []byte, err error) string {
	if err != nil {
		return ""
	}
	return string(bs)
}

// String unsafe convert bytes to string
func String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// ToString convert bytes to string
func ToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// ToBytes convert any value to []byte. return error on convert failed.
func ToBytes(v any) ([]byte, error) {
	return ToBytesWithFunc(v, nil)
}

// SafeBytes convert any value to []byte. use fmt.Sprint() on convert failed.
func SafeBytes(v any) []byte {
	bs, _ := ToBytesWithFunc(v, func(v any) ([]byte, error) {
		return []byte(fmt.Sprint(v)), nil
	})
	return bs
}

// ToBytesFunc convert any value to []byte
type ToBytesFunc = func(v any) ([]byte, error)

// ToBytesWithFunc convert any value to []byte with custom fallback func.
//
// refer the strutil.ToStringWithFunc
//
// On not convert:
//   - If usrFn is nil, will return comdef.ErrConvType.
//   - If usrFn is not nil, will call it to convert.
func ToBytesWithFunc(v any, usrFn ToBytesFunc) ([]byte, error) {
	if v == nil {
		return nil, nil
	}

	switch val := v.(type) {
	case []byte:
		return val, nil
	case string:
		return []byte(val), nil
	case int:
		return []byte(strconv.Itoa(val)), nil
	case int8:
		return []byte(strconv.Itoa(int(val))), nil
	case int16:
		return []byte(strconv.Itoa(int(val))), nil
	case int32: // same as `rune`
		return []byte(strconv.Itoa(int(val))), nil
	case int64:
		return []byte(strconv.FormatInt(val, 10)), nil
	case uint:
		return []byte(strconv.FormatUint(uint64(val), 10)), nil
	case uint8:
		return []byte(strconv.FormatUint(uint64(val), 10)), nil
	case uint16:
		return []byte(strconv.FormatUint(uint64(val), 10)), nil
	case uint32:
		return []byte(strconv.FormatUint(uint64(val), 10)), nil
	case uint64:
		return []byte(strconv.FormatUint(val, 10)), nil
	case float32:
		return []byte(strconv.FormatFloat(float64(val), 'f', -1, 32)), nil
	case float64:
		return []byte(strconv.FormatFloat(val, 'f', -1, 64)), nil
	case bool:
		return []byte(strconv.FormatBool(val)), nil
	case time.Duration:
		return []byte(strconv.FormatInt(int64(val), 10)), nil
	case fmt.Stringer:
		return []byte(val.String()), nil
	case error:
		return []byte(val.Error()), nil
	default:
		if usrFn == nil {
			return nil, comdef.ErrConvType
		}
		return usrFn(val)
	}
}
