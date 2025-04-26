# ScafGin

[Gin](https://gin-gonic.com/) によるWebアプリ開発のテンプレートです。  
ログイン / サインアップ機能など、すぐに使える土台を用意しています。  

### 必要なツール
- **Docker**
- **make**

---

## 🚀 使い方
[webscaf](https://github.com/kodaimura/webscaf)  を使って簡単にセットアップできます。  
セットアップ後、以下のコマンドで起動できます。  
   ```bash
   make up
   ```
ログイン・サインアップ機能付きの Genie アプリが立ち上がります。  
http://localhost:8000

---

## コマンド一覧（Makefile）

```bash
make up        # コンテナ起動
make down      # コンテナ停止 & 破棄
make reup      # コンテナ停止 & 破棄 & 起動
make build     # コンテナの再ビルド
make stop      # コンテナ停止のみ
make in        # appコンテナ内にbashで入る
make log       # ログ監視
make ps        # コンテナの状態確認
```

環境を切り替えたいときは `ENV` を指定
```bash
make up ENV=prod      # 本番環境で起動
```
