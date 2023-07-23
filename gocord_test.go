package gocord

import (
	"fmt"
	"os"
	"testing"

	"github.com/skifli/gocord/api"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func TestMain(t *testing.T) {
	data, err := os.ReadFile("_token")
	check(err)

	client := &Client{
		SelfBot: SelfBotToken(string(data)),
	}

	check(client.Init())

	client.AddEventHandler(api.GatewayEventNameReady, func(event *api.GatewayEventReady) {
		fmt.Printf("%#v\n\n", event)
	})

	check(client.Connect())
}
