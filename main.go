package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"time"
	"websocket/impl"
)

/**
 * 时间       :  18:01
 * 创建人     : Leor_Cao
 * 功能描述   :
 **/

var (
	upgrader = websocket.Upgrader{
		// 允许跨域请求
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		conn       *websocket.Conn
		connection impl.Connection
		err        error
		data       []byte
	)
	if conn, err = upgrader.Upgrade(w, r, nil); err != nil {
		return
	}
	connection.InitConnection(conn)
	go func() {
		for {
			connection.WriteMessage([]byte("Heartbeat"))
			time.Sleep(1 * time.Second)
		}
	}()
	for {
		if data, err = connection.ReadMessage(); err != nil {
			goto ERR
		}
		if err = connection.WriteMessage(data); err != nil {
			goto ERR
		}
	}

ERR:
	connection.Closed()
}

func main() {

	http.HandleFunc("/ws", wsHandler)

	http.ListenAndServe("0.0.0.0:8080", nil)
}
