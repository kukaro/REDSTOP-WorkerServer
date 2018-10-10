package api

import (
	"encoding/json"
	"github.com/labstack/echo"
	"io/ioutil"
	"net/http"
)

type UrlData struct {
	Username     string      `json:"username"`
	Name         string      `json:"name"`
	Url          string      `json:"url"`
	Method       string      `json:"method"`
	Data         interface{} `json:"data"`
	Isp          string      `json:"isp"`
	Platform     string      `json:"platform"`
	Region       string      `json:"region"`
	ResponseTime int         `json:"responsetime"`
	Status       int         `json:"status"`
};

func Routers() *echo.Echo {
	e := echo.New()
	g := e.Group("/api/v1")
	{
		g.POST("/send-url", getSendUrl)
	}
	return e
}

func getSendUrl(c echo.Context) error {
	body, _ := ioutil.ReadAll(c.Request().Body)
	jsonData := &UrlData{}
	json.Unmarshal(body, jsonData)
	return c.String(http.StatusOK, "OK")
}
