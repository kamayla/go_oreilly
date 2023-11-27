package main

import (
	"chat/trace"
	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
	"log"
	"net/http"
)

type room struct {
	forward chan *message
	join    chan *client
	leave   chan *client
	clients map[*client]bool
	tracer  trace.Tracer
}

func newRoom() *room {
	return &room{
		// forwardは他のクライアントに転送するためのメッセージを保持するチャネル
		forward: make(chan *message),
		// joinはチャットルームに参加しようとしているクライアントのためのチャネル
		join: make(chan *client),
		// leaveはチャッルームを退出しようとしているクライアントのためのチャネル
		leave: make(chan *client),
		// clientsは在室している全てのクライアントが保持される
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// roomに参加
			r.clients[client] = true
			r.tracer.Trace("新しいクライアントが参加しました")
		case client := <-r.leave:
			// roomを退出
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("クライアントが退出しました")
		case msg := <-r.forward:
			r.tracer.Trace("メッセージを受信しました: ", msg.Message)
			// 全てのクライアントにメッセージを転送
			for client := range r.clients {
				select {
				case client.send <- msg:
					r.tracer.Trace("-- クライアントに送信されました")
				// メッセージを送信
				default:
					// 送信に失敗
					delete(r.clients, client)
					close(client.send)
					r.tracer.Trace("-- 送信に失敗しました。クライアントをクリーンアップします")
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: messageBufferSize,
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	authCoolie, err := req.Cookie("auth")

	if err != nil {
		log.Fatalln(err)
		return
	}

	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCoolie.Value),
	}

	r.join <- client

	defer func() { r.leave <- client }()

	go client.write()
	client.read()
}
