package routes

import (
	middleware "api-users/middleware"
	models "api-users/models"
	"encoding/json"
	"net/http"

	cryptutils "github.com/devcoons/go-cryptutils"
	"github.com/gin-gonic/gin"
)

func RoutePOSTRegister(c *gin.Context) {

	srv, ok := c.MustGet("service").(*middleware.Service)

	if !ok || srv.Database == nil {
		c.IndentedJSON(http.StatusExpectationFailed, nil)
		return
	}

	values := UnmashalBody(c.Request.Body)

	if values == nil {
		c.AbortWithStatus(409)
		return
	}

	if values["security"] != "XYZABCFAKE" || len(values["username"].(string)) < 5 || len(values["password"].(string)) < 5 {
		c.IndentedJSON(http.StatusNotAcceptable, nil)
		return
	}

	var user = models.User{Username: values["username"].(string), Password: values["password"].(string), Role: 1, Nonce: cryptutils.RandString(6)}

	if !user.Create(srv.Database) {
		c.IndentedJSON(http.StatusExpectationFailed, nil)
		return
	}

	var perms = models.UsersPermissions{UserId: user.Id}
	if !perms.Create(srv.Database) {
		c.IndentedJSON(http.StatusExpectationFailed, nil)
		return
	}

	r, _ := json.Marshal(struct {
		Id       int    `json:"id"`
		Username string `json:"username"`
		Role     int    `json:"role"`
	}{user.Id, user.Username, user.Role})
	c.Data(http.StatusOK, gin.MIMEJSON, (r))
}
