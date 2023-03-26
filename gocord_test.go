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
		User: UserToken(""),
	}

	err := client.Init()
	check(err)

	client.AddEventHandler(api.GatewayEventNameHello, func(event *api.GatewayEventHello) {
		fmt.Printf("%#v\n", event) // Testing the event handler
	})

	err = client.Run()
	check(err)
}
