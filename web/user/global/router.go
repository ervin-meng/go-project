package global

import "github.com/gin-gonic/gin"

type ApiRouter struct {
	*gin.RouterGroup
}

var Router *ApiRouter
