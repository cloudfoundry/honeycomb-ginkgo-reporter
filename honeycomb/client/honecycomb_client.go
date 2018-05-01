package client

import "github.com/honeycombio/libhoney-go"

//go:generate counterfeiter . Client
type Client interface {
	SendEvent(data interface{}) error
}

type honeyCombClient struct {
	config libhoney.Config
}

func New(config libhoney.Config) honeyCombClient {
	return honeyCombClient{config: config}
}

// We created this because the way the honey comb go client is written makes
// it impossible to test in a reasonable way
// We need the init method to return a client so that we can test it with fakes.
func (hc honeyCombClient) SendEvent(data interface{}) error {
	err := libhoney.Init(hc.config)
	defer libhoney.Close()

	if err != nil {
		return err
	}

	return libhoney.SendNow(data)
}
