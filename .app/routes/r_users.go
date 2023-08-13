package routes

import (
	models "api-users/models"
	"encoding/json"
	"net/http"
	"time"

	ztm "github.com/devcoons/go-ztm"
	"github.com/gin-gonic/gin"
)

func RouteGETUserById(c *gin.Context) {

	claims, srv, ok := ztm.InitServiceSJWT(c)

	if !ok || srv == nil || claims == nil {
		c.AbortWithStatus(401)
		return
	}

	if !claims.Auth || claims.UserId == -1 {
		c.AbortWithStatus(401)
		return
	}

	str_id := c.Param("id")

	id := ztm.ConvertToInt(str_id, -1)

	if id == -1 {
		c.AbortWithStatus(401)
		return
	}

	user := models.UsersGetById(srv.Database, id)

	if user == nil {
		c.IndentedJSON(http.StatusOK, nil)
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

func RouteGETUserByIdComplete(c *gin.Context) {

	claims, srv, ok := ztm.InitServiceSJWT(c)

	if !ok || srv == nil || claims == nil {
		c.AbortWithStatus(401)
		return
	}

	if !claims.Auth || claims.UserId == -1 {
		c.AbortWithStatus(401)
		return
	}

	str_id := c.Param("id")

	id := ztm.ConvertToInt(str_id, -1)

	if id == -1 {
		c.AbortWithStatus(401)
		return
	}

	user := models.UsersGetById(srv.Database, id)

	if user == nil {
		c.IndentedJSON(http.StatusOK, nil)
		return
	}

	r, _ := json.Marshal(struct {
		Id          int       `json:"id"`
		Username    string    `json:"username"`
		Role        int       `json:"role"`
		FirstName   string    `json:"firstname"`
		LastName    string    `json:"lastname"`
		Image       []byte    `json:"image"`
		Company     string    `json:"company"`
		Email       string    `json:"email"`
		MobilePhone string    `json:"cellphone"`
		LandLine    string    `json:"landline"`
		Country     string    `json:"country"`
		Province    string    `json:"province"`
		City        string    `json:"city"`
		Address     string    `json:"address"`
		IsEnabled   bool      `json:"is_enabled"`
		LastLogin   time.Time `json:"last_login"`
		CreatedAt   time.Time `json:"cr_at"`
		UpdatedAt   time.Time `json:"up_at"`
	}{user.Id, user.Username, user.Role, user.FirstName, user.LastName, user.Image, user.Company, user.Email, user.MobilePhone, user.LandLine, user.Country, user.Province, user.City, user.Address, user.IsEnabled, user.LastLogin, user.CreatedAt, user.UpdatedAt})
	c.Data(http.StatusOK, gin.MIMEJSON, (r))
}

func RouteGETUsers(c *gin.Context) {

	claims, srv, ok := ztm.InitServiceSJWT(c)

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

	claims, srv, ok := ztm.InitServiceSJWT(c)

	if !ok {
		c.AbortWithStatus(401)
		return
	}

	if !claims.Auth {
		c.AbortWithStatus(401)
		return
	}

	users := models.UsersGetAll(srv.Database)

	var usersAll []models.UserJson

	for _, b := range *users {
		nu := models.UserJson{Id: b.Id, Username: b.Username, Role: b.Role, Image: b.Image, FirstName: b.FirstName, LastName: b.LastName, Email: b.Email, Company: b.Company, Country: b.Country, Province: b.Province, LastLogin: b.LastLogin, CreatedAt: b.CreatedAt, UpdatedAt: b.UpdatedAt}
		usersAll = append(usersAll, nu)
	}

	c.IndentedJSON(http.StatusOK, usersAll)
}

func RoutePUTUsersPasswordRecovery(c *gin.Context) {

	claims, srv, ok := ztm.InitServiceSJWT(c)

	if !ok {
		c.AbortWithStatus(401)
		return
	}

	if claims.Auth {
		c.AbortWithStatus(401)
		return
	}

	var values = ztm.UnmashalBody(c.Request.Body)

	if values["rec_token"] == nil || values["password"] == nil {
		c.AbortWithStatus(403)
		return
	}

	if len(values["password"].(string)) < 7 {
		c.AbortWithStatus(403)
		return
	}

	recToken := models.UsersRecoveryGetByTokenId(srv.Database, values["rec_token"].(string))

	if recToken == nil {
		c.AbortWithStatus(403)
		return
	}

	if time.Now().After(recToken.ExpireAt) {
		models.UsersRecoveryDeleteById(srv.Database, recToken.Id)
		c.AbortWithStatus(402)
		return
	}

	user := models.UsersGetById(srv.Database, recToken.UserId)
	models.UsersRecoveryDeleteById(srv.Database, recToken.Id)

	if user == nil {
		c.AbortWithStatus(401)
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

	claims, srv, ok := ztm.InitServiceSJWT(c)

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

	user := models.UsersGetByUsername(srv.Database, uname)

	if user == nil {
		c.AbortWithStatus(401)
		return
	}
	models.UsersRecoveryDeleteByUserId(srv.Database, user.Id)
	recToken := models.UsersRecovery{UserId: user.Id}

	if !recToken.Create(srv.Database) {
		c.AbortWithStatus(401)
		return
	}

	r, _ := json.Marshal(struct {
		Email         string    `json:"email"`
		RecoveryToken string    `json:"rec_token"`
		CreatedAt     time.Time `json:"cr_at"`
		ExpireAtint   time.Time `json:"ex_at"`
	}{user.Email, recToken.RecoveryToken, recToken.CreatedAt, recToken.ExpireAt})

	c.Data(http.StatusOK, gin.MIMEJSON, (r))
}
