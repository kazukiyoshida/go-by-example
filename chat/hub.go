package main

type Hub struct {
  // client一覧
  clients map[*Client]bool

  // clients からのメッセージ バイパス
  broadcast chan []byte

  // client からの登録要求バイパス
  register chan *Client

  // client からの登録解除要求
  unregister chan *Client
}

func newHub() *Hub {
  return &Hub{
    broadcast:  make(chan []byte),
    register:   make(chan *Client),
    unregister: make(chan *Client),
    clients:    make(map[*Client]bool),
  }
}

func (h *Hub) run() {
  for {
    select {

    // clientからの登録要求
    case client := <-h.register:
      h.clients[client] = true

    // clientからの登録解除要求
    case client := <-h.unregister:
      if _, ok := h.clients[client]; ok {
        delete(h.clients, client)
        close(client.send)
      }

    // clientからの、broadcastチャンネルを通したメッセージ
    // 個別のチャンネルの send チャンネルにメッセージを送り返す
    case message := <-h.broadcast:
      for client := range h.clients {
        select {
        case client.send <- message:
        default:
          close(client.send)
          delete(h.clients, client)
        }
      }
    }
  }
}
