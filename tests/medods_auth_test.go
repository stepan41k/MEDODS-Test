package tests

import (
	"net/http"
	"net/url"
	"testing"
	"github.com/gavv/httpexpect/v2"
)

const (
	host = "localhost:8020"
)

func TestMedodsAuth_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host: host,
	}

	e := httpexpect.Default(t, u.String())

	user := e.POST("/user/new").
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		Value("data").
		String().
		Raw()

	e.POST("/auth/new/{guid}").
		WithPath("guid", user).
		Expect().
		Status(http.StatusOK)

	e.POST("/auth/refresh").
		Expect().
		Status(http.StatusOK)
}

func TestCreate_RepeatCreateFailCase(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host: host,
	}

	e := httpexpect.Default(t, u.String())

	user := e.POST("/user/new").
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		Value("data").
		String().
		Raw()


	e.POST("/auth/new/{guid}").
		WithPath("guid", user).
		Expect().JSON().Object()

	resp := e.POST("/auth/new/{guid}").
		WithPath("guid", user).
		Expect().JSON().Object()

	if resp.ContainsKey("error") != nil {
	
		resp.Value("error").String().IsEqual("tokens already created")

		return
	}
}
		
func TestCreate_NonExistentUserFailCase(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host: host,
	}

	e := httpexpect.Default(t, u.String())

	resp := e.POST("/auth/new/{guid}").
		WithPath("guid", "non-existed-user").
		Expect().JSON().Object()

	if resp.ContainsKey("error") != nil {
	
		resp.Value("error").String().IsEqual("user not found")

		return
	}
}

func TestCreate_EmptyGUIDFailCase(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host: host,
	}

	e := httpexpect.Default(t, u.String())

	e.POST("/auth/new/{guid}").
		WithPath("guid", "").
		Expect().Status(http.StatusNotFound)
}

func TestRefresh_RefreshBeforeCreateFailCase(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host: host,
	}

	e := httpexpect.Default(t, u.String())

	resp := e.POST("/auth/refresh").
		Expect().JSON().Object()

	if resp.ContainsKey("error") != nil {
	
		resp.Value("error").String().IsEqual("first you need to generate a token")

		return
	}
}