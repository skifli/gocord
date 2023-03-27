package gocord

import (
	"testing"
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

	check(client.Connect())
}
