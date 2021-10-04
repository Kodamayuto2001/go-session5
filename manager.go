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
		セッションの生成・保存・破棄
*/