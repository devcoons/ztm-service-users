package routes

import (
	models "api-users/models"
	"encoding/json"
	"net/http"

	cryptutils "github.com/devcoons/go-cryptutils"
	"github.com/gin-gonic/gin"
)

func RoutePOSTLogin(c *gin.Context) {

	_, srv, ok := InitServiceSJWT(c)

	if !ok || srv.SJwt == nil {
		c.IndentedJSON(http.StatusExpectationFailed, nil)
		return
	}

	values := UnmashalBody(c.Request.Body)

	if values == nil {
		c.AbortWithStatus(401)
		return
	}

	if len(values["username"].(string)) < 5 || len(values["password"].(string)) < 5 {
		c.IndentedJSON(http.StatusNotAcceptable, nil)
		return
	}

	user := models.UsersGetOneByUsernamePassword(srv.Database, values["username"].(string), values["password"].(string))
	if user == nil {
		r, _ := json.Marshal(struct {
			Code    string `json:"code"`
			Details string `json:"details"`
		}{"api-users:login:nil", "Authentication failed. Please check your credentials."})
		c.Data(http.StatusNotAcceptable, gin.MIMEJSON, (r))
		return
	}

	if !user.IsEnabled {
		r, _ := json.Marshal(struct {
			Code    string `json:"code"`
			Details string `json:"details"`
		}{"api-users:login:disabled", "Authentication failed. Your account is disabled."})
		c.Data(http.StatusNotAcceptable, gin.MIMEJSON, (r))
		return
	}
	var nnonce = cryptutils.RandString(6)

	for nnonce == user.Nonce {
		nnonce = cryptutils.RandString(6)
	}
	user.UpdateMapped(srv.Database, map[string]interface{}{"nonce": nnonce})

	r, _ := json.Marshal(struct {
		Id       int    `json:"id"`
		Username string `json:"username"`
		Role     int    `json:"role"`
		Nonce    string `json:"nonce"`
	}{user.Id, user.Username, user.Role, nnonce})
	c.Data(http.StatusOK, gin.MIMEJSON, (r))
}
