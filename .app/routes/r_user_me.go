package routes

import (
	models "api-users/models"
	"encoding/json"
	"net/http"
	"strings"
	"time"

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

func RoutePUTMe(c *gin.Context) {

	claims, srv, ok := InitServiceSJWT(c)

	if !ok || claims.UserId == -1 {
		c.AbortWithStatus(401)
		return
	}

	var values = UnmashalBody(c.Request.Body)
	var updatevals map[string]interface{} = make(map[string]interface{})

	if values["username"] != nil {
		v := values["username"].(string)
		if len(v) > 7 && strings.Contains(v, "@") {
			updatevals["username"] = v
		}
	}

	if values["fname"] != nil {
		updatevals["first_name"] = values["fname"].(string)
	}

	if values["lname"] != nil {
		updatevals["last_name"] = values["lname"].(string)
	}

	if values["img"] != nil {
		updatevals["image"] = values["img"]
	}

	if values["company"] != nil {
		updatevals["company"] = values["company"].(string)
	}

	if values["email"] != nil {
		updatevals["email"] = values["email"].(string)
	}

	if values["cellphone"] != nil {
		updatevals["mobile_phone"] = values["cellphone"].(string)
	}

	if values["landline"] != nil {
		updatevals["land_line"] = values["landline"].(string)
	}

	if values["country"] != nil {
		updatevals["country"] = values["country"].(string)
	}

	if values["province"] != nil {
		updatevals["province"] = values["province"].(string)
	}
	if values["city"] != nil {
		updatevals["city"] = values["city"].(string)
	}

	if values["address"] != nil {
		updatevals["address"] = values["address"].(string)
	}
	user := models.UsersGetById(srv.Database, claims.UserId)

	if user == nil {
		c.AbortWithStatus(404)
		return
	}

	user.UpdateMapped(srv.Database, updatevals)

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

func RouteGETMeComplete(c *gin.Context) {

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

func RouteGETMeCompleteWPermissions(c *gin.Context) {

	claims, srv, ok := InitServiceSJWT(c)

	if !ok {
		c.AbortWithStatus(401)
		return
	}

	user := models.UsersGetById(srv.Database, claims.UserId)
	perms := models.UsersPermissionsGetByUserId(srv.Database, claims.UserId)

	if user == nil {
		c.AbortWithStatus(404)
		return
	}

	if perms == nil {
		var perms = models.UsersPermissions{UserId: user.Id}
		if !perms.Create(srv.Database) {
			c.IndentedJSON(http.StatusExpectationFailed, nil)
			return
		}
	}

	r, _ := json.Marshal(struct {
		Id          int                     `json:"id"`
		Username    string                  `json:"username"`
		Role        int                     `json:"role"`
		Permissions models.UsersPermissions `json:"permissions"`
		FirstName   string                  `json:"firstname"`
		LastName    string                  `json:"lastname"`
		Image       []byte                  `json:"image"`
		Company     string                  `json:"company"`
		Email       string                  `json:"email"`
		MobilePhone string                  `json:"cellphone"`
		LandLine    string                  `json:"landline"`
		Country     string                  `json:"country"`
		Province    string                  `json:"province"`
		City        string                  `json:"city"`
		Address     string                  `json:"address"`
		IsEnabled   bool                    `json:"is_enabled"`
		LastLogin   time.Time               `json:"last_login"`
		CreatedAt   time.Time               `json:"cr_at"`
		UpdatedAt   time.Time               `json:"up_at"`
	}{user.Id, user.Username, user.Role, *perms, user.FirstName, user.LastName, user.Image, user.Company, user.Email, user.MobilePhone, user.LandLine, user.Country, user.Province, user.City, user.Address, user.IsEnabled, user.LastLogin, user.CreatedAt, user.UpdatedAt})
	c.Data(http.StatusOK, gin.MIMEJSON, (r))
}
