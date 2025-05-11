package read

import "github.com/gin-gonic/gin"

type Reader interface {
	Read(c *gin.Context)
}
