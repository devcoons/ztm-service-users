package routes

import (
	models "api-users/models"
	"encoding/json"
	"net/http"

	ztm "github.com/devcoons/go-ztm"
	"github.com/gin-gonic/gin"
)

func RouteGETAdmin(c *gin.Context) {

	claims, srv, ok := ztm.InitServiceSJWT(c)

	if !ok || claims == nil || !claims.Auth || claims.UserId == -1 {
		c.AbortWithStatus(401)
		return
	}

	user := models.UsersGetById(srv.Database, claims.UserId)

	if user == nil {
		c.AbortWithStatus(404)
		return
	}

	r, _ := json.Marshal(struct {
		Id        int    `json:"id"`
		Username  string `json:"username"`
		IsEnabled bool   `json:"enabled"`
		IsAdmin   bool   `json:"admin"`
	}{user.Id, user.Username, user.IsEnabled, user.IsAdmin})
	c.Data(http.StatusOK, gin.MIMEJSON, (r))
}
