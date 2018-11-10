package main

import (
  "time"
  "fmt"
)

func main() {
  timer1 := time.NewTimer(2 * time.Second)


  // TimerはチャンネルC をデフォルトで保持している
  // timer1 によるブロック
  fmt.Println("Timer 1 started !")
  <-timer1.C
  fmt.Println("Timer 1 expired")


  // time.Sleep でも上記のようなことはできるが
  // Timerは途中で破棄することができる点で優れている
  timer2 := time.NewTimer(3 * time.Second)
  go func() {
    fmt.Println("Timer 2 started !")
    <-timer2.C
    fmt.Println("Timer 2 expired")
  }()

  time.Sleep(time.Second)
  stop2 := timer2.Stop()
  if stop2 {
    fmt.Println("Timer 2 stopped")
  }
}
