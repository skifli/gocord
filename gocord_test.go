package gocord

import (
	"fmt"
	"os"
	"testing"

	"github.com/switchupcb/dasgo/dasgo"
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

	client.AddEventHandler(dasgo.FlagGatewayEventNameMessageCreate, func(event *dasgo.Message) {
		fmt.Printf("%#v\n\n", event)
	})

	check(client.Connect())
}
