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
// func (m *Manager) New(r *http.Request, cookieName string) (*Session, error) {
// 	cookie, err := r.Cookie(cookieName)
// }