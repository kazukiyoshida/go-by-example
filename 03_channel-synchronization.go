package main

import (
  "fmt"
  "time"
)

func worker(done chan bool) {
  fmt.Print("working...")
  time.Sleep(time.Second)
  fmt.Println("done")

  done <- true
}

func main() {
  done := make(chan bool, 1)
  go worker(done)

  // <-done で、doneから通知が来るまでbolckされる
  // つまり、done channel は他のgoroutineへの通知機能を果たす
  <-done
}
