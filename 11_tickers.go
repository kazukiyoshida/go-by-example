package main

import (
  "time"
  "fmt"
)

func main() {
  ticker := time.NewTicker(500 * time.Millisecond)
  go func() {
    // ticker.C チャンネルに値が送られる
    // ticker.C チャンネルに値が送られる度に、range で値を取り出している
    for t := range ticker.C {
      fmt.Println("Tick at", t)
    }
  }()

  time.Sleep(1600 * time.Millisecond)
  ticker.Stop()
  fmt.Println("Ticker stopped")
}
