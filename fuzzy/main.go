package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	control "oksof.com/fuzzy-therm/control"
	"oksof.com/fuzzy-therm/mqtt"
)

func main() {
	godotenv.Load(".env")

	client, err := mqtt.ConnectMQTT()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(250)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	for i := range 51 {
		waterTemperature := i + 20
		CH_Temp := 60
		pelletLevel := 100

		decisionFuzzy := new(control.DecisionFuzzy)
		fuzzyStrategy := decisionFuzzy.DecideStrategy(float64(waterTemperature), float64(CH_Temp), float64(pelletLevel))

		fmt.Printf("CWU:%d,\tCO:%d,\tFS:%s\n", waterTemperature, CH_Temp, fuzzyStrategy)
	}

	<-shutdown
}
