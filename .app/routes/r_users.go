package routes

import (
	models "api-users/models"
	"net/http"

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

	_, _, ok := InitServiceSJWT(c)

	if !ok {
		c.AbortWithStatus(401)
		return
	}

	c.IndentedJSON(http.StatusOK, nil)
}

func RouteGETUsersPasswordRecovery(c *gin.Context) {

	_, _, ok := InitServiceSJWT(c)

	if !ok {
		c.AbortWithStatus(401)
		return
	}

	c.IndentedJSON(http.StatusOK, nil)
}
