package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SwaggerExplorerRedirect(c *gin.Context) {
	c.Redirect(http.StatusPermanentRedirect, "/explorer/")
}

func SwaggerAPI(c *gin.Context) {
	c.File("./public/swagger.yaml")
}
