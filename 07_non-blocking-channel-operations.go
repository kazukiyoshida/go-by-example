package main

import "fmt"

func main() {
  messages := make(chan string)
  signals := make(chan bool)

  // non-blockingなselect
  select {
  case msg := <-messages:
    fmt.Println("received message", msg)
  default:
    fmt.Println("no message received")
  }

  // non-blocking な message の送信
  // messages には msg を送れない。なぜなら buffer がないから。
  // messages に make(chan string, 1) とすればmessageが送信される
  msg := "hi"
  select {
  case messages <- msg:
    fmt.Println("sent message", msg)
  default:
    fmt.Println("no message sent")
  }

  // 複数の non-blocking なselect
  select {
  case msg := <-messages:
    fmt.Println("received message", msg)
  case sig := <-signals:
    fmt.Println("received signal", sig)
  default:
    fmt.Println("no activity")
  }
}
