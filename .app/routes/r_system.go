package routes

import (
	models "api-users/models"
	"net/http"

	ztm "github.com/devcoons/go-ztm"

	"github.com/gin-gonic/gin"
)

func RouteDELDatabase(c *gin.Context) {

	srv, ok := c.MustGet("service").(*ztm.Service)

	if !ok || srv == nil {
		c.IndentedJSON(http.StatusInternalServerError, ztm.ErrorMsg{ErrorCode: "US-S-1000", Message: "Internal Issue. m.srv-[Users] could not properly init the internal services."})
		return
	}

	models.ResetMigration(srv.Database)

	c.IndentedJSON(http.StatusOK, nil)
}
