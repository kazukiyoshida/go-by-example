package main

import (
  "bytes"
  "log"
  "net/http"
  "time"

  ws "github.com/gorilla/websocket"
)

const (
  writeWait = 10 * time.Second
  pongWait = 60 * time.Second
  pingPeriod = (pongWait * 9) / 10
  maxMessageSize = 512
)

var (
  newline = []byte{'\n'}
  space = []byte{' '}
)

var upgrader = ws.Upgrader{
  ReadBufferSize: 1024,
  WriteBufferSize: 1024,
}

type Client struct {
  hub *Hub
  conn *ws.Conn

  // 個別の client に用意された専用のチャンネル
  send chan []byte
}


func (c *Client) readPump() {
  defer func() {
    c.hub.unregister <- c
    c.conn.Close()
  }()

  c.conn.SetReadLimit(maxMessageSize)
  c.conn.SetReadDeadline(time.Now().Add(pongWait))
  c.conn.SetPongHandler(func(string) error {
    c.conn.SetReadDeadline(time.Now().Add(pongWait))
    return nil
  })

  for {
    _, message, err := c.conn.ReadMessage()
    if err != nil {
      if ws.IsUnexpectedCloseError(err, ws.CloseGoingAway, ws.CloseAbnormalClosure) {
        log.Printf("error: %v", err)
      }
      break
    }
    message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
    c.hub.broadcast <- message
  }
}

func (c *Client) writePump() {
  ticker := time.NewTicker(pingPeriod)
  defer func() {
    ticker.Stop()
    c.conn.Close()
  }()

  for {
    select {
    // Hubから個別のClient専用チャンネルsendを通してメッセージが送られた場合
    case message, ok := <-c.send:
      c.conn.SetWriteDeadline(time.Now().Add(writeWait))
      if !ok {
        c.conn.WriteMessage(ws.CloseMessage, []byte{})
        return
      }

      // WebSocket の書き込みオブジェクトを作成
      w, err := c.conn.NextWriter(ws.TextMessage)
      if err != nil {
        return
      }
      w.Write(message)

      n := len(c.send)
      for i := 0; i < n; i++ {
        // 書き込みオブジェクトを用いて、改行してからメッセージを送信
        w.Write(newline)
        w.Write(<-c.send)
      }

      if err := w.Close(); err != nil {
        return
      }

    case <-ticker.C:
      c.conn.SetWriteDeadline(time.Now().Add(writeWait))
      if err := c.conn.WriteMessage(ws.PingMessage, nil); err != nil {
        return
      }
    }

  }
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
  conn, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Println(err)
    return
  }
  client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
  client.hub.register <- client

  go client.writePump()
  go client.readPump()
}
