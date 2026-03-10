# Validate

[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/gookit/validate)](https://github.com/gookit/validate)
[![GoDoc](https://pkg.go.dev/badge/github.com/gookit/validate.svg)](https://pkg.go.dev/github.com/gookit/validate)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gookit/goutil?style=flat-square)
[![Coverage Status](https://coveralls.io/repos/github/gookit/validate/badge.svg?branch=master)](https://coveralls.io/github/gookit/validate?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/gookit/validate)](https://goreportcard.com/report/github.com/gookit/validate)
[![Actions Status](https://github.com/gookit/validate/workflows/Unit-Tests/badge.svg)](https://github.com/gookit/validate/actions)

`validate` is a generic Go data validate and filter tool library.

- Support quick validate `Map`, `Struct`, `Request`(`Form`, `JSON`, `url.Values`, `UploadedFile`) data
  - Validating `http.Request` automatically collects data based on the request `Content-Type` value
  - Supports checking each child value in a slice. eg: `v.StringRule("tags.*", "required|string")`
- Support filter/sanitize/convert data before validate
- Support add custom filter/validator func
- Support scene settings, verify different fields in different scenes
- Support custom error messages, field translates.
  - Can use `message`, `label` tags in struct
- Customizable i18n aware error messages, built in `en`, `zh-CN`, `zh-TW`
- Built-in common data type filter/converter. see [Built In Filters](#built-in-filters)
- Many commonly used validators have been built in(**> 70**), see [Built In Validators](#built-in-validators)
- Can use `validate` in any frameworks, such as Gin, Echo, Chi and more
- Supports direct use of rules to validate value. eg: `validate.Val("xyz@mail.com", "required|email")`

## [中文说明](README.zh-CN.md)

中文说明请查看 **[README.zh-CN](README.zh-CN.md)**

## Go Doc

- [godoc for gopkg](https://pkg.go.dev/gopkg.in/gookit/validate.v1)
- [godoc for github](https://pkg.go.dev/github.com/gookit/validate)

## Validate Struct

Use the `validate` tag of the structure, you can quickly config a structure.

### Config the struct use tags

Field translations and error messages for structs can be quickly configured using the `message` and `label` tags.

- Support configuration field mapping through structure tag, read the value of `json` tag by default
- Support configuration error message via structure's `message` tag
- Support configuration field translation via structure's `label` tag

```go
package main

import (
	"fmt"
	"time"

	"github.com/gookit/validate"
)

// UserForm struct
type UserForm struct {
	Name     string    `validate:"required|min_len:7" message:"required:{field} is required" label:"User Name"`
	Email    string    `validate:"email" message:"email is invalid" label:"User Email"`
	Age      int       `validate:"required|int|min:1|max:99" message:"int:age must int|min:age min value is 1"`
	CreateAt int       `validate:"min:1"`
	Safe     int       `validate:"-"`
	UpdateAt time.Time `validate:"required" message:"update time is required"`
	Code     string    `validate:"customValidator"`
	// ExtInfo nested struct
	ExtInfo struct{
		Homepage string `validate:"required" label:"Home Page"`
		CityName string
	} `validate:"required" label:"Home Page"`
}

// CustomValidator custom validator in the source struct.
func (f UserForm) CustomValidator(val string) bool {
	return len(val) == 4
}
```

### Config validate use struct methods

`validate` provides extended functionality:

The struct can implement three interfaces methods, which is convenient to do some customization:

- `ConfigValidation(v *Validation)` will be called after the validator instance is created
- `Messages() map[string]string` can customize the validator error message
- `Translates() map[string]string` can customize field translation

```go
package main

import (
	"fmt"
	"time"

	"github.com/gookit/validate"
)

// UserForm struct
type UserForm struct {
	Name     string    `validate:"required|min_len:7"`
	Email    string    `validate:"email"`
	Age      int       `validate:"required|int|min:1|max:99"`
	CreateAt int       `validate:"min:1"`
	Safe     int       `validate:"-"`
	UpdateAt time.Time `validate:"required"`
	Code     string    `validate:"customValidator"`
	// ExtInfo nested struct
	ExtInfo struct{
		Homepage string `validate:"required"`
		CityName string
	} `validate:"required"`
}

// CustomValidator custom validator in the source struct.
func (f UserForm) CustomValidator(val string) bool {
	return len(val) == 4
}

// ConfigValidation config the Validation
// eg:
// - define validate scenes
func (f UserForm) ConfigValidation(v *validate.Validation) {
	v.WithScenes(validate.SValues{
		"add":    []string{"ExtInfo.Homepage", "Name", "Code"},
		"update": []string{"ExtInfo.CityName", "Name"},
	})
}

// Messages you can custom validator error messages. 
func (f UserForm) Messages() map[string]string {
	return validate.MS{
		"required": "oh! the {field} is required",
		"email": "email is invalid",
		"Name.required": "message for special field",
		"Age.int": "age must int",
		"Age.min": "age min value is 1",
	}
}

// Translates you can custom field translates. 
func (f UserForm) Translates() map[string]string {
	return validate.MS{
		"Name": "User Name",
		"Email": "User Email",
		"ExtInfo.Homepage": "Home Page",
	}
}
```

### Create and validating

Can use `validate.Struct(ptr)` quick create a validation instance. then call `v.Validate()` for validating.

```go
package main

import (
  "fmt"

  "github.com/gookit/validate"
)

func main() {
	u := &UserForm{
		Name: "inhere",
	}
	
	v := validate.Struct(u)
	// v := validate.New(u)

	if v.Validate() { // validate ok
		// do something ...
	} else {
		fmt.Println(v.Errors) // all error messages
		fmt.Println(v.Errors.One()) // returns a random error message text
		fmt.Println(v.Errors.OneError()) // returns a random error
		fmt.Println(v.Errors.Field("Name")) // returns error messages of the field 
	}
}
```

## Validate Map

You can also validate a MAP data directly.

```go
package main

import (
"fmt"

"github.com/gookit/validate"
)

func main()  {
	m := map[string]any{
		"name":  "inhere",
		"age":   100,
		"oldSt": 1,
		"newSt": 2,
		"email": "some@email.com",
		"tags": []string{"go", "php", "java"},
	}

	v := validate.Map(m)
	// v := validate.New(m)
	v.AddRule("name", "required")
	v.AddRule("name", "minLen", 7)
	v.AddRule("age", "max", 99)
	v.AddRule("age", "min", 1)
	v.AddRule("email", "email")
	
	// can also
	v.StringRule("age", "required|int|min:1|max:99")
	v.StringRule("name", "required|minLen:7")
	v.StringRule("tags", "required|slice|minlen:1")
	// feat: support check sub-item in slice
	v.StringRule("tags.*", "required|string|min_len:7")

	// v.WithScenes(map[string]string{
	//	 "create": []string{"name", "email"},
	//	 "update": []string{"name"},
	// })
	
	if v.Validate() { // validate ok
		safeData := v.SafeData()
		// do something ...
	} else {
		fmt.Println(v.Errors) // all error messages
		fmt.Println(v.Errors.One()) // returns a random error message text
	}
}
```

## Validate Request

If it is an HTTP request, you can quickly validate the data and pass the verification.
Then bind the secure data to the structure.

```go
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gookit/validate"
)

// UserForm struct
type UserForm struct {
	Name     string
	Email    string
	Age      int
	CreateAt int
	Safe     int
	UpdateAt time.Time
	Code     string
}

func main()  {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := validate.FromRequest(r)
		if err != nil {
			panic(err)
		}

		v := data.Create()
		// setting rules
		v.FilterRule("age", "int") // convert value to int
		
		v.AddRule("name", "required")
		v.AddRule("name", "minLen", 7)
		v.AddRule("age", "max", 99)
		v.StringRule("code", `required|regex:\d{4,6}`)

		if v.Validate() { // validate ok
			// safeData := v.SafeData()
			userForm := &UserForm{}
			v.BindSafeData(userForm)

			// do something ...
			fmt.Println(userForm.Name)
		} else {
			fmt.Println(v.Errors) // all error messages
			fmt.Println(v.Errors.One()) // returns a random error message text
		}
	})

	http.ListenAndServe(":8090", handler)
}
```

## Quick Method

Quick create `Validation` instance.

- `New(data any, scene ...string) *Validation`
- `Request(r *http.Request) *Validation`
- `JSON(s string, scene ...string) *Validation`
- `Struct(s any, scene ...string) *Validation`
- `Map(m map[string]any, scene ...string) *Validation`

Quick create `DataFace` instance.

- `FromMap(m map[string]any) *MapData`
- `FromStruct(s any) (*StructData, error)`
- `FromJSON(s string) (*MapData, error)`
- `FromJSONBytes(bs []byte) (*MapData, error)`
- `FromURLValues(values url.Values) *FormData`
- `FromRequest(r *http.Request, maxMemoryLimit ...int64) (DataFace, error)`

> Create `Validation` from `DataFace`

```go
d := FromMap(map[string]any{"key": "val"})
v := d.Validation()
```

### Methods In Validation

- `func (v *Validation) Validate(scene ...string) bool` Do validating and return is success.
- `func (v *Validation) ValidateE(scene ...string) Errors` Do validating and return error.

## More Usage

### Validate Error

`v.Errors` is map data, top key is field name, value is `map[string]string`.

```go
// do validating
if v.Validate() {
	return nil
}

// get errors
es := v.Errors

// check
es.Empty() // bool

// returns an random error, if no error returns nil
fmt.Println(v.Errors.OneError())
fmt.Println(v.Errors.ErrOrNil())

fmt.Println(v.Errors) // all error messages
fmt.Println(v.Errors.One()) // returns a random error message text
fmt.Println(v.Errors.Field("Name")) // returns error messages of the field 
```

**Encode to JSON**:

- `StopOnError=true`(default), will only one error

```json
{
    "field1": {
        "required": "error msg0"
    }
}
```

- if `StopOnError=false`, will get multi error

```json
{
    "field1": {
        "minLen": "error msg1",
        "required": "error msg0"
    },
    "field2": {
        "min": "error msg2"
    }
}
```

### Global Option

You can adjust some processing logic of the validator by changing the global option settings.

```go
// GlobalOption settings for validate
type GlobalOption struct {
	// FilterTag name in the struct tags.
	//
	// default: filter
	FilterTag string
	// ValidateTag in the struct tags.
	//
	// default: validate
	ValidateTag string
	// FieldTag the output field name in the struct tags.
	// it as placeholder on error message.
	//
	// default: json
	FieldTag string
	// LabelTag the display name in the struct tags.
	// use for define field translate name on error.
	//
	// default: label
	LabelTag string
	// MessageTag define error message for the field.
	//
	// default: message
	MessageTag string
	// StopOnError If true: An error occurs, it will cease to continue to verify
	StopOnError bool
	// SkipOnEmpty Skip check on field not exist or value is empty
	SkipOnEmpty bool
	// UpdateSource Whether to update source field value, useful for struct validate
	UpdateSource bool
	// CheckDefault Whether to validate the default value set by the user
	CheckDefault bool
	// CheckZero Whether validate the default zero value. (intX,uintX: 0, string: "")
	CheckZero bool
	// CheckSubOnParentMarked True: only collect sub-struct rule on current field has rule.
	CheckSubOnParentMarked bool
	// ValidatePrivateFields Whether to validate private fields or not, especially when inheriting other other structs.
	//
	//  type foo struct {
	//	  Field int `json:"field" validate:"required"`
	//  }
	//  type bar struct {
	//    foo // <-- validate this field
	//    Field2 int `json:"field2" validate:"required"`
	//  }
	//
	// default: false
	ValidatePrivateFields bool
}
```

**Usage**:

```go
// change global opts
validate.Config(func(opt *validate.GlobalOption) {
	opt.StopOnError = false
	opt.SkipOnEmpty = false
})
```

### Validating Private (Unexported fields)
By default, private fields are skipped. It is not uncommon to find code such as the following

```go
type foo struct {
	somefield int
}

type Bar struct {
	foo
	SomeOtherField string
}
```
In order to have `foo.somefield` validated, enable the behavior by setting `GlobalOption.ValidatePrivateFields` to `true`.

```go
validate.Config(func(opt *validate.GlobalOption) {
	opt.ValidatePrivateFields = true
})

```

### Custom Error Messages

- Register language messages

```go
import "github.com/gookit/validate/locales/zhcn"

// for all Validation.
// NOTICE: must be registered before on validate.New(), it only need call at once.
zhcn.RegisterGlobal()

// ... ...

v := validate.New()

// only for current Validation
zhcn.Register(v)
```

- Manual add global messages

```go
validate.AddGlobalMessages(map[string]string{
    "minLength": "OO! {field} min length is %d",
})
```

- Add messages for current validation

```go
v := validate.New(map[string]any{
    "name": "inhere",
})
v.StringRule("name", "required|string|minLen:7|maxLen:15")

v.AddMessages(map[string]string{
    "minLength": "OO! {field} min length is %d",
    "name.minLen": "OO! username min length is %d",
})
```

- Use struct tags: `message, label`

```go
type UserForm struct {
    Name  string `validate:"required|minLen:7" label:"User Name"`
    Email string `validate:"email" message:"email is invalid" label:"User Email"`
}
```

- Use struct method `Messages()`

```go
// Messages you can custom validator error messages. 
func (f UserForm) Messages() map[string]string {
	return validate.MS{
		"required": "oh! the {field} is required",
		"Name.required": "message for special field",
	}
}
```

### Add Custom Validator

`validate` supports adding custom validators, and supports adding `global validator` and `temporary validator`.

- **Global Validator** is globally valid and can be used everywhere
- **Temporary Validator** added to the current validation instance, only the current validation is available
- Add verification method to the structure. How to use please see the structure verification example above

> Note: The validator method must return a `bool` to indicate whether the validation was successful.
> The first parameter is the corresponding field value. If there are additional parameters, they will be appended automatically.

#### Add Global Validator

You can add one or more custom validators at once.

```go
validate.AddValidator("myCheck0", func(val any) bool {
	// do validate val ...
	return true
})
validate.AddValidators(validate.M{
	"myCheck1": func(val any) bool {
		// do validate val ...
		return true
	},
})
```

#### Add Temporary Validator

Again, you can add one or more custom validators at once.

```go
v := validate.Struct(u)
v.AddValidator("myFunc3", func(val any) bool {
	// do validate val ...
	return true
})
v.AddValidators(validate.M{
	"myFunc4": func(val any) bool {
		// do validate val ...
		return true
	},
})
```

### Add Custom Filter

`validate` can also support adding custom filters, and supports adding `global filter` and `temporary filter`.

- **Global Filter** is globally valid and can be used everywhere
- **Temporary Filter** added to the current validation instance, only the current validation is available

> TIP: for filter func, we allow functions with 1 result or 2 results where the second is an error.

#### Add Global Filter

You can add one or more custom validators at once.

```go
package main

import "github.com/gookit/validate"

func init() {
	validate.AddFilter("myToIntFilter0", func(val any) int {
		// do filtering val ...
		return 1
	})
	validate.AddFilters(validate.M{
		"myToIntFilter1": func(val any) (int, error) {
			// do filtering val ...
			return 1, nil
		},
	})
}
```

#### Add Temporary Filter

Again, you can add one or more custom filters at once.

```go
package main

import "github.com/gookit/validate"

func main() {
	v := validate.New(&someStrcut{})

	v.AddFilter("myToIntFilter0", func(val any) int {
		// do filtering val ...
		return 1
	})
	v.AddFilters(validate.M{
		"myToIntFilter1": func(val any) (int, error) {
			// do filtering val ...
			return 1, nil
		},
	})
	// use the added filter
	v.FilterRule("field", "myToIntFilter0")
}
```

### Custom `required` validation

Allows a custom `required` validator to customize whether the validation is empty.
However, note that the validator name must start with `required`, e.g. `required_custom`.

```go
	type Data struct {
		Age  int    `validate:"required_custom" message:"age is required"`
		Name string `validate:"required"`
	}

	v := validate.New(&Data{
		Name: "tom",
		Age:  0,
	})

	v.AddValidator("required_custom", func(val any) bool {
		// do check value
		return false
	})

	ok := v.Validate()
	assert.False(t, ok)
```

## Use on gin framework

Can use `validate` in any frameworks, such as Gin, Echo, Chi and more.

**Examples on gin:**

```go
package main
import (
    "github.com/gin-gonic/gin/binding"
    "github.com/gookit/validate"
)

// implements the binding.StructValidator
type customValidator struct {}

func (c *customValidator) ValidateStruct(ptr any) error {
    v := validate.Struct(ptr)
    v.Validate() // do validating
    
    if v.Errors.Empty() {
	return nil
    }

    return v.Errors
}

func (c *customValidator) Engine() any {
    return nil
}

func main()  {
	// ...

    // after init gin, set custom validator
    binding.Validator = &customValidator{}
}
```

<a id="built-in-validators"></a>
## Built In Validators

> Camel-style validator names now have underlined aliases. `endsWith` can also be written as `ends_with`

validator/aliases | description
-------------------|-------------------------------------------
`required`  | Check value is required and cannot be empty. 
`required_if/requiredIf`  | `required_if:anotherfield,value,...` The field under validation must be present and not empty if the `anotherField` field is equal to any value.
`requiredUnless`  | `required_unless:anotherfield,value,...` The field under validation must be present and not empty unless the `anotherField` field is equal to any value. 
`requiredWith`  | `required_with:foo,bar,...` The field under validation must be present and not empty only if any of the other specified fields are present.
`requiredWithAll`  | `required_with_all:foo,bar,...` The field under validation must be present and not empty only if all of the other specified fields are present.
`requiredWithout`  | `required_without:foo,bar,...` The field under validation must be present and not empty only when any of the other specified fields are not present.
`requiredWithoutAll`  | `required_without_all:foo,bar,...` The field under validation must be present and not empty only when all of the other specified fields are not present. 
`-/safe`  | The field values are safe and do not require validation
`int/integer/isInt`  | Check value is `intX` `uintX` type, And support size checking. eg: `"int"` `"int:2"` `"int:2,12"`
`uint/isUint`  |  Check value is uint(`uintX`) type, `value >= 0`
`bool/isBool`  |  Check value is bool string(`true`: "1", "on", "yes", "true", `false`: "0", "off", "no", "false").
`string/isString`  |  Check value is string type.
`float/isFloat`  |  Check value is float(`floatX`) type
`slice/isSlice`  |  Check value is slice type(`[]intX` `[]uintX` `[]byte` `[]string` ...).
`in/enum`  |  Check if the value is in the given enumeration `"in:a,b"`
`not_in/notIn`  |  Check if the value is not in the given enumeration `"contains:b"`
`contains`  |  Check if the input value contains the given value
`not_contains/notContains`  |  Check if the input value not contains the given value
`string_contains/stringContains`  |  Check if the input string value is contains the given sub-string
`starts_with/startsWith`  |  Check if the input string value is starts with the given sub-string
`ends_with/endsWith`  |  Check if the input string value is ends with the given sub-string
`range/between`  |  Check that the value is a number and is within the given range
`max/lte`  |  Check value is less than or equal to the given value
`min/gte`  |  Check value is greater than or equal to the given value(for `intX` `uintX` `floatX`)
`eq/equal/isEqual`  |  Check that the input value is equal to the given value
`ne/notEq/notEqual`  |  Check that the input value is not equal to the given value
`lt/lessThan`  |  Check value is less than the given value(use for `intX` `uintX` `floatX`)
`gt/greaterThan`  |  Check value is greater than the given value(use for `intX` `uintX` `floatX`)
`email/isEmail`  |   Check value is email address string.
`intEq/intEqual`  |  Check value is int and equals to the given value.
`len/length`  |  Check value length is equals to the given size(use for `string` `array` `slice` `map`).
`regex/regexp`  |  Check if the value can pass the regular verification
`arr/list/array/isArray`  |   Check value is array, slice type
`map/isMap`  |  Check value is a MAP type
`strings/isStrings`  |  Check value is string slice type(only allow `[]string`).
`ints/isInts`  |  Check value is int slice type(only allow `[]int`).
`min_len/minLen/minLength`  |  Check the minimum length of the value is the given size
`max_len/maxLen/maxLength`  |  Check the maximum length of the value is the given size
`eq_field/eqField`  |  Check that the field value is equals to the value of another field
`ne_field/neField`  |  Check that the field value is not equals to the value of another field
`gte_field/gteField`  |  Check that the field value is greater than or equal to the value of another field
`gt_field/gtField`  |  Check that the field value is greater than the value of another field
`lte_field/lteField`  |  Check if the field value is less than or equal to the value of another field
`lt_field/ltField`  |  Check that the field value is less than the value of another field
`file/isFile`  |  Verify if it is an uploaded file
`image/isImage`  |  Check if it is an uploaded image file and support suffix check
`mime/mimeType/inMimeTypes`  |  Check that it is an uploaded file and is in the specified MIME type
`date/isDate` | Check the field value is date string. eg `2018-10-25`
`gt_date/gtDate/afterDate` | Check that the input value is greater than the given date string.
`lt_date/ltDate/beforeDate` | Check that the input value is less than the given date string
`gte_date/gteDate/afterOrEqualDate` | Check that the input value is greater than or equal to the given date string.
`lte_date/lteDate/beforeOrEqualDate` | Check that the input value is less than or equal to the given date string.
`has_whitespace/hasWhitespace` | Check value string has Whitespace.
`ascii/ASCII/isASCII` | Check value is ASCII string.
`alpha/isAlpha` | Verify that the value contains only alphabetic characters
`alphaNum/isAlphaNum` | Check that only letters, numbers are included
`alphaDash/isAlphaDash` | Check to include only letters, numbers, dashes ( - ), and underscores ( _ )
`multiByte/isMultiByte` | Check value is MultiByte string.
`base64/isBase64` | Check value is Base64 string.
`dns_name/dnsName/DNSName/isDNSName` | Check value is DNSName string.
`data_uri/dataURI/isDataURI` | Check value is DataURI string.
`empty/isEmpty` | Check value is Empty string.
`hex_color/hexColor/isHexColor` | Check value is Hex color string.
`hexadecimal/isHexadecimal` | Check value is Hexadecimal string.
`json/JSON/isJSON` | Check value is JSON string.
`lat/latitude/isLatitude` | Check value is Latitude string.
`lon/longitude/isLongitude` | Check value is Longitude string.
`mac/isMAC` | Check value is MAC string.
`num/number/isNumber` | Check value is number string. `>= 0`
`cn_mobile/cnMobile/isCnMobile` | Check value is china mobile number string.
`printableASCII/isPrintableASCII` | Check value is PrintableASCII string.
`rgb_color/rgbColor/RGBColor/isRGBColor` | Check value is RGB color string.
`url/isURL` | Check value is URL string.
`fullUrl/isFullURL` | Check value is full URL string(_must start with http,https_).
`ip/isIP`  |  Check value is IP(v4 or v6) string.
`ipv4/isIPv4`  |  Check value is IPv4 string.
`ipv6/isIPv6`  |  Check value is IPv6 string.
`CIDR/isCIDR` | Check value is CIDR string.
`CIDRv4/isCIDRv4` | Check value is CIDRv4 string.
`CIDRv6/isCIDRv6` | Check value is CIDRv6 string.
`uuid/isUUID` | Check value is UUID string.
`uuid3/isUUID3` | Check value is UUID3 string.
`uuid4/isUUID4` | Check value is UUID4 string.
`uuid5/isUUID5` | Check value is UUID5 string.
`filePath/isFilePath` | Check value is an existing file path
`unixPath/isUnixPath` | Check value is Unix Path string.
`winPath/isWinPath` | Check value is Windows Path string.
`isbn10/ISBN10/isISBN10` | Check value is ISBN10 string.
`isbn13/ISBN13/isISBN13` | Check value is ISBN13 string.

**Notice:**

- `intX` is contains: int, int8, int16, int32, int64
- `uintX` is contains: uint, uint8, uint16, uint32, uint64
- `floatX` is contains: float32, float64

<a id="built-in-filters"></a>
## Built In Filters

> Filters powered by: [gookit/filter](https://github.com/gookit/filter)

filter/aliases | description 
-------------------|-------------------------------------------
`int/toInt`  | Convert value(string/intX/floatX) to `int` type `v.FilterRule("id", "int")`
`uint/toUint`  | Convert value(string/intX/floatX) to `uint` type `v.FilterRule("id", "uint")`
`int64/toInt64`  | Convert value(string/intX/floatX) to `int64` type `v.FilterRule("id", "int64")`
`float/toFloat`  | Convert value(string/intX/floatX) to `float` type
`bool/toBool`   | Convert string value to bool. (`true`: "1", "on", "yes", "true", `false`: "0", "off", "no", "false")
`trim/trimSpace`  | Clean up whitespace characters on both sides of the string
`ltrim/trimLeft`  | Clean up whitespace characters on left sides of the string
`rtrim/trimRight`  | Clean up whitespace characters on right sides of the string
`int/integer`  | Convert value(string/intX/floatX) to int type `v.FilterRule("id", "int")`
`lower/lowercase` | Convert string to lowercase
`upper/uppercase` | Convert string to uppercase
`lcFirst/lowerFirst` | Convert the first character of a string to lowercase
`ucFirst/upperFirst` | Convert the first character of a string to uppercase
`ucWord/upperWord` | Convert the first character of each word to uppercase
`camel/camelCase` | Convert string to camel naming style
`snake/snakeCase` | Convert string to snake naming style
`escapeJs/escapeJS` | Escape JS string.
`escapeHtml/escapeHTML` | Escape HTML string.
`str2ints/strToInts` | Convert string to int slice `[]int` 
`str2time/strToTime` | Convert date string to `time.Time`.
`str2arr/str2array/strToArray` | Convert string to string slice `[]string`

## Gookit packages

- [gookit/ini](https://github.com/gookit/ini) Go config management, use INI files
- [gookit/rux](https://github.com/gookit/rux) Simple and fast request router for golang HTTP 
- [gookit/gcli](https://github.com/gookit/gcli) build CLI application, tool library, running CLI commands
- [gookit/event](https://github.com/gookit/event) Lightweight event manager and dispatcher implements by Go
- [gookit/cache](https://github.com/gookit/cache) Generic cache use and cache manager for golang. support File, Memory, Redis, Memcached.
- [gookit/config](https://github.com/gookit/config) Go config management. support JSON, YAML, TOML, INI, HCL, ENV and Flags
- [gookit/color](https://github.com/gookit/color) A command-line color library with true color support, universal API methods and Windows support
- [gookit/filter](https://github.com/gookit/filter) Provide filtering, sanitizing, and conversion of golang data
- [gookit/validate](https://github.com/gookit/validate) Use for data validation and filtering. support Map, Struct, Form data
- [gookit/goutil](https://github.com/gookit/goutil) Some utils for the Go: string, array/slice, map, format, cli, env, filesystem, test and more
- More please see https://github.com/gookit

## See also

- https://github.com/albrow/forms
- https://github.com/asaskevich/govalidator
- https://github.com/go-playground/validator
- https://github.com/inhere/php-validate

## License

**[MIT](LICENSE)**
