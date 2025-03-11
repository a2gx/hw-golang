package main

import (
	"fmt"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/config"
)

func main() {
	conf := config.GetSenderConfig()
	fmt.Printf("%+v\n", conf)
}
