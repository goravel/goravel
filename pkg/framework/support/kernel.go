package support

import "github.com/gin-gonic/gin"

type Kernel interface {
	Middleware() []gin.HandlerFunc
}