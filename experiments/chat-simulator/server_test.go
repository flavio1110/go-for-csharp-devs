package main

import (
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel)

	t.Run("clients receives welcome message", func(t *testing.T) {
		s := newServer()

		client1 := newTestClient("client1")
		err := client1.enter(s)
		require.NoError(t, err)

		client2 := newTestClient("client2")
		err = client2.enter(s)
		require.NoError(t, err)

		client3 := newTestClient("client3")
		err = client3.enter(s)
		require.NoError(t, err)

		time.Sleep(time.Millisecond * 100)
		assert.Equal(t, []string{"FROM \"server\": Hello client1!!! Welcome to the chat."}, client1.received)
		assert.Equal(t, []string{"FROM \"server\": Hello client2!!! Welcome to the chat."}, client2.received)
		assert.Equal(t, []string{"FROM \"server\": Hello client3!!! Welcome to the chat."}, client3.received)
	})

	t.Run("broadcast messaage with a single client", func(t *testing.T) {
		s := newServer()

		client1 := newTestClient("client1")
		err := client1.enter(s)
		require.NoError(t, err)

		client1.broadcast("hello from 1")

		time.Sleep(time.Millisecond * 300)

		assert.ElementsMatch(t, []string{
			"FROM \"server\": Hello client1!!! Welcome to the chat.",
			`FROM "client1": hello from 1`},
			client1.received)
	})

	t.Run("broadcast messaage", func(t *testing.T) {
		s := newServer()

		client1 := newTestClient("client1")
		err := client1.enter(s)
		require.NoError(t, err)

		client2 := newTestClient("client2")
		err = client2.enter(s)
		require.NoError(t, err)

		client3 := newTestClient("client3")
		err = client3.enter(s)
		require.NoError(t, err)

		client1.broadcast("hello from 1")

		time.Sleep(time.Millisecond * 100)

		assert.ElementsMatch(t, []string{
			"FROM \"server\": Hello client1!!! Welcome to the chat.",
			"FROM \"client1\": hello from 1"},
			client1.received)

		assert.ElementsMatch(t, []string{
			"FROM \"server\": Hello client2!!! Welcome to the chat.",
			"FROM \"client1\": hello from 1"},
			client2.received)

		assert.ElementsMatch(t, []string{
			"FROM \"server\": Hello client3!!! Welcome to the chat.",
			"FROM \"client1\": hello from 1"},
			client3.received)
	})

	t.Run("multiple broadcast messaages", func(t *testing.T) {
		s := newServer()

		client1 := newTestClient("client1")
		err := client1.enter(s)
		require.NoError(t, err)

		client2 := newTestClient("client2")
		err = client2.enter(s)
		require.NoError(t, err)

		client3 := newTestClient("client3")
		err = client3.enter(s)
		require.NoError(t, err)

		client1.broadcast("hello from 1")
		client2.broadcast("hello from 2")
		client3.broadcast("hello from 3")

		time.Sleep(time.Millisecond * 100)

		assert.ElementsMatch(t, []string{
			"Hello client1!!! Welcome to the chat.",
			"hello from 1",
			"hello from 2",
			"hello from 3"},
			client1.received)

		assert.ElementsMatch(t, []string{
			"Hello client2!!! Welcome to the chat.",
			"hello from 1",
			"hello from 2",
			"hello from 3"},
			client2.received)

		assert.ElementsMatch(t, []string{
			"Hello client3!!! Welcome to the chat.",
			"hello from 1",
			"hello from 2",
			"hello from 3"},
			client3.received)
	})

	t.Run("send private message to someone not connected", func(t *testing.T) {
		s := newServer()

		client1 := newTestClient("client1")
		err := client1.enter(s)
		require.NoError(t, err)

		client1.sendPM("AN_ID", "Hello 1")

		time.Sleep(time.Millisecond * 100)

		assert.ElementsMatch(t, []string{
			"Hello client1!!! Welcome to the chat.",
			"Hello 1",
			`"AN_ID" is not in the chat`},
			client1.received)

	})

	t.Run("send private message", func(t *testing.T) {
		s := newServer()

		client1 := newTestClient("client1")
		err := client1.enter(s)
		require.NoError(t, err)

		client2 := newTestClient("client2")
		err = client2.enter(s)
		require.NoError(t, err)

		client3 := newTestClient("client3")
		err = client3.enter(s)
		require.NoError(t, err)

		client1.sendPM(client2.cs.id, "hello from 1")
		client2.sendPM(client1.cs.id, "hello from 2")

		time.Sleep(time.Millisecond * 100)

		assert.ElementsMatch(t, []string{
			"Hello client1!!! Welcome to the chat.",
			"hello from 1",
			"hello from 2"},
			client1.received)

		assert.ElementsMatch(t, []string{
			"Hello client2!!! Welcome to the chat.",
			"hello from 2",
			"hello from 1"},
			client2.received)

		assert.ElementsMatch(t, []string{
			"Hello client3!!! Welcome to the chat."},
			client3.received)
	})

}
