package main

import (
	"flag"
	"log"
	"os"

	"github.com/b4t3ou/cc/service"
	"github.com/b4t3ou/cc/service/client"
	"github.com/b4t3ou/cc/service/handler"
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	c, err := client.NewStream(
		os.Getenv("API_KEY"),
		service.JSONParams{Action: "SubAdd", Subs: []string{"8~Binance~ETH~EUR"}},
	)
	if err != nil {
		log.Fatal(err)
	}

	topPosition := handler.NewTopPositions(10, 15)

	go c.Process(topPosition.Process)

	if err := c.Init(); err != nil {
		log.Fatal(err)
	}
}
