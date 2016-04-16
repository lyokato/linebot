# Golang LINE Bot Application Framework

## 注意

THIS IS NOT STABLE VERSION

まだAPIガンガン変えていくと思います

各種別のメッセージがまだテストしきれてなくてrobustじゃないです

## Building Bot Server

### Step 1: イベントを受け取る

まずはlinebot.EventHandlerのインターフェースを実装したクラスを用意しましょう

インターフェースの仕様は次のようになっています。
メソッド名を見れば、それがどんなイベントを処理するものか想像できると思います。

```
type EventHandler interface {
	OnAddedAsFriendOperation(MIDs []string)
	OnBlockedAccountOperation(MIDs []string)
	OnTextMessage(from, text string)
	OnImageMessage(from string)
	OnVideoMessage(from string)
	OnAudioMessage(from string)
	OnLocationMessage(from, title, address string, latitude, longitude float64)
	OnStickerMessage(from, stickerPackageId, stickerId, stickerVersion, stickerText string)
	OnContactMessage(from, MID, displayName string)
}
```

このインターフェースを実装し、イベントを受け取ったらログを出力するだけのハンドラを用意してみましょう
```
package your_event_handler

import (
	"github.com/lyokato/linebot"
	log "github.com/Sirupsen/logrus"
)

type YourEventHandler struct {
}

func New() *YourEventHandler {
	return &YourEventHandler{}
}

func (h *YourEventHandler) OnAddedAsFriendOperation(MIDs []string) {
	log.Infof("OnAddedAsFriendOperation: %v %v", MIDs)
}

func (h *YourEventHandler) OnBlockedAccountOperation(MIDs []string) {
	log.Infof("OnBlockedAccountOperation: %v", MIDs)
}

func (h *YourEventHandler) OnTextMessage(from, text string) {
	log.Infof("OnTextMesssage: %s %s", from, text)
}

func (h *YourEventHandler) OnImageMessage(from string) {
	log.Infof("OnImageMesssage: %s", from)
}

func (h *YourEventHandler) OnVideoMessage(from string) {
	log.Infof("OnVideoMesssage: %s", from)
}

func (h *YourEventHandler) OnAudioMessage(from string) {
	log.Infof("OnAudioMesssage: %s", from)
}

func (h *YourEventHandler) OnLocationMessage(from, title, address string, latitude, longitude float64)
	log.Infof("OnLocationMesssage: %s %s %s", from, title, address)
}

func (h *YourEventHandler) OnStickerMessage(from, stickerPackageId, stickerId, stickerVersion, stickerText string)
	log.Infof("OnStickerMesssage: %s %s %s", from, stickerPackageId, stickerId)
}

func (h *YourEventHandler) OnContactMessage(from, MID, displayName string)
	log.Infof("OnContactMesssage: %s %s %s", from, MID, displayName)
}
```

以下のようにlinebot.NewServerを使ってbotServerを作成し。
そしてHTTPHandlerに以下の三つを渡し、http.HandleFuncを準備します。

- channelSecret: あらかじめLINEのDeveloper Channelで発行された*channel secret*の文字列を指定
- eventHandler: 上で準備したEventHandlerインターフェースを実装したオブジェクト
- queueSize: eventを扱うchannelのサイズ

```
package main

import (
  "github.com/lyokato/linebot"
  "your_event_handler"
)

const cnannelSecret = "your channel secret"

func main() {
  ...
  evh := your_event_handler.New()
  botServer := linebot.NewServer()
  http.HandleFunc("/callback", botServer.HTTPHandler(channelSecret, evh, eventQueueSize)) 
  http.ListenAndServe(address, nil)
  ...
}
```

お使いに環境に合わせて
Routingの設定をするとよいでしょう。

GolangのWeb Application Frameworkであれば、多くの場合、
http.HandleFuncをその環境用のハンドラにラップする手段が用意されています。

以下はGinの例です。gin.WrapFを使ってhttp.HandleFuncを変換しています。

```
g := gin.Default()
g.POST("/", gin.WrapF(botServer.HTTPHandler(botServer.HTTPHandler(channelSecret, evh, queueSize))))
```

### Step 2: メッセージを送信する

これまでは受信だけの説明をしてきました。
ユーザーにメッセージを送らなければ対話は成り立ちません。
そのためには、Clientを利用します


LINE側であらかじめ取得した*channel id*, *channel secret*, *channel MID*を引数で指定し、
クライアントを作成します。


あとは以下のように*PostText*、あるいはその他のClientのメソッドで、
簡単にコンテンツを送信することが可能です。
```
c := linebot.NewClient(channelId, channelSecret, mid)
c.PostText(to, "だめだハゲ")
```

Clientのインターフェース定義は以下のようになっています

```
type Client interface {
		PostEvent(r *PostEventRequest)
		PostText(to, text string)
		PostImage(to, imageUrl, thumbnailUrl string)
		PostVideo(to, movieUrl, thumbnailUrl string)
		PostAudio(to, audioUrl string, playTimeMilliSeconds int)
		PostLocation(to, locationTitle, address string, latitude, longitude float64)
		PostSticker(to, stickerId, stickerPackageId, stickerVersion string)
	}
```

### Step 3: ユーザーからのメッセージに対して返信する

Step1で作成したYourEventHandlerの一部を次のように変更します

```
package your_event_handler

import (
	"github.com/lyokato/linebot"
	log "github.com/Sirupsen/logrus"
)

type YourEventHandler struct {
  client linebot.Client
}

func New(client linebot.Client) *YourEventHandler {
	return &YourEventHandler{
    client: client, 
  }
}

func (h *YourEventHandler) OnTextMessage(from, text string) {
	log.Infof("OnTextMesssage: %s %s", from, text)

  h.client.PostText(from, text)
}

...
```

これにより、送信者に発言をそのまま返すエコーサーバーが実現できます
メイン関数も次のように、変えておきましょう。
```
func main() {
...
  c := linebot.NewClient(channelId, channelSecre, mid)

  evh := your_event_handler.New(c)
  botServer := linebot.NewServer()

  g := gin.Default()
  g.POST("/", gin.WrapF(botServer.HTTPHandler(channelSecret, evh, queueSize)))
...
}

```

## Example

exampleディレクトリ以下にサンプルアプリケーションを用意しています。
ボイラープレートとしても使えます。

```
cd example
```

依存を解決しておきます
```
glide up
```

以下のハンドラを編集します
```
handler/handler.go
```

以下のコマンドでアプリケーションをビルドできます
```
./support/build/bot_server.sh
```

これでLinux環境用にビルドして、./support/deploy/roles/deploy-app/files/以下に
実行バイナリがコピーされます。以下で説明するAnsibleでのデプロイ時にそのまま利用されます。

### Deploy

Example環境で、Ansibleを利用して簡単にデプロイできるようにplaybookが準備されています
接続先のOSはCentOSを想定しています (systemdを前提としているため)

```
cd ./support/deploy
```

ここでhostsファイルを編集しておきましょう

```
vim hosts
```
自分の環境に合わせてIPアドレスを指定しておきましょう
```
[bot_server]
XXX.XXX.XXX.XXX
```

次にPlaybookを編集します
```
vim bot_server.yml
```

```
- hosts: bot_server
  become: yes
  roles: 
    - { role: add-user, name: replaceme, group: replaceme, shell: /bin/bash, pass: replaceme, pubkey: ~/.ssh/id_rsa.pub }
    - { role: deploy-app, app_name: bot_server, bin_name: bot_server, config_file: bot_server.toml, log_level: debug }
```

接続先のLinux(CentOS)でユーザーを追加するようになっています。
ユーザー名などを変更するか、必要なければこのrole(add-user)を削除してください

Configファイルを編集しておきます
```
vim roles/deploy-app/files/bot_server.toml
```
ChannelID, ChannelSecret, MIDなどを、LINEの管理コンソール上で取得したものに置き換えましょう
```
[web]

host = "0.0.0.0"
port = 80

[bot]

channel_id = 100000000
channel_secret = "your_channel_secret"
channel_mid = "your_channel_mid"
client_worker_queue_size = 10
event_dispatcher_queue_size = 10
```

次のコマンドでデプロイできますが、環境変数EC2_KEYに
あなたのEC2の鍵のパスを指定しておいてください
```
./support/deploy/bot_server.sh
```

うまくいけばsystemdにサービス登録されてサーバ上で再起動されます

## TODO

Clientのバッチリクエストなど

## AUTHOR

lyo.kato __at__ gmail.com

## LICENSE 

MIT
