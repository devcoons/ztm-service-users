package routes

import (
	models "api-users/models"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RouteGETMePermissions(c *gin.Context) {

	claims, srv, ok := InitServiceSJWT(c)

	if !ok {
		c.AbortWithStatus(401)
		return
	}

	perms := models.UsersPermissionsGetByUserId(srv.Database, claims.UserId)

	if perms == nil {
		c.AbortWithStatus(404)
		return
	}

	r, _ := json.Marshal(struct {
		Id        int       `json:"id"`
		UserId    int       `json:"user_id"`
		Perm01    int       `json:"p_01"`
		Perm02    int       `json:"p_02"`
		Perm03    int       `json:"p_03"`
		Perm04    int       `json:"p_04"`
		Perm05    int       `json:"p_05"`
		Perm06    int       `json:"p_06"`
		Perm07    int       `json:"p_07"`
		Perm08    int       `json:"p_08"`
		Perm09    int       `json:"p_09"`
		Perm10    int       `json:"p_10"`
		Perm11    int       `json:"p_11"`
		Perm12    int       `json:"p_12"`
		Perm13    int       `json:"p_13"`
		Perm14    int       `json:"p_14"`
		Perm15    int       `json:"p_15"`
		Perm16    int       `json:"p_16"`
		CreatedAt time.Time `json:"cr_at"`
		UpdatedAt time.Time `json:"up_at"`
	}{perms.Id, perms.UserId, perms.Perm01, perms.Perm02, perms.Perm03, perms.Perm04, perms.Perm05, perms.Perm06, perms.Perm07, perms.Perm08, perms.Perm09, perms.Perm10, perms.Perm11, perms.Perm12, perms.Perm13, perms.Perm14, perms.Perm15, perms.Perm16, perms.CreatedAt, perms.UpdatedAt})
	c.Data(http.StatusOK, gin.MIMEJSON, (r))
}
