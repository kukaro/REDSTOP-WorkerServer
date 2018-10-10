package main

import (
	"./conf"
	"./router"
	"fmt"
)

func main() {
	if err := conf.Init(""); err == nil {
		fmt.Println("config success")
	}
	router.RunSubDomains()
}
