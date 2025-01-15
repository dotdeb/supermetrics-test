package middlewares

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dotdeb/supermetrics-test/internals/configs"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func TestSuccess(t *testing.T) {
	rec, route, envs, req := initServer()
	req.Header.Set("Authorization", createToken(envs.Jwt.ApiSecret,
		"HS512",
		envs.Jwt.Issuer,
		envs.Jwt.Audience,
		time.Now().Add(10*time.Second).Unix(),
		time.Now().Unix(),
	))
	route.ServeHTTP(rec, req)

	defer rec.Result().Body.Close()
	body, _ := io.ReadAll(rec.Result().Body)
	if string(body) != "reader writer" || rec.Code != http.StatusOK {
		panic("Auth middleware doesn't return correct response")
	}
}

func TestWrongSecret(t *testing.T) {
	rec, route, envs, req := initServer()
	req.Header.Set("Authorization", createToken("wrong-secret",
		"HS512",
		envs.Jwt.Issuer,
		envs.Jwt.Audience,
		time.Now().Add(-10*time.Second).Unix(),
		time.Now().Unix(),
	))
	route.ServeHTTP(rec, req)
	if rec.Code != 401 {
		panic("Wrong secret didn't return 401")
	}
}

func TestExpired(t *testing.T) {
	rec, route, envs, req := initServer()
	req.Header.Set("Authorization", createToken(envs.Jwt.ApiSecret,
		"HS512",
		envs.Jwt.Issuer,
		envs.Jwt.Audience,
		time.Now().Add(-10*time.Second).Unix(),
		time.Now().Unix(),
	))
	route.ServeHTTP(rec, req)
	if rec.Code != 401 {
		panic("Expired token didn't return 401")
	}
}

func TestWrongNbf(t *testing.T) {
	rec, route, envs, req := initServer()
	req.Header.Set("Authorization", createToken(envs.Jwt.ApiSecret,
		"HS512",
		envs.Jwt.Issuer,
		envs.Jwt.Audience,
		time.Now().Add(10*time.Second).Unix(),
		time.Now().Add(10*time.Second).Unix(),
	))
	route.ServeHTTP(rec, req)
	if rec.Code != 401 {
		panic(fmt.Sprintf(
			"Wrong issuer method didn't return 403: %v",
			rec.Code,
		))
	}
}

func TestWrongIss(t *testing.T) {
	rec, route, envs, req := initServer()
	req.Header.Set("Authorization", createToken(envs.Jwt.ApiSecret,
		"HS512",
		"wrong-issuer",
		envs.Jwt.Audience,
		time.Now().Add(10*time.Second).Unix(),
		time.Now().Unix(),
	))
	route.ServeHTTP(rec, req)
	fmt.Println(rec.Code)
	if rec.Code != 401 {
		panic(fmt.Sprintf(
			"Wrong issuer method didn't return 403: %v",
			rec.Code,
		))
	}
}

func TestWrongAud(t *testing.T) {
	rec, route, envs, req := initServer()
	req.Header.Set("Authorization", createToken(envs.Jwt.ApiSecret,
		"HS512",
		envs.Jwt.Issuer,
		"wront-audience",
		time.Now().Add(10*time.Second).Unix(),
		time.Now().Unix(),
	))
	route.ServeHTTP(rec, req)
	if rec.Code != 401 {
		panic(fmt.Sprintf(
			"Wrong issuer method didn't return 403: %v",
			rec.Code,
		))
	}
}

func initServer() (*httptest.ResponseRecorder, *gin.Engine, configs.Env, *http.Request) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	e := configs.Env{
		Jwt: configs.Jwt{
			ApiSecret: "unittest-secret",
			Issuer:    "unittest-issuer",
			Audience:  "unittest-audience",
		},
		Variables: configs.Variables{
			Port:    "",
			Timeout: 0,
		},
	}

	r.Use(BearerAuthMiddleware(e.Jwt))
	r.GET("/unittest", func(c *gin.Context) {
		y := c.MustGet("user_roles").([]string)
		c.String(http.StatusOK, strings.Join(y, " "))
	})

	req, _ := http.NewRequest(http.MethodGet, "/unittest", nil)
	return w, r, e, req
}

func createToken(secret string, alg string, iss string, aud string, exp int64, nbf int64) string {
	key := []byte(secret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS512,
		jwt.MapClaims{
			"alg":   alg,
			"iss":   iss,
			"aud":   aud,
			"roles": []string{"reader", "writer"},
			"exp":   exp,
			"nbf":   nbf,
		})
	signedJwt, _ := token.SignedString(key)
	return "Bearer " + signedJwt
}
