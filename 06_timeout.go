package main

import (
  "time"
  "fmt"
)

// select timeout パターン
// この実装のためにはchannelを通して通信の結果をやりとりする必要がある
func main() {

  fmt.Println("-------- section 1 -------")

  c1 := make(chan string, 1)
  go func() {
    // これは例えば外部通信など..
    time.Sleep(2 * time.Second)
    c1 <- "result 1"
  }()

  select {
  case res := <-c1:
    fmt.Println(res)
  case <-time.After(1 * time.Second):
    fmt.Println("timeout 1")
  }

  fmt.Println("-------- section 2 -------")
  c2 := make(chan string, 1)
  go func() {
    // これは例えば外部通信など..
    time.Sleep(2 * time.Second)
    c2 <- "result 2"
  }()

  select {
  case res := <-c2:
    fmt.Println(res)
  case <-time.After(3 * time.Second):
    fmt.Println("timeout 2")
  }
}
