package routes

import (
	models "api-users/models"
	"encoding/json"
	"net/http"
	"net/mail"

	ztm "github.com/devcoons/go-ztm"

	cryptutils "github.com/devcoons/go-cryptutils"
	"github.com/gin-gonic/gin"
)

func RoutePOSTRegister(c *gin.Context) {

	srv, ok := c.MustGet("service").(*ztm.Service)

	if !ok || srv.Database == nil {
		c.IndentedJSON(http.StatusExpectationFailed, nil)
		return
	}

	values := ztm.UnmashalBody(c.Request.Body)

	if values == nil {
		c.IndentedJSON(http.StatusNotAcceptable, ztm.ErrorMsg{ErrorCode: "US-R-0000", Message: "Required field are missing"})
		return
	}

	if values["security"] != "XYZABCFAKE" || len(values["username"].(string)) < 5 || len(values["password"].(string)) < 5 {
		c.IndentedJSON(http.StatusNotAcceptable, ztm.ErrorMsg{ErrorCode: "US-R-0001", Message: "Security Code, Username or Password are invalid"})
		return
	}

	if values["email"] == nil {
		c.IndentedJSON(http.StatusNotAcceptable, ztm.ErrorMsg{ErrorCode: "US-R-0002", Message: "E-mail address is required"})
		return
	}

	_, ern := mail.ParseAddress(values["email"].(string))

	if ern != nil {
		c.IndentedJSON(http.StatusNotAcceptable, ztm.ErrorMsg{ErrorCode: "US-R-0003", Message: "E-mail address is invalid"})
		return
	}

	var user = models.User{Username: values["username"].(string), Password: values["password"].(string), Role: 1, Nonce: cryptutils.RandString(6), Email: values["email"].(string)}

	if !user.Create(srv.Database) {
		c.IndentedJSON(http.StatusNotAcceptable, ztm.ErrorMsg{ErrorCode: "US-R-0004", Message: "User already exists"})
		return
	}

	var perms = models.UsersPermissions{UserId: user.Id}
	if !perms.Create(srv.Database) {
		c.IndentedJSON(http.StatusNotAcceptable, ztm.ErrorMsg{ErrorCode: "US-R-0005", Message: "User Permissions already exists"})
		return
	}

	r, _ := json.Marshal(struct {
		Id       int    `json:"id"`
		Username string `json:"username"`
		Role     int    `json:"role"`
	}{user.Id, user.Username, user.Role})
	c.Data(http.StatusOK, gin.MIMEJSON, (r))
}
