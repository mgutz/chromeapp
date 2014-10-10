package main

import (
	"log"

	"github.com/googollee/go-socket.io"
	"github.com/mgutz/chromeapp"
)

func listener(server *socketio.Server) {
	server.On("connection", func(conn socketio.Socket) {
		log.Println("connected")

		conn.On("echo", func(msg string) {
			conn.Emit("echo", msg)
		})

		conn.On("disconnection", func() {
			log.Println("disconnected")
		})
	})

	server.On("error", func(socket socketio.Socket, err error) {
		log.Println("ERR", err)
	})
}

func main() {
	chromeapp.Simple(listener)
}
