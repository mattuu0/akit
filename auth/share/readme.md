# 使い方

## .envを設定する

テンプレート
```
# 使いたいプロバイダごとに増やす
# Oauth プロバイダのキーなど
DISCORD_CLIENT_ID = 
DISCORD_CLIENT_SECRET = 
DISCORD_CALLBACK_URL = 

# ランダムな文字列 (64文字くらい)
Session_Key = 
Session_MaxAge = 2592000

# ランダムな文字列 (64文字くらい)
Transaction_Store_Secret = 
Token_Secret = 

# Redis に接続するための情報
Redis_Host = 
Redis_Port = 
Redis_Password = 
Token_Redis_DB = 0

# ランダムな文字列 (64文字くらい)
JWT_Secret = 
```

## init.go の Provider_init 関数を編集する
```go
func Provider_init() {
	goth.UseProviders(
		//ここにプロバイダを増やしていく
		discord.New(os.Getenv("DISCORD_CLIENT_ID"), os.Getenv("DISCORD_CLIENT_SECRET"), os.Getenv("DISCORD_CALLBACK_URL"), discord.ScopeIdentify, discord.ScopeEmail),
	)
}

```
