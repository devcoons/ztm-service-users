package main

import (
	middleware "api-users/middleware"
	models "api-users/models"
	routes "api-users/routes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	c "github.com/devcoons/go-fmt-colors"
	"github.com/gin-gonic/gin"
)

var APIService middleware.Service

func main() {
	runtime.GOMAXPROCS(4)
	fmt.Println(c.FmtFgBgWhiteLBlue+"[ IMS ]"+c.FmtReset, c.FmtFgBgWhiteBlue+" INFO "+c.FmtReset, c.FmtFgBgWhiteBlack+"Initializing microservice."+c.FmtReset)

	cfgfile, present := os.LookupEnv("IMSCFGFILE")

	if !present {
		fmt.Println(c.FmtFgBgWhiteLBlue+"[ IMS ]"+c.FmtReset, c.FmtFgBgWhiteRed+" ERRN "+c.FmtReset, c.FmtFgBgWhiteBlack+"Configuration file env.variable `IMSCFGFILE` does not exist"+c.FmtReset)
		return
	}

	if !APIService.Initialize(cfgfile) {
		fmt.Println(c.FmtFgBgWhiteLBlue+"[ IMS ]"+c.FmtReset, c.FmtFgBgWhiteRed+" ERRN "+c.FmtReset, c.FmtFgBgWhiteBlack+"Initialization failed. Exiting application.."+c.FmtReset)
		return
	}

	fmt.Println(c.FmtFgBgWhiteLBlue+"[ IMS ]"+c.FmtReset, c.FmtFgBgWhiteBlue+" INFO "+c.FmtReset, c.FmtFgBgWhiteBlack+"Models Database auto-migration"+c.FmtReset)
	models.AutoMigrate(APIService.Database)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(middleware.AddUSEService(&APIService))

	router.GET("/", func(c *gin.Context) {

		claims := APIService.ValidateServiceJWT(c.Request)
		if claims == nil {
			c.Data(503, "application/json", nil)
			return
		} else {
			claims.Hop = claims.Hop - 1
		}
		c.IndentedJSON(200, nil)
	})

	router.POST("/auth", routes.RoutePOSTLogin)
	router.POST("/register", routes.RoutePOSTRegister)
	router.PATCH("/nonce", routes.RoutePATCHMeNonce)
	router.GET("/nonce", routes.RouteGETMeNonce)
	router.GET("/users/me", routes.RouteGETMeOverview)
	router.PUT("/users/me", routes.RoutePUTMe)
	router.GET("/users/me/complete", routes.RouteGETMeComplete)
	router.GET("/users/me/complete-perms", routes.RouteGETMeCompleteWPermissions)
	router.GET("/users/me/permissions", routes.RouteGETMePermissions)
	router.GET("/users", routes.RouteGETUsers)
	router.GET("/users/complete", routes.RouteGETUsersComplete)
	router.PUT("users/recovery", routes.RoutePUTUsersPasswordRecovery)
	router.GET("users/recovery/:uname", routes.RouteGETUsersPasswordRecovery)
	router.NoRoute(RequestForwarder)
	fmt.Println("[GIN] Starting service at [0.0.0.0:8080]")
	router.Run("0.0.0.0:8080")

}

func RequestForwarder(c *gin.Context) {

	var requestedPath = strings.TrimRight(c.Request.URL.Path, "/")
	var requestedUrlQuery = c.Request.URL.RawQuery

	for _, nodeDetails := range APIService.Config.Services {

		m, _ := regexp.MatchString(nodeDetails.URL, requestedPath)
		if m {
			claims := APIService.ValidateServiceJWT(c.Request)
			if claims == nil {
				c.Data(503, "application/json", nil)
				return
			} else {
				claims.Hop = claims.Hop - 1
			}

			if claims.Hop == 0 {
				c.Data(503, "application/json", nil)
				return
			}

			token := APIService.SJwt.GenerateJWT(claims)
			client := &http.Client{}
			req, _ := http.NewRequest(c.Request.Method, nodeDetails.Host+":"+strconv.Itoa(nodeDetails.Port)+requestedPath+"?"+requestedUrlQuery, c.Request.Body)
			req.Header = c.Request.Header
			req.Header.Del("Authorization")
			req.Header.Add("Authorization", APIService.SJwt.AuthType+" "+token)
			res, errn := client.Do(req)
			if errn == nil {
				body, _ := ioutil.ReadAll(res.Body)
				c.Data(res.StatusCode, res.Header.Get("Content-Type"), body)
			} else {
				c.Data(503, "application/json", nil)
			}
			return
		}
	}
}
