package api

import (
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/streadway/amqp"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	EXCHANGE_NAME string = "myExchange"
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

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func getSendUrl(c echo.Context) error {
	body, _ := ioutil.ReadAll(c.Request().Body)
	jsonData := &UrlData{}
	json.Unmarshal(body, jsonData)

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to declare an exchange")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		EXCHANGE_NAME,
		"fanout",
		true,
		false,
		false,
		false,
		nil)
	failOnError(err, "Failed to declare an exchange")

	err = ch.Publish(
		EXCHANGE_NAME,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] sent %s", "OK")

	return c.String(http.StatusOK, "OK")
}
