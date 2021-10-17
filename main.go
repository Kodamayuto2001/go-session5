package main

import (
	"fmt"
	"go-session5/sessions"
)

func main() {
	fmt.Println("ok")
	var x sessions.Manager
	
	//	初期化が必要（初期化しないとエラーが出る）
	// x.Database = map[string]interface{}{}
	// x.Database = map[string]interface{}{
	// 	"hoge":"hoge",
	// }
	// fmt.Println(x)

	fmt.Println(sessions.NewManager())
	fmt.Println(x.NewSessionID())
}