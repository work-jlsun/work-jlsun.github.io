package main

import (
    "fmt"
    "math/rand"
    "sort"
    "time"
    "strconv"
)


func sliceAvg( s  []int ) int {
    sum := 0
    for i := 0; i != len(s); i++ {
        sum += s[i]
    }
    return sum/len(s)
}

func main() {
    r := rand.New(rand.NewSource(int64(time.Now().Second())))
    zipf := rand.NewZipf(r, 2.7, 25, 300)
    data := make([]int, 0)
    data2 := make([]int, 0)
    data3 := make([]int, 0)

    data4 := make([]int, 0)
    data5 := make([]int, 0)

    N := 100000
    for i := 0; i != N; i++ {

        item  :=  int(zipf.Uint64() + 30)
        item2  :=  int(zipf.Uint64() + 30)
        item3  :=  int(zipf.Uint64() + 30)

        data = append(data, item)
        data2 = append(data2, item2)
        data3 = append(data3, item3)

        items := make([]int,0)
        items = append(items, item)
        items = append(items, item2)
        items = append(items, item3)
        sort.Ints(items)

        //get the biggest in three
        data4 = append(data4, items[2])

        //get the sendond  bigger one
        data5 = append(data5, items[1])
    }

    sort.Ints(data)
    sort.Ints(data2)
    sort.Ints(data3)
    sort.Ints(data4)
    sort.Ints(data5)

    str := ""
    str2 := ""
    str3 := ""
    str4 := ""
    str5 := ""

    // dist
    for  i := 1; i <= 20; i++ {
        index  := (N/100)*5*i
        str = str + ","  + strconv.Itoa(data[index-1])
        str2 = str2 + ","  + strconv.Itoa(data2[index-1])
        str3 = str3 + ","  + strconv.Itoa(data3[index-1])
        str4 = str4 + ","  + strconv.Itoa(data4[index-1])
        str5 = str5 + ","  + strconv.Itoa(data5[index-1])
        //fmt.Printf(" %d percent:  %d\n", 5*i,data[index-1])
    }

    // avg data
    fmt.Printf("%d\n", sliceAvg(data))
    fmt.Printf("%d\n", sliceAvg(data2))
    fmt.Printf("%d\n", sliceAvg(data3))
    fmt.Printf("%d\n", sliceAvg(data4))
    fmt.Printf("%d\n", sliceAvg(data5))

    fmt.Printf("%s\n", str)
    fmt.Printf("%s\n", str2)
    fmt.Printf("%s\n", str3)
    fmt.Printf("%s\n", str4)
    fmt.Printf("%s\n", str5) 
}
