package sessions 

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/context"
)

//	大文字の場合は、外部パッケージからアクセスできる
//	小文字の場合は、外部パッケージからアクセスできない
//	Managerはmanager.goファイルのManager構造体のことを指す
type Session struct {
	cookieName	string
	ID 			string
	manager		*Manager		 
	request		*http.Request 
	writer		http.ResponseWriter
	Values 		map[string]interface{}
}

func NewSession(manager *Manager, cookieName string) *Session {
	return &Session{
		cookieName:	cookieName,
		manager:	manager,
		Values:		map[string]interface{}{},
	}
}