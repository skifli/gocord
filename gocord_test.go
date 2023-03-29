package gocord

import (
	"fmt"
	"testing"

	"github.com/skifli/gocord/api"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func TestMain(t *testing.T) {
	client := &Client{
		SelfBot: SelfBotToken(""),
	}

	check(client.Init())

	client.AddEventHandler(api.GatewayEventNameReady, func(event *api.GatewayEventReady) {
		fmt.Printf("%#v\n", event)
	})

	check(client.Connect())
}
