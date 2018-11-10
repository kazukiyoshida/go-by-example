package main

import (
  "fmt"
  "time"
)

// worker
//
// jobチャンネルにある値を取得し、2倍にしてresultへ返し、
// それが終わったら再びjobチャンネルから値を取り出し..という作業を行う
func worker(id int, jobs <-chan int, results chan<- int) {
  for j := range jobs {
    fmt.Println("worker", id, "start job", j)
    time.Sleep(time.Second)
    fmt.Println("worker", id, "fin job", j)
    results <- j * 2
  }
}

func main() {
  jobs := make(chan int, 100) // workerに処理してもら数字を流し込むチャネル
  results := make(chan int, 100) // workerが結果をこのresultsチャネルに流し込む

  // workerを3つ起動する
  // 起動した直後は jobs に特に何の値も入っていないので仕事をしない
  for w := 1; w <= 3; w++ {
    go worker(w, jobs, results)
  }

  // jobsチャンネルに処理してもらう数字を流し込む
  // closeすることで全てのjobが流し込まれたことを意味する
  for j := 1; j <= 5; j++ {
    jobs <- j
  }
  close(jobs)

  // resultsから通知が5回くるまでブロックする
  for a := 1; a <= 5; a ++ {
    <-results
  }
}
