package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	aJWT "github.com/devcoons/go-auth-jwt"
	c "github.com/devcoons/go-fmt-colors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/logrusorgru/aurora"
	"github.com/mitchellh/mapstructure"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Service struct {
	SJwt     *aJWT.AuthJWT
	Config   *ServiceConfiguration
	Database *gorm.DB
	Rdb      *redis.Client
}

func (u *Service) Initialize(cfgpath string) bool {
	var err error

	u.Config = &ServiceConfiguration{}
	r := u.Config.Load(cfgpath)
	fmt.Println(c.FmtFgBgWhiteLBlue+"[ IMS ]"+c.FmtReset, c.FmtFgBgWhiteBlue+" INFO "+c.FmtReset, c.FmtFgBgWhiteBlack+"Loading service configuration"+c.FmtReset)

	if !r {
		return false
	}

	if u.Config.RedisDB.Host != "" {
		fmt.Println(c.FmtFgBgWhiteLBlue+"[ IMS ]"+c.FmtReset, c.FmtFgBgWhiteBlue+" INFO "+c.FmtReset, c.FmtFgBgWhiteBlack+"Redis Instance will be available"+c.FmtReset)

		u.Rdb = redis.NewClient(&redis.Options{
			Addr:     u.Config.RedisDB.Host + ":" + strconv.Itoa(u.Config.RedisDB.Port),
			Username: u.Config.RedisDB.Username,
			Password: u.Config.RedisDB.Password,
			DB:       u.Config.RedisDB.DB,
		})
	} else {
		fmt.Println(aurora.BgBrightYellow("[ IMS ] Redis Instance will NOT be available.."))
	}

	if u.Config.Secrets != nil {
		for _, s := range u.Config.Secrets {
			if strings.ToLower(s.Name) == "sjwt" {
				fmt.Println(c.FmtFgBgWhiteLBlue+"[ IMS ]"+c.FmtReset, c.FmtFgBgWhiteBlue+" INFO "+c.FmtReset, c.FmtFgBgWhiteBlack+"SJWT Token will be available"+c.FmtReset)
				u.SJwt = &aJWT.AuthJWT{}
				u.SJwt.SecretKey = s.Secret
				u.SJwt.TokenDuration = time.Duration(s.Duration) * time.Second
				u.SJwt.AuthType = s.AuthType
			}
		}
	} else {
		fmt.Println(aurora.BgRed("[ IMS ] Microservice cannot work without Secrets"))
		return false
	}

	if u.Config.Database.Host != "" {
		u.Database = &gorm.DB{}

		dsn := u.Config.Database.Username + ":" + u.Config.Database.Password + "@tcp("
		dsn += u.Config.Database.Host + ":" + strconv.Itoa(u.Config.Database.Port) + ")/"
		dsn += u.Config.Database.DbName + "?parseTime=true"

		for i := 1; i <= 5; i++ {
			fmt.Println(c.FmtFgBgWhiteLBlue+"[ IMS ]"+c.FmtReset, c.FmtFgBgWhiteBlue+" INFO "+c.FmtReset, c.FmtFgBgWhiteBlack+"Connecting SQL database: "+u.Config.Database.Host+":"+strconv.Itoa(u.Config.Database.Port)+c.FmtReset)

			u.Database, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
			if err != nil {
				fmt.Println(aurora.BgBrightYellow("[ IMS ] Connection failed. Retring in 7 seconds.."))
				time.Sleep(7 * time.Second)
			} else {
				fmt.Println(c.FmtFgBgWhiteLBlue+"[ IMS ]"+c.FmtReset, c.FmtFgBgWhiteBlue+" INFO "+c.FmtReset, c.FmtFgBgWhiteBlack+"Connection succesfully completed"+c.FmtReset)
				break
			}
		}
	} else {
		fmt.Println(aurora.BgBrightYellow("[ IMS ] Sql Database will NOT be available.."))
	}

	return err == nil
}

func AddUSEService(u *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("service", u)
		c.Next()
	}
}

func (u *Service) ValidateServiceJWT(r *http.Request) *SJWTClaims {

	if u.Rdb == nil {
		return nil
	}

	iclaims, ok := u.SJwt.IsAuthorized(r)

	if !ok {
		return nil
	}

	var claims SJWTClaims

	err := mapstructure.Decode(iclaims, &claims)

	if err != nil {
		return nil
	}

	return &claims
}

func (u *Service) UpdateUserNonce(userId int, userNonce string) bool {

	var ctx = context.Background()
	if u.Rdb == nil {
		return false
	}

	_, err := u.Rdb.Get(ctx, strconv.Itoa(userId)).Result()
	if err == nil {
		_, _ = u.Rdb.Del(ctx, strconv.Itoa(userId)).Result()
	}
	u.Rdb.Set(ctx, strconv.Itoa(userId), userNonce, 0)

	for _, gw := range u.Config.Gateways {

		var sclaims SJWTClaims
		sclaims.Auth = true
		sclaims.UserId = userId
		sclaims.Role = 0
		sclaims.Service = "api-users"

		token := u.SJwt.GenerateJWT(sclaims)

		client := &http.Client{}
		req, _ := http.NewRequest("GET", gw.Host+":"+strconv.Itoa(gw.Port)+"/syncunc", nil)
		req.Header.Del("Authorization")
		req.Header.Add("Authorization", "X-Fowarder "+token)
		req.Body = nil
		client.Do(req)
	}
	return true
}

func (u *Service) ReloadUserNonceFromDB(userId int, userNonce string) bool {

	var ctx = context.Background()
	if u.Rdb == nil {
		return false
	}

	_, err := u.Rdb.Get(ctx, strconv.Itoa(userId)).Result()
	if err == nil {
		_, _ = u.Rdb.Del(ctx, strconv.Itoa(userId)).Result()
	}
	u.Rdb.Set(ctx, strconv.Itoa(userId), userNonce, 0)

	return true
}
