package main

import
(
	"fmt"
	"time"
	"errors"
	//"strconv"
)

var err error

func main(){
	err = errors.New("This is an error msg")
	
	fmt.Println(err)
}

func getNanoSecTime() int64 {
	return (time.Now().UnixNano() - (time.Now().UnixNano()/100000)*100000)
}