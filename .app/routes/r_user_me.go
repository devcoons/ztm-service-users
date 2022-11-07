package routes

import (
	models "api-users/models"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RouteGETMeOverview(c *gin.Context) {

	claims, srv, ok := InitServiceSJWT(c)

	if !ok {
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
		Role      int    `json:"role"`
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		Image     []byte `json:"image"`
	}{user.Id, user.Username, user.Role, user.FirstName, user.LastName, user.Image})
	c.Data(http.StatusOK, gin.MIMEJSON, (r))
}
