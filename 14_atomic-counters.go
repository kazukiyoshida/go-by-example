package main

import (
  "fmt"
  "time"
  "sync/atomic"
)

func main() {

  // managing state in Go の基本
  //
  // goroutine で並行処理を行う場合、数値のカウントアップなど
  // 単純なものでも変数の取り扱いに気をつける必要がある
  //
  // sync/atomicパッケージのAddUintを使った方法を取り上げる

  // カウントしていく変数
  var ops uint64

  // 50個のgoroutineで並行してカウントアップしていく
  for i := 0; i < 50; i++ {
    go func() {
      for {
        // アトミックなカウントアップを行うためAddUintを用いる
        // 引数にはカウンターのアドレスを与える
        atomic.AddUint64(&ops, 1)

        time.Sleep(time.Millisecond)
      }
    }()

    time.Sleep(time.Second)

    // 読み出し時にもカウントアップが行われている可能性があるので
    // 実行時点でのopsのアドレスを参照して値をcopyし、opsFinalに格納する
    opsFinal := atomic.LoadUint64(&ops)
    fmt.Println("ops: ", opsFinal)
  }
}
