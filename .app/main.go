package main

import (
	models "api-users/models"
	routes "api-users/routes"
	"fmt"
	"os"
	"runtime"

	ztm "github.com/devcoons/go-ztm"

	c "github.com/devcoons/go-fmt-colors"
	"github.com/gin-gonic/gin"
)

var APIService ztm.Service

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
	router.Use(ztm.AddUSEService(&APIService))
	router.GET("/", routes.RouteGET)

	////
	////

	router.POST("/auth", routes.RoutePOSTLogin)
	router.POST("/register", routes.RoutePOSTRegister)
	router.PATCH("/nonce", routes.RoutePATCHMeNonce)
	router.GET("/admin", routes.RouteGETAdmin)
	router.GET("/nonce", routes.RouteGETMeNonce)
	router.GET("/users/me", routes.RouteGETMeOverview)
	router.PUT("/users/me", routes.RoutePUTMe)
	router.GET("/users/me/complete", routes.RouteGETMeComplete)
	router.GET("/users/me/complete-perms", routes.RouteGETMeCompleteWPermissions)
	router.GET("/users/me/permissions", routes.RouteGETMePermissions)
	router.GET("/users", routes.RouteGETUsers)
	router.GET("/users/:id", routes.RouteGETUserById)
	router.GET("/users/:id/complete", routes.RouteGETUserByIdComplete)
	router.GET("/users/complete", routes.RouteGETUsersComplete)
	router.PUT("/users/recovery", routes.RoutePUTUsersPasswordRecovery)
	router.GET("/users/recovery/:uname", routes.RouteGETUsersPasswordRecovery)
	router.DELETE("/users/system/database/reset", routes.RouteDELDatabase)

	////
	////

	APIService.Start(router)
	fmt.Println("[GIN] Starting service at [0.0.0.0:8080]")
	router.Run("0.0.0.0:8080")

}
