package main

import (
	"fmt"

	"github.com/robbykun/languageGPRestApi/api"
	"github.com/robbykun/languageGPRestApi/db"
)

func main() {
	fmt.Println("main開始")
	db.Init()
	api.Init()
}
