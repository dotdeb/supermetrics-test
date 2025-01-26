package controllers

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/dotdeb/supermetrics-test/internals/data"
	"github.com/dotdeb/supermetrics-test/internals/utils"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Users []data.User
}

func TestSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set(utils.LOG_ID, "unittest")
	c.Set(utils.CALLER_USER_ROLE, []string{"reader"})
	c.Set(utils.GET_USERS_FUNC, func() []data.User {
		return []data.User{
			{
				Id:       "1",
				Username: "john.doe",
			},
		}
	})
	GetUsers(c)

	defer w.Result().Body.Close()
	b, _ := io.ReadAll(w.Body)
	if w.Code != 200 {
		t.Error(w.Code, string(b))
	}

	data := Response{}
	json.Unmarshal(b, &data)
	if data.Users[0].Id != "1" {
		panic("Response id is not correct")
	}
	if data.Users[0].Username != "john.doe" {
		panic("Response username is not correct")
	}
}

func TestWrongRole(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set(utils.LOG_ID, "unittest")
	c.Set(utils.CALLER_USER_ROLE, []string{"writer"})
	GetUsers(c)
	if w.Code != 403 {
		t.Error(w.Code, "Wrong user role is passing")
	}
}

func TestEmptyRole(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set(utils.LOG_ID, "unittest")
	c.Set(utils.CALLER_USER_ROLE, []string{""})
	GetUsers(c)
	if w.Code != 403 {
		t.Error(w.Code, "Empty user role is passing")
	}
}
