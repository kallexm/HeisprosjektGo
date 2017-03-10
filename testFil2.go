package main

import
(
	"fmt"
	"time"
)



func main(){
	var str string
	chan2 := make(chan string)
	
	go thread1(chan2)
	
	for {
		select {
		case chan2 <- "verden":
		
		case str = <- chan2:
			fmt.Println(str)
		//default:
			
		}
		time.Sleep(500*time.Millisecond)
	}
	
	
}


func getNanoSecTime() int64 {
	return (time.Now().UnixNano() - (time.Now().UnixNano()/100000)*100000)
}


func thread1(ch chan string) {
	var str string
	for {
		select {
		case ch <- "Hei":
			
		case str = <- ch:
			fmt.Println(str)
		default:
			//time.Sleep(2000*time.Millisecond)
		}
		
	}
	
}