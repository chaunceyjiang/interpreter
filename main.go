package main

import (
	"Interpreter/repl"
	"fmt"
	"os"
	"os/user"
)

func main() {
	u,err:=user.Current()
	if err!=nil{
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Monkey programming language!\n",u.Username)
	fmt.Printf("Feel free to type commands\n")
	repl.Start(os.Stdin,os.Stdout)
}
