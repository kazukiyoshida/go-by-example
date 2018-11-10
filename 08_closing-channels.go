package main

import "fmt"
import "time"

func main() {
  fmt.Println("====== start =====")

  jobs := make(chan int, 5)
  done := make(chan bool)

  // Worker
  //
  // このworker =「func() の goroutine」は jobsチャンネルから
  // 値を受信するまで待ち続ける
  go func() {
    for {
      // jobsから2つの返り値を受ける
      // もしも jobsチャンネルがcloseしていたらmoreはfalseになる
      j, more := <-jobs
      if more {
        fmt.Println("received job", j)
      } else {
        fmt.Println("received all jobs")
        done <- true
        return
      }
    }
  }()

  // 本体側では jobチャンネルを通してworkerにメッセージを送り、
  // 終了したらチェンネルをクローズする
  for j := 1; j <= 3; j++ {
    jobs <- j
    fmt.Println("sent job", j)
    time.Sleep(time.Second)
  }

  close(jobs)
  fmt.Println("sent all jobs")

  <-done
}
