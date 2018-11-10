package main

import (
  "fmt"
  "time"
)

func main() {
  messages := make(chan string)

  go func() {
    messages <- "ping"
  }()

  // 2つあると pong にしかならない
  // By default sends and receives block
  // until both the sender and receiver are ready.
  // なので、全ての sends が終わるまで待っているぽい
  // 
  //go func() {
  //  messages <- "pong"
  //  time.Sleep(2)
  //}()

  msg := <-messages
  fmt.Println(msg)
}
