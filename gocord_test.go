package gocord

import "testing"

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

	err = client.Run()
	check(err)
}
