package main

import (
	"./conf"
	"./router"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"net/url"
)

type VisitData struct {
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
}

const (
	EXCHANGE_NAME string = "myExchange"
)

func main() {

	go func() {
		rabbitmqRec()
	}()

	if err := conf.Init(""); err == nil {
		fmt.Println("config success")
	}
	router.RunSubDomains()
}

func rabbitmqRec() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		EXCHANGE_NAME,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,
		"",
		EXCHANGE_NAME,
		false,
		nil,
	)
	failOnError(err, "failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			visitData := VisitData{}
			json.Unmarshal(d.Body, &visitData)
			//log.Println(visitData)
			visitData.Isp = url.PathEscape(visitData.Isp)
			visitData.Platform = url.PathEscape(visitData.Platform)
			visitData.Region = url.PathEscape(visitData.Region)

			bodyString := `testresult,isp=%s,platform=%s,region=%s name="%s",url="%s",status=%d,method="%s",responsetime=%d`
			bodyString = fmt.Sprintf(bodyString, visitData.Isp, visitData.Platform, visitData.Region, visitData.Name, visitData.Url, visitData.Status, visitData.Method, visitData.ResponseTime)
			resp, err := http.Post("http://localhost:8086/write?db=redstop", "", bytes.NewBuffer([]byte(bodyString)))
			log.Println(resp, &err == nil)
			log.Println(" [x] receive", !(&d == nil))
		}
	}()
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
