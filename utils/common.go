package utils

import "github.com/gin-gonic/gin"

func ErrorResponse(ctx *gin.Context, err error, status int) {
	ctx.JSON(status, gin.H{"error": err.Error()})
}
