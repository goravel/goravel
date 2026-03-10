package console

type Stubs struct {
}

func (r Stubs) Rule() string {
	return `package DummyPackage

import (
	"context"

	"github.com/goravel/framework/contracts/validation"
)

type DummyRule struct {
}

// Signature The name of the rule.
func (receiver *DummyRule) Signature() string {
	return "DummySignature"
}

// Passes Determine if the validation rule passes.
func (receiver *DummyRule) Passes(ctx context.Context, data validation.Data, val any, options ...any) bool {
	return true
}

// Message Get the validation error message.
func (receiver *DummyRule) Message(ctx context.Context) string {
	return ""
}
`
}

func (r Stubs) Filter() string {
	return `package DummyPackage

import "context"

type DummyFilter struct {
}

// Signature The signature of the filter.
func (receiver *DummyFilter) Signature() string {
	return "DummySignature"
}

// Handle defines the filter function to apply.
//
// The Handle method should return a function that processes an input and
// returns a transformed value. The function can either return the
// transformed value alone or a tuple of the transformed value and an error.
// The input to the filter function is flexible: the first input is the value
// of the key on which the filter is applied, and the rest of the inputs are
// the arguments passed to the filter.
//
// Example usages:
//
// 1. Return only the transformed value:
//    func (val string) int {
//        // conversion logic
//        return 1
//    }
//
// 2. Return the transformed value and an error:
//    func (val int) (int, error) {
//        // conversion logic with error handling
//        return 1, nil
//    }
//
// 3. Take additional arguments:
//    func (val string, def ...string) string {
//        if val == "" && len(def) > 0 {
//            return def[0]
//        }
//        return val
//    }
//
func (receiver *DummyFilter) Handle(ctx context.Context) any {
	return nil
}
`
}
