
package main

import (
    "fmt"
    "math/rand"
    "sort"
    "time"
)

func main() {
    r := rand.New(rand.NewSource(int64(time.Now().Second())))
    zipf := rand.NewZipf(r, 2.7, 25, 300)

    data := make([]int, 0)

    N := 20
    for i := 0; i != N; i++ {
        item  :=  int(zipf.Uint64() + 30)
        data = append(data, item)
    }
    sort.Ints(data)
    fmt.Printf("%+v", data)
}
