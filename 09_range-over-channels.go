package main

import "fmt"

func main() {
  queue := make(chan string, 2)
  queue <- "one"
  queue <- "two"
  close(queue)

  // queueチャネルはcloseされても値を保持し続ける
  // チャネルのバッファに対して range がかけられる
  for elem := range queue {
    fmt.Println(elem)
  }
  // 一度取り出されたら queue は空になっている
  for elem := range queue {
    fmt.Println(elem)
  }
}
