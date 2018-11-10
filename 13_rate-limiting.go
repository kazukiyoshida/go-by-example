package main

import (
  "time"
  "fmt"
)

// rate limiting = 律速制御
func main() {

  fmt.Println("===== 律速制御なし ====")

  requests := make(chan int, 10)
  for i := 1; i <= 10; i++ {
    requests <- i
  }
  close(requests)

  for req := range requests {
    fmt.Println("request", req, time.Now())
  }



  fmt.Println("===== 律速制御あり(Limitter) ====")

  limiter := time.Tick(400 * time.Millisecond)

  requests2 := make(chan int, 10)
  for i := 1; i <= 10; i++ {
    requests2 <- i
  }
  close(requests2)

  for req := range requests2 {
    <-limiter
    fmt.Println("request", req, time.Now())
  }



  fmt.Println("===== 律速制御あり(burstyLimitter) ====")
  // 最初の3つのリクエストは一度に処理し、そのあとのリクエストは
  // 等間隔に処理を行っていく

  burstyLimitter := make(chan time.Time, 3)

  for i := 0; i < 3; i++ {
    burstyLimitter <- time.Now()
  }

  go func() {
    for t := range time.Tick(400 * time.Millisecond) {
      burstyLimitter <- t
    }
  }()

  burstyRequests := make(chan int, 10)
  for i := 1; i <= 10; i++ {
    burstyRequests <- i
  }
  close(burstyRequests)

  for req := range burstyRequests {
    <-burstyLimitter
    fmt.Println("request", req, time.Now())
  }
}
