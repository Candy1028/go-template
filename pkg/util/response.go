package util

import (
	"github.com/Candy1028/go-template/pkg/comment/response"

	"github.com/gin-gonic/gin"
)

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"code":    response.OK,
		"message": response.GetMessage(response.OK),
		"data":    data,
	})
}

func ErrorResponse(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(200, gin.H{
		"code":    code,
		"message": msg,
		"data":    data,
	})
}
