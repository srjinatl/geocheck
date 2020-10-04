package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthGET(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "UP",
	})
}
