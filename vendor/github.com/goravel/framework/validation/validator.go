package validation

import (
	"net/url"
	"reflect"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/gookit/validate"
	"github.com/spf13/cast"

	httpvalidate "github.com/goravel/framework/contracts/validation"
	"github.com/goravel/framework/support/carbon"
	"github.com/goravel/framework/support/maps"
)

func init() {
	validate.Config(func(opt *validate.GlobalOption) {
		opt.StopOnError = false
		opt.SkipOnEmpty = true
		opt.FieldTag = "form"
		opt.RestoreRequestBody = true
	})
}

type Validator struct {
	instance *validate.Validation
	data     validate.DataFace
}

func NewValidator(instance *validate.Validation, data validate.DataFace) *Validator {
	instance.Validate()

	return &Validator{instance: instance, data: data}
}

func (v *Validator) Bind(ptr any) error {
	// Don't bind if there are errors
	if v.Fails() {
		return nil
	}

	// SafeData only contains the data that is defined in the rules,
	// we want user can the original data that is not defined in the rules,
	// so that user doesn't need to define rules for all fields.
	data := v.instance.SafeData()
	prtType := reflect.TypeOf(ptr)
	if prtType.Kind() == reflect.Ptr {
		prtType = prtType.Elem()
	}

	if formData, ok := v.data.(*validate.FormData); ok {
		if values, ok := v.data.Src().(url.Values); ok {
			for key, value := range values {
				if _, exist := data[key]; !exist {
					data[key] = value[0]
				}
			}

			for key, value := range formData.Files {
				if _, exist := data[key]; !exist {
					for i := 0; i < prtType.NumField(); i++ {
						field := prtType.Field(i)
						if field.Tag.Get("form") == key {
							if field.Type.Kind() == reflect.Slice {
								data[key] = value
							} else {
								data[key] = value[0]
							}
						}
					}
				}
			}
		}
	} else if _, ok := v.data.(*validate.MapData); ok {
		values := v.data.Src().(map[string]any)
		for key, value := range values {
			if _, exist := data[key]; !exist {
				data[key] = value
			}
		}
	} else {
		if srcMap := maps.FromStruct(v.data.Src()); len(srcMap) > 0 {
			for key, value := range srcMap {
				if _, exist := data[key]; !exist {
					data[key] = value
				}
			}
		}
	}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName:    "form",
		Result:     &ptr,
		DecodeHook: v.castValue(),
		Squash:     true,
	})
	if err != nil {
		return err
	}

	return decoder.Decode(data)
}

func (v *Validator) Errors() httpvalidate.Errors {
	if len(v.instance.Errors) == 0 {
		return nil
	}

	return NewErrors(v.instance.Errors)
}

func (v *Validator) Fails() bool {
	return v.instance.IsFail()
}

func (v *Validator) castValue() mapstructure.DecodeHookFunc {
	return func(from reflect.Value, to reflect.Value) (any, error) {
		var (
			err error

			castedValue = from.Interface()
		)

		switch to.Kind() {
		case reflect.String:
			castedValue = cast.ToString(from.Interface())
		case reflect.Int:
			castedValue, err = cast.ToIntE(from.Interface())
		case reflect.Int8:
			castedValue, err = cast.ToInt8E(from.Interface())
		case reflect.Int16:
			castedValue, err = cast.ToInt16E(from.Interface())
		case reflect.Int32:
			castedValue, err = cast.ToInt32E(from.Interface())
		case reflect.Int64:
			castedValue, err = cast.ToInt64E(from.Interface())
		case reflect.Uint:
			castedValue, err = cast.ToUintE(from.Interface())
		case reflect.Uint8:
			castedValue, err = cast.ToUint8E(from.Interface())
		case reflect.Uint16:
			castedValue, err = cast.ToUint16E(from.Interface())
		case reflect.Uint32:
			castedValue, err = cast.ToUint32E(from.Interface())
		case reflect.Uint64:
			castedValue, err = cast.ToUint64E(from.Interface())
		case reflect.Bool:
			castedValue, err = cast.ToBoolE(from.Interface())
		case reflect.Float32:
			castedValue, err = cast.ToFloat32E(from.Interface())
		case reflect.Float64:
			castedValue, err = cast.ToFloat64E(from.Interface())
		case reflect.Slice, reflect.Array:
			switch to.Type().Elem().Kind() {
			case reflect.String:
				castedValue, err = cast.ToStringSliceE(from.Interface())
			case reflect.Int:
				castedValue, err = cast.ToIntSliceE(from.Interface())
			case reflect.Bool:
				castedValue, err = cast.ToBoolSliceE(from.Interface())
			default:
				castedValue, err = cast.ToSliceE(from.Interface())
			}
		case reflect.Struct:
			switch to.Type() {
			case reflect.TypeOf(carbon.Carbon{}):
				castedValue = castCarbon(from, nil)
			case reflect.TypeOf(carbon.DateTime{}):
				castedValue = castCarbon(from, func(c *carbon.Carbon) any {
					return carbon.NewDateTime(c)
				})
			case reflect.TypeOf(carbon.DateTimeMilli{}):
				castedValue = castCarbon(from, func(c *carbon.Carbon) any {
					return carbon.NewDateTimeMilli(c)
				})
			case reflect.TypeOf(carbon.DateTimeMicro{}):
				castedValue = castCarbon(from, func(c *carbon.Carbon) any {
					return carbon.NewDateTimeMicro(c)
				})
			case reflect.TypeOf(carbon.DateTimeNano{}):
				castedValue = castCarbon(from, func(c *carbon.Carbon) any {
					return carbon.NewDateTimeNano(c)
				})
			case reflect.TypeOf(carbon.Date{}):
				castedValue = castCarbon(from, func(c *carbon.Carbon) any {
					return carbon.NewDate(c)
				})
			case reflect.TypeOf(carbon.DateMilli{}):
				castedValue = castCarbon(from, func(c *carbon.Carbon) any {
					return carbon.NewDateMilli(c)
				})
			case reflect.TypeOf(carbon.DateMicro{}):
				castedValue = castCarbon(from, func(c *carbon.Carbon) any {
					return carbon.NewDateMicro(c)
				})
			case reflect.TypeOf(carbon.DateNano{}):
				castedValue = castCarbon(from, func(c *carbon.Carbon) any {
					return carbon.NewDateNano(c)
				})
			case reflect.TypeOf(carbon.Timestamp{}):
				castedValue = castCarbon(from, func(c *carbon.Carbon) any {
					return carbon.NewTimestamp(c)
				})
			case reflect.TypeOf(carbon.TimestampMilli{}):
				castedValue = castCarbon(from, func(c *carbon.Carbon) any {
					return carbon.NewTimestampMilli(c)
				})
			case reflect.TypeOf(carbon.TimestampMicro{}):
				castedValue = castCarbon(from, func(c *carbon.Carbon) any {
					return carbon.NewTimestampMicro(c)
				})
			case reflect.TypeOf(carbon.TimestampNano{}):
				castedValue = castCarbon(from, func(c *carbon.Carbon) any {
					return carbon.NewTimestampNano(c)
				})
			case reflect.TypeOf(time.Time{}):
				castedValue = castCarbon(from, func(c *carbon.Carbon) any {
					return c.StdTime()
				})
			}
		default:
			castedValue = from.Interface()
		}

		// Only return casted value if there was no error
		if err == nil {
			return castedValue, nil
		}

		return from.Interface(), nil
	}
}

func castCarbon(from reflect.Value, transfrom func(carbon *carbon.Carbon) any) any {
	var c *carbon.Carbon

	switch len(cast.ToString(from.Interface())) {
	case 10:
		fromInt64, err := cast.ToInt64E(from.Interface())
		if err != nil {
			c = carbon.Parse(cast.ToString(from.Interface()))
		}
		if fromInt64 > 0 {
			c = carbon.FromTimestamp(fromInt64)
		}
	case 13:
		fromInt64, err := cast.ToInt64E(from.Interface())
		if err != nil {
			c = carbon.ParseByFormat(cast.ToString(from.Interface()), "Y-m-d H")
		}
		if fromInt64 > 0 {
			c = carbon.FromTimestampMilli(fromInt64)
		}
	case 16:
		fromInt64, err := cast.ToInt64E(from.Interface())
		if err != nil {
			c = carbon.ParseByFormat(cast.ToString(from.Interface()), "Y-m-d H:i")
		}
		if fromInt64 > 0 {
			c = carbon.FromTimestampMicro(fromInt64)
		}
	case 19:
		fromInt64, err := cast.ToInt64E(from.Interface())
		if err != nil {
			c = carbon.Parse(cast.ToString(from.Interface()))
		}

		if fromInt64 > 0 {
			c = carbon.FromTimestampNano(fromInt64)
		}
	default:
		c = carbon.Parse(cast.ToString(from.Interface()))
	}

	if transfrom != nil {
		return transfrom(c)
	}

	return c
}
