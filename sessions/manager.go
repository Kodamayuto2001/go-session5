/*
	セッションマネージャ
		複数のセッションを並行して管理するための機能
		1つのセッションが作られるとセッション情報を保存するための領域をメモリ上に確保する必要がある。
		それぞれのセッション情報を識別する必要がある

	実装
		セッションマネージャの構造定義
			key		interface{}	(void*型)
			value	interface{}	(void*型)
			
			個別のセッション情報をvalueに保存する
				ex).	x[key] = session

			key変数とはなにか
				複数のセッションを識別するための変数
				sessionIDなどで識別

			session変数とは何か
				ex).
				type Session struct {
					cookieName	string
					ID			string
					manager		*Manager
					request		*http.Request
					writer		http.ResponseWriter
					Values		map[string]interface{}
				}
				var session Session
				

		ユニークなセッションIDの発行
			uuidなど非常に長く、限りなく一意に近い値を生成し、返す関数
			ex).
				make([]byte,64)
				io.ReadFull(rand.Reader, b)
				base64.URLEncoding.EncodeToString(b)

		クライアント情報とセッション情報の関連付け
			クライアントから送信されるCookie情報からセッションIDを抽出する
			サーバーに保存されているセッション情報のうちクッキーから抽出したセッションIDがkeyのものを取得できるかどうか検証する
				
		セッションの生成・保存・破棄
			セッションを生成、保存、破棄する
*/

package sessions

import (
	"crypto/rand"
	"encoding/base64"
	_ "errors"
	"io"
	_ "net/http"
)

type Manager struct {
	//	小文字の場合は、外部パッケージからのアクセスができない
	//	大文字の場合は、外部パッケージからアクセスができる
	// Database map[string]interface{}
	database map[string]interface{}
}

var mg Manager

func NewManager() * Manager {
	return &mg
}

//	セッションIDの発行
func (m *Manager) NewSessionID() string {
	b := make([]byte, 64)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

//	新規セッションの生成
func (m *Manager) New(r *http.Request, cookieName string) (*Session, error) {
	//	Go言語で簡単にCookie操作ができる
	//	Cookieの属性
	//	Domain
	//	Path
	//	Expires
	//	Secure
	//	HttpOnlyなど
	cookie, err := r.Cookie(cookieName)
	//	cookie.Valueでcookieを取得する。
	if err == nil && m.Exists(cookie.Value) {
		return nil, errors.New("sessionIDはすでに発行されています")
	}

	//	session.goファイルの中のNewSession関数を使用する。
	session := NewSession(m, cookieName)
	//	IDは大文字から始まっているので、外部パッケージからでもアクセスできる
	session.ID = m.NewSessionID()
	//	requestは小文字なので、外部パッケージからアクセスできない
	session.request = r 
	
	return session, nil 
}

//	既存セッションの存在チェック
func (m *Manager) Exists(sessionID string) bool {
	_, r := m.database[sessionID]
	return r
}

//	セッション情報の保存
func (m *Manager) Save(r *http.Request, w http.ResponseWriter, session *Session) error {
	m.database[session.ID] = session 

	//	Set-Cookie header field
	c := &http.Cookie {
		//	HTTPレスポンスヘッダのSetCookieフィールドを記述
		//	Name:	cookieの名前を記述
		//	Value:	cookieの値を記述
		Name:	session.Name(),
		Value:	session.ID,
		Path:	"/",
	}

	//	サーバー側からUser AgentにSet-Cookieヘッダフィールドを設定したものを送る。
	http.SetCookie(session.writer, c)

	return session, nil 
}

//	既存セッションの取得
func (m *Manager) Get(r *http.Request, cookieName string) (*Session, error) {
	//	User Agentから送られてきたCookie header fieldを取得する。
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		//	リクエストからcookie情報を取得できない場合
		return nil ,err 
	}

	sessionID := cookie.Value
	//	cookie情報からセッション情報を取得
	//	mapの戻り値は、(value,bool)
	//	v, ok := m["apple"]
	//	appleが存在するとき、v = 値、ok = true
	//	appleが存在しないとき v = 0、ok = false
	buffer, exists := m.database[sessionID]
	if !exists {
		return nil, errors.New("無効なセッションIDです")
	}

	//	bufferという変数はinterface型
	//	interface.(type)
	//	interface型を「*Session」型にしている。
	session := buffer.(*Session)
	//	interface.(*Session)にしたので、Session構造体の変数を使うことができる
	session.request = r 
	return session, nil 
}

//	セッションの破棄
func (m *Manager) Destroy(sessionID string) {
	//	delete関数でkeyのvalueを消す。
	//	key:	sessionID
	//	map:	m.database
	//	delete(map,key)
	delete(m.database, sessionID)
}