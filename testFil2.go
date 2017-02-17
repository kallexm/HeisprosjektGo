package main

import
(
	"fmt"
	"time"
	//"strconv"
)

func main(){

}

func getNanoSecTime() int64 {
	return (time.Now().UnixNano() - (time.Now().UnixNano()/100000)*100000)
}