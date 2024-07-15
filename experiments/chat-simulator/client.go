package main

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

type testClient struct {
	received []string
	name     string
	cs       *client

	ctx       context.Context
	cancelCtx func()
}

func newTestClient(name string) *testClient {
	ctx, cancel := context.WithCancel(context.Background())
	return &testClient{
		name:      name,
		ctx:       ctx,
		cancelCtx: cancel,
	}
}

func (c *testClient) enter(s *server) error {
	cs, err := s.connectClient(c.name)
	if err != nil {
		return fmt.Errorf("connect to server: %w", err)
	}
	c.cs = cs

	go c.listenToServer()
	return nil
}

func (c *testClient) listenToServer() {
	ticker := time.NewTicker(time.Millisecond * 200)

	for {
		select {
		case <-c.ctx.Done():
			log.Debug().Msgf("%s: \t Disconnected ", c.name)
			return
		case <-ticker.C:
			// c.cs.chMessage <- message{
			// 	messageType: ping,
			// }
			//log.Debug().Msgf("%s: \t sent ping", c.name)
		case msg := <-c.cs.chToClient:
			c.received = append(c.received, msg.content)
			log.Debug().Msgf("%s: \t received message: \t %s", c.name, msg.content)
		}
	}
}

func (c *testClient) leave() {
	c.cs.chToServer <- message{
		messageType: disconnect,
	}
	log.Debug().Msgf("%s: \t sent disconnect ", c.name)

	c.cancelCtx()
}

func (c *testClient) sendPM(to string, content string) {
	c.cs.chToServer <- message{
		messageType: private,
		content:     content,
		to:          to,
		from:        c.cs.id,
	}

	log.Debug().Msgf("%s: \t sent pm \t to: %s \t content: %s", c.name, to, content)
}

func (c *testClient) broadcast(content string) {
	c.cs.chToServer <- message{
		messageType: public,
		content:     content,
		from:        c.cs.id,
	}
	log.Debug().Msgf("%s: \t sent \t to all \t content: %s", c.name, content)
}
