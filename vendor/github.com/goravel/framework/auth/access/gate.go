package access

import (
	"context"
	"fmt"

	"github.com/goravel/framework/contracts/auth/access"
)

type Gate struct {
	ctx             context.Context
	abilities       map[string]func(ctx context.Context, arguments map[string]any) access.Response
	beforeCallbacks []func(ctx context.Context, ability string, arguments map[string]any) access.Response
	afterCallbacks  []func(ctx context.Context, ability string, arguments map[string]any, result access.Response) access.Response
}

func NewGate(ctx context.Context) *Gate {
	return &Gate{
		ctx:       ctx,
		abilities: make(map[string]func(ctx context.Context, arguments map[string]any) access.Response),
	}
}

func (r *Gate) WithContext(ctx context.Context) access.Gate {
	return &Gate{
		ctx:             ctx,
		abilities:       r.abilities,
		beforeCallbacks: r.beforeCallbacks,
		afterCallbacks:  r.afterCallbacks,
	}
}

func (r *Gate) Allows(ability string, arguments map[string]any) bool {
	return r.Inspect(ability, arguments).Allowed()
}

func (r *Gate) Denies(ability string, arguments map[string]any) bool {
	return !r.Allows(ability, arguments)
}

func (r *Gate) Inspect(ability string, arguments map[string]any) access.Response {
	result := r.callBeforeCallbacks(r.ctx, ability, arguments)
	if result == nil {
		if _, exist := r.abilities[ability]; exist {
			result = r.abilities[ability](r.ctx, arguments)
		} else {
			result = NewDenyResponse(fmt.Sprintf("ability doesn't exist: %s", ability))
		}
	}

	return r.callAfterCallbacks(r.ctx, ability, arguments, result)
}

func (r *Gate) Define(ability string, callback func(ctx context.Context, arguments map[string]any) access.Response) {
	r.abilities[ability] = callback
}

func (r *Gate) Any(abilities []string, arguments map[string]any) bool {
	var res bool
	for _, ability := range abilities {
		res = res || r.Allows(ability, arguments)
	}

	return res
}

func (r *Gate) None(abilities []string, arguments map[string]any) bool {
	return !r.Any(abilities, arguments)
}

func (r *Gate) Before(callback func(ctx context.Context, ability string, arguments map[string]any) access.Response) {
	r.beforeCallbacks = append(r.beforeCallbacks, callback)
}

func (r *Gate) After(callback func(ctx context.Context, ability string, arguments map[string]any, result access.Response) access.Response) {
	r.afterCallbacks = append(r.afterCallbacks, callback)
}

func (r *Gate) callBeforeCallbacks(ctx context.Context, ability string, arguments map[string]any) access.Response {
	for _, before := range r.beforeCallbacks {
		result := before(ctx, ability, arguments)
		if result != nil {
			return result
		}
	}

	return nil
}

func (r *Gate) callAfterCallbacks(ctx context.Context, ability string, arguments map[string]any, result access.Response) access.Response {
	for _, after := range r.afterCallbacks {
		afterResult := after(ctx, ability, arguments, result)
		if result == nil && afterResult != nil {
			result = afterResult
		}
	}

	return result
}
