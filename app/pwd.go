package main

import (
	"fmt"
	"os"
)

// var wd = getWorkingDir()

// func getWorkingDir() string {
// 	pwd, err := os.Getwd()
// 	if err != nil {
// 		fmt.Println("Error getting working directory:", err)
// 		os.Exit(1)
// 	}
//
// 	return pwd
// }

type pwd struct{}

func (p pwd) execute(args []string) {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory:", err)
		os.Exit(1)
	}

	fmt.Println(pwd)
}
