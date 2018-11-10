package main

import (
  "fmt"
  "math/rand"
  "sync/atomic"
  "time"
)

// Stateful Goroutine とのインターフェース
// 読み取り
type readOp struct {
  key int
  resp chan int
}
// 書き込み
type writeOp struct {
  key int
  val int
  resp chan bool
}

func main() {
  var readOps uint64
  var writeOps uint64


  // 読み取り・書き込み用の 2 つのチャンネル
  // 
  // 足し算Goroutines, State更新Goroutine, Stateful-Goroutineの
  // 全てがこの2つのchannelにアクセスできる。
  // いわばバイパスのようなもの。
  // このバイパスを通って readOp, writeOp がやり取りされる
  reads := make(chan *readOp)
  writes := make(chan *writeOp)


  // Stateful Goroutine
  //
  // state はこの stateful-goroutine の中で private になっている
  go func() {
    var state = make(map[int]int)
    for {
      // 2つのチャンネルのどちらかにメッセージが来たら対応する
      select {
      case read := <-reads:
        read.resp <- state[read.key]
      case write := <-writes:
        state[write.key] = write.val
        write.resp <- true
      }
    }
  }()


  // ランダムに足し算を行い続けるGopherが100匹
  for r := 0; r < 100; r++ {
    go func() {
      for {
        read := &readOp{
          key: rand.Intn(5),
          resp: make(chan int)}
          reads <- read
          <-read.resp
          // stateful-goroutine と 1匹のGopherの間に
          // 専用のchannelを作成している
          atomic.AddUint64(&readOps, 1)
          time.Sleep(time.Millisecond)
      }
    }()
  }


  // ランダムにデータ追加・更新を行い続けるGopherが10匹
  for w := 0; w < 10; w++ {
    go func() {
      for {
        write := &writeOp{
          key: rand.Intn(5),
          val: rand.Intn(100),
          resp: make(chan bool)}
        writes <- write
        <-write.resp
        atomic.AddUint64(&writeOps, 1)
        time.Sleep(time.Millisecond)
      }
    }()
  }

  time.Sleep(time.Second)

  // 結果発表
  readOpsFinal := atomic.LoadUint64(&readOps)
  fmt.Println("readOps:", readOpsFinal)
  writeOpsFinal := atomic.LoadUint64(&writeOps)
  fmt.Println("writeOps:", writeOpsFinal)
}
