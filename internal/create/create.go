package create

import "github.com/gin-gonic/gin"

type Creator interface {
	Create(c *gin.Context)
}
