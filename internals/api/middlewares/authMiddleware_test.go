package middlewares

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	crypto "crypto/rand"
	"crypto/rsa"

	tokenparser "github.com/dotdeb/supermetrics-test/internals/api/middlewares/tokenParser"
	"github.com/dotdeb/supermetrics-test/internals/configs"
	"github.com/dotdeb/supermetrics-test/internals/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func TestSuccessHmac(t *testing.T) {
	rec, route, envs, req := initHmacServer()
	req.Header.Set("Authorization", createHmacToken(envs.Jwt.ApiSecret,
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

func TestSuccessRsa(t *testing.T) {
	// Rsa must overwrite env first as private/public key must exists during server init
	e := configs.Env{
		Jwt: configs.Jwt{
			Issuer:    "unittest-issuer",
			Audience:  "unittest-audience",
			RsaPublic: &rsa.PublicKey{},
		},
		Variables: configs.Variables{
			Port:    "",
			Timeout: 0,
		},
	}
	publicKey, token := createRsaToken(e.Jwt)
	e.Jwt.RsaPublic = publicKey

	rec, route, req := initRsaServer(e.Jwt)

	req.Header.Set("Authorization", "Bearer "+token)
	route.ServeHTTP(rec, req)

	defer rec.Result().Body.Close()
	body, _ := io.ReadAll(rec.Result().Body)

	if string(body) != "reader" || rec.Code != http.StatusOK {
		panic("Auth middleware doesn't return correct response")
	}
}

func TestWrongSecret(t *testing.T) {
	rec, route, envs, req := initHmacServer()
	req.Header.Set("Authorization", createHmacToken("wrong-secret",
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
	rec, route, envs, req := initHmacServer()
	req.Header.Set("Authorization", createHmacToken(envs.Jwt.ApiSecret,
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
	rec, route, envs, req := initHmacServer()
	req.Header.Set("Authorization", createHmacToken(envs.Jwt.ApiSecret,
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
	rec, route, envs, req := initHmacServer()
	req.Header.Set("Authorization", createHmacToken(envs.Jwt.ApiSecret,
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
	rec, route, envs, req := initHmacServer()
	req.Header.Set("Authorization", createHmacToken(envs.Jwt.ApiSecret,
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

func initHmacServer() (*httptest.ResponseRecorder, *gin.Engine, configs.Env, *http.Request) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	e := configs.Env{
		Jwt: configs.Jwt{
			ApiSecret: "unittest-secret",
			Issuer:    "unittest-issuer",
			Audience:  "unittest-audience",
			RsaPublic: nil,
		},
		Variables: configs.Variables{
			Port:    "",
			Timeout: 0,
		},
	}

	r.Use(LoggingMiddleware())
	r.Use(BearerAuthMiddleware(e.Jwt, tokenparser.HmacParser{
		Secret: []byte(e.Jwt.ApiSecret),
	}, tokenparser.RsaParser{}))
	r.GET("/unittest", func(c *gin.Context) {
		y := c.MustGet(utils.CALLER_USER_ROLE).([]string)
		c.String(http.StatusOK, strings.Join(y, " "))
	})

	req, _ := http.NewRequest(http.MethodGet, "/unittest", nil)
	return w, r, e, req
}

func initRsaServer(e configs.Jwt) (*httptest.ResponseRecorder, *gin.Engine, *http.Request) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	r.Use(LoggingMiddleware())
	r.Use(BearerAuthMiddleware(e, tokenparser.HmacParser{
		Secret: []byte{},
	}, tokenparser.RsaParser{
		PublicCert: e.RsaPublic,
	}))
	r.GET("/unittest", func(c *gin.Context) {
		y := c.MustGet(utils.CALLER_USER_ROLE).([]string)
		c.String(http.StatusOK, strings.Join(y, " "))
	})

	req, _ := http.NewRequest(http.MethodGet, "/unittest", nil)
	return w, r, req
}

func createHmacToken(secret string, alg string, iss string, aud string, exp int64, nbf int64) string {
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

func createRsaToken(e configs.Jwt) (*rsa.PublicKey, string) {
	key, _ := rsa.GenerateKey(crypto.Reader, 2048)
	claims := &jwt.MapClaims{
		"alg":   "HS256",
		"iss":   e.Issuer,
		"aud":   e.Audience,
		"roles": []string{"reader"},
		"exp":   jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		"nbf":   jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, _ := token.SignedString(key)
	return &key.PublicKey, tokenString
}
