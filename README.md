# goutils

自分用 Go ユーティリティ集。破壊的変更あり。

## インポート

```
github.com/itoken417/goutils/<パッケージ名>
```

---

## chrome

`github.com/sclevine/agouti` を使ったHeadless Chrome操作ラッパー。

**前提:** `google-chrome-stable` がインストール済みであること。

| 関数 | 説明 |
|------|------|
| `Init()` | ChromeDriverを起動する |
| `GetDriver()` | `*agouti.WebDriver` を返す |
| `GetNewPage(sslCheck bool)` | 新しいページを開く。`false` でSSL検証をスキップ |
| `Destroy()` | ChromeDriverを停止する |

---

## chromedp

`github.com/chromedp/chromedp` を使ったHeadless Chrome操作ラッパー。

**前提:** `google-chrome-stable` がインストール済みであること。

| 関数 | 説明 |
|------|------|
| `Init()` | グローバルコンテキストを初期化する |
| `Destroy()` | グローバルコンテキストを破棄する |
| `GetChromeDPContext()` | グローバルコンテキストを返す（未初期化なら自動Init） |
| `GetNewPage(sslCheck bool)` | 新しいページコンテキストを返す。`false` でSSL検証をスキップ |
| `NewChromeDPContext(parent, sslCheck)` | 独立したコンテキストを生成する |
| `NewChromeDPContextWithTimeout(parent, timeout, sslCheck)` | タイムアウト付きコンテキストを生成する |
| `RunChromeDP(ctx, actions...)` | 指定コンテキストでアクションを実行する |
| `RunDefaultChromeDP(actions...)` | グローバルコンテキストでアクションを実行する |

---

## converter

`github.com/ktnyt/go-moji` の文字変換ラッパー。よく使うものだけ収録。

| 関数 | 説明 |
|------|------|
| `HK2ZK(str string) string` | 半角カタカナ → 全角カタカナ |
| `ZE2HE(str string) string` | 全角英数 → 半角英数 |

---

## regex

よく使う正規表現操作をまとめたユーティリティ。

| 関数 | 説明 |
|------|------|
| `RmAllSpace(str string) string` | 半角・全角スペースをすべて除去する |
| `RmSpace(str string) string` | 半角スペースを除去する |
| `RmZenSpace(str string) string` | 全角スペースを除去する |
| `RegC(rule, str string) bool` | 正規表現にマッチするか確認する |
| `RegM(rule, str string) []string` | 正規表現にマッチした文字列をすべて返す |
| `RegSM(rule, str string) [][]string` | サブマッチを含む全マッチを返す |
| `RegR(rule, repl, str string) string` | 正規表現にマッチした箇所を置換する |

---

## logger

見やすさ優先のカスタムロガー。

| 関数 | 説明 |
|------|------|
| `Init(isDebug, writeCaller bool) bool` | ロガーを初期化する。`isDebug=false` でリリースモード（ファイル出力） |
| `Destory()` | ログファイルをクローズする |
| `Log(a ...interface{})` | 通常ログを出力する |
| `Dump(a ...interface{})` | `%#v` 形式で詳細ログを出力する |
| `ErrLog(a ...interface{})` | エラーログを出力してpanicする |
| `ErrCheck(err error)` | `err != nil` の場合にpanicする |

**モード:**
- `isDebug=true` (デバッグ): 標準出力へ出力
- `isDebug=false` (リリース): `./logs/<呼び出し元ディレクトリ名>.log` へ書き込み（上書き）

**writeCaller=true** にすると各ログに呼び出し元のファイル名・行番号が付く（出力が増えるため通常は `false` 推奨）。

---

## mailsender

SMTPメール送信パッケージ。`base64` / `quoted-printable` エンコーディングに対応。

```go
cfg := mailsender.Config{
    Host:     "smtp.example.com",
    Port:     "587",
    User:     "user@example.com",
    Password: "password",
    From:     "user@example.com",
    Encoding: "base64", // "base64" / "quoted-printable" / "" (未設定)
}
s := mailsender.New(cfg)
err := s.Send([]string{"to@example.com"}, "件名", "本文")
```

| 型/関数 | 説明 |
|---------|------|
| `Config` | SMTP接続設定を保持する構造体 |
| `New(cfg Config) *Sender` | Senderを生成する |
| `(*Sender).Send(to []string, subject, body string) error` | メールを送信する |

---

## dumper

PerlのData::Dumperに相当する、デバッグ用ダンプユーティリティ。
任意の値をインデント付きの人間が読みやすい形式で出力する。

```go
dumper.Dd(42, "hello", []int{1, 2, 3})
// $var1 = 42
// $var2 = "hello"
// $var3 = []int{
//   1,
//   2,
//   3,
// }

fmt.Print(dumper.DumpNamed("person", p))
// $person = Person{
//   Name: "田中",
//   Age: 30,
//   ...
// }
```

| 関数 | 説明 |
|------|------|
| `Dump(vars ...any) string` | 値を `$var1, $var2, ...` 形式の文字列にして返す |
| `DumpNamed(name string, v any) string` | 変数名付きでダンプ文字列を返す |
| `Dd(vars ...any)` | ダンプ結果を標準出力に書き出す |
| `DdNamed(name string, v any)` | 変数名付きで標準出力に書き出す |

struct・map・slice・pointer・循環参照・nil など主要な型に対応。非エクスポートフィールドは `<型名>` として表示。
