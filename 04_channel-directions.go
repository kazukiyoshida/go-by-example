package main

import "fmt"

// channel への送受信方向
//
// 関数の中では送信専用のchannel、受信専用のchannelを
// 引数にできる。（channel自体にその機能はない？）
// 方向を間違えるとcompile-time errorになる

// 送信するためだけに channel を受け取る関数
func ping(pings chan<- string, msg string) {
  pings <- msg
}

// 受信するためだけに channel を受け取る関数
func pong(pings <-chan string, pongs chan<- string) {
  msg := <-pings
  pongs <- msg
}

func main() {
  pings := make(chan string, 1)
  pongs := make(chan string, 1)

  ping(pings, "passed message")
  pong(pings, pongs)
  fmt.Println(<-pongs)
}
