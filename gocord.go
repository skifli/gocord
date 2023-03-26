package gocord

import (
	"errors"

	"github.com/skifli/gocord/api"
)

// Client represents a Discord self-bot.
type Client struct {
	Gateway *api.Gateway // Gateway contains data relating to a Discord WebSocket connection.
	User    *api.User    // User contains data relating to the user.
}

// Creates a user struct, used when initializing a Client struct.
func UserToken(token string) *api.User {
	return &api.User{
		Discriminator: "",
		ID:            "",
		Locale:        "",
		Token:         token,
		Username:      "",
	}
}

func (client *Client) Init() error {
	if client.User == nil {
		return errors.New("user not supplied")
	} else if client.User.Token == "" {
		return errors.New("user's token not supplied")
	}

	if err := client.User.GetData(); err != nil {
		return err
	}

	client.Gateway = api.CreateGateway(client.User)

	return nil
}

func (client *Client) Run() error {
	return client.Gateway.Connect()
}

func (client *Client) AddEventHandler(event api.GatewayEventName, function any) error {
	return client.Gateway.Handlers.Add(event, function)
}
