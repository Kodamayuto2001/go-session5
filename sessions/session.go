package sessions 

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/context"
	"fmt"
)

const (
	DefaultSessionName	= "default-session"
	DefaultCookieName	= "default-cookie"
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

//	新規セッション生成
func NewSession(manager *Manager, cookieName string) *Session {
	return &Session{
		cookieName:	cookieName,
		manager:	manager,
		//	mapの初期化
		//	type{key:value,key:value,....}
		//	type{}	←からでもOK
		//	interface{}{}
		Values:		map[string]interface{}{},
	}
}

//	セッションの開始
func StartSession(sessionName, cookieName string, manager *Manager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var session *Session
		var err error 
		
		//	既存セッションの取得
		session, err = manager.Get(ctx.Request, cookieName)
		//	既存セッションを取得できなかった時
		if err != nil {
			//	セッションを新しく作成する。
			session, err = manager.New(ctx.Request, cookieName)
			if err != nil {
				fmt.Println(err.Error())
				
				//	errorが発生したときは、middlewareにおいてctx.Abort以降の処理を通さないということをする必要がある。
				//	ctx.Abort()関数で簡単に、middlewareにおいて処理を通さないことが実現できる。
				ctx.Abort()
			}
		}
		session.writer = ctx.Writer 
		ctx.Set(sessionName, session)
		defer context.Clear(ctx.Request)
		ctx.Next()
	}
}

//	デフォルトセッションの開始
func StartDefaultSession(manager *Manager) gin.HandlerFunc {
	//	constで定義してあるものをただ引数に入れて渡している
	return StartSession(DefaultSessionName, DefaultCookieName, manager)
}

//	セッションの取得
func GetSession(c *gin.Context, sessionName string) *Session {
	//	MustGet returns the value for the given key if it exists, otherwise if panics
	return c.MustGet(sessionName).(*Session)
}

//	デフォルトセッションの取得
func GetDefaultSession(c *gin.Context) *Session {
	return GetSession(c, DefaultSessionName)
}

//	セッションの保存
func (s *Session) Save() error {
	return s.manager.Save(s.request, s.writer, s)
}

//	セッション名の取得
func (s *Session) Name() string {
	return s.cookieName
}

//	セッション変数値の取得
func (s *Session) Get(key string) (interface{}, bool) {
	ret, exists := s.Values[key]
	return ret, exists
}

//	セッション変数値のセット
func (s *Session) Set(key string, val interface{}) {
	s.Values[key] = val 
}

//	セッション変数の削除
func (s *Session) Delete(key string) {
	delete(s.Values, key)
}

//	セッションの削除
func (s *Session) Terminate() {
	s.manager.Destroy(s.ID)
}