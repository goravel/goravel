package access

import "github.com/goravel/framework/contracts/auth/access"

func NewAllowResponse() access.Response {
	return &ResponseImpl{allowed: true}
}

func NewDenyResponse(message string) access.Response {
	return &ResponseImpl{message: message}
}

type ResponseImpl struct {
	message string
	allowed bool
}

func (r *ResponseImpl) Allowed() bool {
	return r.allowed
}

func (r *ResponseImpl) Message() string {
	return r.message
}
