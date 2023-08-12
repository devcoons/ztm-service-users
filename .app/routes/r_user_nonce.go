package routes

import (
	models "api-users/models"
	"encoding/json"
	"fmt"
	"net/http"

	cryptutils "github.com/devcoons/go-cryptutils"
	ztm "github.com/devcoons/go-ztm"
	"github.com/gin-gonic/gin"
)

func RouteGETMeNonce(c *gin.Context) {

	claims, srv, ok := ztm.InitServiceSJWT(c)

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
		Id    int    `json:"id"`
		Nonce string `json:"nonce"`
	}{user.Id, user.Nonce})
	c.Data(http.StatusOK, gin.MIMEJSON, (r))
}

func RoutePATCHMeNonce(c *gin.Context) {

	claims, srv, ok := ztm.InitServiceSJWT(c)

	if !ok {
		c.AbortWithStatus(401)
		return
	}

	user := models.UsersGetById(srv.Database, claims.UserId)

	fmt.Println(claims)

	if user == nil {
		c.AbortWithStatus(404)
		return
	}

	var nnonce = cryptutils.RandString(6)
	for nnonce == user.Nonce {
		nnonce = cryptutils.RandString(6)
	}

	user.UpdateMapped(srv.Database, map[string]interface{}{"nonce": nnonce})

	r, _ := json.Marshal(struct {
		Id    int    `json:"id"`
		Nonce string `json:"nonce"`
	}{user.Id, user.Nonce})
	c.Data(http.StatusOK, gin.MIMEJSON, (r))
}
