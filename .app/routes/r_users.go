package routes

import (
	models "api-users/models"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RouteGETUsers(c *gin.Context) {

	claims, srv, ok := InitServiceSJWT(c)

	if !ok {
		c.AbortWithStatus(401)
		return
	}

	if !claims.Auth {
		c.AbortWithStatus(401)
		return
	}

	users := models.UsersGetAll(srv.Database)

	if users == nil {
		c.IndentedJSON(http.StatusOK, nil)
		return
	}

	var usersOverview []models.UserJsonOverview

	for _, b := range *users {
		nu := models.UserJsonOverview{Id: b.Id, Username: b.Username, Role: b.Role, Image: b.Image, FirstName: b.FirstName, LastName: b.LastName, LastLogin: b.LastLogin, CreatedAt: b.CreatedAt, UpdatedAt: b.UpdatedAt}
		usersOverview = append(usersOverview, nu)
	}

	c.IndentedJSON(http.StatusOK, usersOverview)
}

func RouteGETUsersComplete(c *gin.Context) {

	claims, srv, ok := InitServiceSJWT(c)

	if !ok {
		c.AbortWithStatus(401)
		return
	}

	if !claims.Auth {
		c.AbortWithStatus(401)
		return
	}

	users := models.UsersGetAll(srv.Database)

	c.IndentedJSON(http.StatusOK, users)
}

func RoutePUTUsersPasswordRecovery(c *gin.Context) {

	claims, srv, ok := InitServiceSJWT(c)

	if !ok {
		c.AbortWithStatus(401)
		return
	}

	if claims.Auth {
		c.AbortWithStatus(401)
		return
	}

	var values = UnmashalBody(c.Request.Body)

	if values["rec_token"] == nil || values["password"] == nil {
		c.AbortWithStatus(403)
		return
	}

	if len(values["password"].(string)) < 7 {
		c.AbortWithStatus(403)
		return
	}

	recToken := models.UsersRecoveryGetByTokenId(srv.Database, values["password"].(string))

	if recToken == nil {
		c.AbortWithStatus(403)
		return
	}

	if recToken.ExpireAt.After(time.Now()) {
		models.UsersRecoveryDeleteById(srv.Database, recToken.Id)
		c.AbortWithStatus(402)
		return
	}

	user := models.UsersGetById(srv.Database, recToken.Id)
	models.UsersRecoveryDeleteById(srv.Database, recToken.Id)

	if user == nil {
		c.AbortWithStatus(402)
		return
	}

	var updatevals map[string]interface{} = make(map[string]interface{})

	updatevals["password"] = values["password"].(string)

	if user.UpdateMapped(srv.Database, updatevals) {
		c.IndentedJSON(http.StatusOK, nil)
	} else {
		c.AbortWithStatus(402)
		return
	}
}

func RouteGETUsersPasswordRecovery(c *gin.Context) {

	claims, srv, ok := InitServiceSJWT(c)

	if !ok {
		c.AbortWithStatus(401)
		return
	}

	if claims.Auth {
		c.AbortWithStatus(401)
		return
	}

	uname := c.Param("uname")

	if uname == "" {
		c.AbortWithStatus(401)
		return
	}

	userId := models.UsersGetIdByUsername(srv.Database, uname)

	if userId == -1 {
		c.AbortWithStatus(401)
		return
	}

	recToken := models.UsersRecovery{UserId: userId}

	if !recToken.Create(srv.Database) {
		c.AbortWithStatus(401)
		return
	}

	r, _ := json.Marshal(struct {
		UserId        int       `json:"user_id"`
		RecoveryToken string    `json:"rec_token"`
		CreatedAt     time.Time `json:"cr_at"`
		ExpireAtint   time.Time `json:"ex_at"`
	}{recToken.UserId, recToken.RecoveryToken, recToken.CreatedAt, recToken.ExpireAt})

	c.IndentedJSON(http.StatusOK, r)
}
