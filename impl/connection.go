package impl

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

/**
 * 时间       :  19:11
 * 创建人     : Leor_Cao
 * 功能描述   :
 **/

type Connection struct {
	wsConn    *websocket.Conn
	inChan    chan []byte
	outChan   chan []byte
	closeChan chan byte
	lock      sync.Mutex
	isClosed  bool
}

func (conn *Connection) InitConnection(ws *websocket.Conn) error {
	conn.wsConn = ws
	conn.inChan = make(chan []byte, 1000)
	conn.outChan = make(chan []byte, 1000)
	conn.lock = sync.Mutex{}

	// 启动读websockte 消息的协程
	go conn.readLoop()
	// 启动写websocket 消息的协程
	go conn.writeLoop()
	return nil
}

func (conn *Connection) ReadMessage() (data []byte, err error) {
	select {
	case data = <-conn.inChan:
	case <-conn.closeChan:
		err = errors.New(fmt.Sprintf("Web socket is closed !"))
	}
	return
}

func (conn *Connection) WriteMessage(data []byte) (err error) {
	select {
	case conn.outChan <- data:
	case <-conn.closeChan:
		err = errors.New(fmt.Sprintf("Web socket is closed !"))
	}
	return
}

func (conn *Connection) Closed() {
	conn.wsConn.Close()

	conn.lock.Lock()
	if !conn.isClosed {
		conn.closeChan <- 1
		conn.isClosed = true
	}
	conn.lock.Unlock()
}

func (conn *Connection) readLoop() {
	var (
		data []byte
		err  error
	)
	for {
		if _, data, err = conn.wsConn.ReadMessage(); err != nil {
			goto ERR
		}
		conn.inChan <- data
	}
ERR:
	conn.Closed()
}

func (conn *Connection) writeLoop() {
	var (
		err error
	)
	for {
		if err = conn.wsConn.WriteMessage(websocket.TextMessage, <-conn.outChan); err != nil {
			goto ERR
		}
	}
ERR:
	conn.Closed()
}
