package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type server struct {
	mClients sync.RWMutex
	clients  []*client
}

type messageType int

const (
	public = iota
	private
	ping
	disconnect
)

type message struct {
	from        string
	to          string
	content     string
	messageType messageType
}

type client struct {
	id         string
	name       string
	lastPing   time.Time
	chToClient chan message
	chToServer chan message

	ctx       context.Context
	cancelCtx func()
}

func newServer() *server {
	return &server{
		mClients: sync.RWMutex{},
	}
}

func (s *server) connectClient(name string) (*client, error) {
	ctx, cancel := context.WithCancel(context.Background())

	c := &client{
		id:         uuid.New().String(),
		name:       name,
		lastPing:   time.Now().UTC(),
		chToServer: make(chan message),
		chToClient: make(chan message),
		ctx:        ctx,
		cancelCtx:  cancel,
	}

	func() {
		s.mClients.Lock()
		defer s.mClients.Unlock()
		s.clients = append(s.clients, c)
	}()
	go s.listenToClient(c)

	return c, nil
}

func (s *server) Stop() {
	// s.broadcastMessage(message{
	// 	messageType: public,
	// 	content:     "bye",
	// })

	s.mClients.Lock()
	defer s.mClients.Unlock()

	for _, c := range s.clients {
		close(c.chToClient)
		close(c.chToServer)
	}

	s.clients = nil

}

func (s *server) disconnetClient(c *client) {
	s.mClients.Lock()
	defer s.mClients.Unlock()

	for ix, cs := range s.clients {
		if cs.id == c.id {
			s.clients = append(s.clients[:ix], s.clients[ix+1:]...)

			s.broadcastMessage(cs, message{
				from:        "server",
				content:     fmt.Sprintf("%s has left.", c.name),
				messageType: private,
			})

			break
		}
	}
}

func (s *server) listenToClient(c *client) {
	welcome := message{
		from:        "server",
		to:          c.id,
		content:     fmt.Sprintf("Hello %s!!! Welcome to the chat.", c.name),
		messageType: private,
	}

	if err := s.sendMessage(welcome); err != nil {
		log.Err(err).Msg("send welcome message")
		return
	}

	for {
		select {
		case <-c.ctx.Done():
			return
		case msg := <-c.chToServer:
			log.Debug().Msgf("Received nessage from %s\t content: %s", c.name, msg.content)
			msg.from = c.id
			if msg.messageType == ping {
				c.lastPing = time.Now().UTC()
			} else if msg.messageType == disconnect {
				s.disconnetClient(c)
			} else if msg.messageType == private {
				s.sendMessage(msg)
			} else if msg.messageType == public {
				s.broadcastMessage(c, msg)
			}
		}
	}
}

func (s *server) sendMessage(m message) error {
	if m.messageType != private {
		return errors.New("sending non private message")
	}

	s.mClients.RLock()
	defer s.mClients.RUnlock()

	foundRecipient := false

	for _, c := range s.clients {
		if !foundRecipient {
			foundRecipient = c.id == m.to
		}

		if c.id == m.to {
			log.Debug().Msgf("Sending message TO: %q \t content: %q", m.to, m.content)
			m.content = fmt.Sprintf("FROM %q: %s", m.from, m.content)
			c.chToClient <- m
		}

		if c.id == m.from {
			log.Debug().Msgf("Sending message TO: %q \t content: %q", m.to, m.content)
			m.content = fmt.Sprintf("TO %q: %s", m.to, m.content)
			c.chToClient <- m
		}
	}

	if !foundRecipient {
		s.sendMessage(message{
			from:        "server",
			to:          m.from,
			messageType: private,
			content:     fmt.Sprintf("%q is not in the chat", m.to),
		})

		return errors.New("recipient not found")
	}
	return nil
}

func (s *server) broadcastMessage(cf *client, m message) error {
	if m.messageType != public {
		return errors.New("broadcasting non public message")
	}

	s.mClients.RLock()
	defer s.mClients.RUnlock()

	for _, c := range s.clients {
		log.Debug().Msgf("Sending public message TO: %q \t content: %q", m.to, m.content)
		m.content = fmt.Sprintf("FROM %q: %s", cf.name, m.content)
		c.chToClient <- m
	}

	return nil
}
