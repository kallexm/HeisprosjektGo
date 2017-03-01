package main

import
(
	"fmt"
	//"errors"
)

type buttonType int 

const (
	c0 buttonType = iota
	c2   
	c3 
)

func main() {
	fmt.Println(test(2)) 
}

func test(foo buttonType) string{
	if(foo == c0){
		return "hei"
	}else{
		return "hade"
	}
}