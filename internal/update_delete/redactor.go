package redactor

import "github.com/gin-gonic/gin"

type Redactor interface {
	Update(c *gin.Context)
	Delete(c *gin.Context)
}
