# Golang LINE Bot Application Framework

## 注意

THIS IS NOT STABLE VERSION

まだAPIガンガン変えていくと思います

各種別のメッセージがまだテストしきれてなくてrobustじゃないです

## Buildgin Bot Server

### Step 1: イベントを受け取る

まずはlinebot.EventHandlerのインターフェースを実装したクラスを用意しましょう

インターフェースの仕様は次のようになっています。
メソッド名を見れば、それがどんなイベントを処理するものか想像できると思います。

```
type EventHandler interface {
		OnAddedAsFriendOperation(e *Event, op *Operation)
		OnBlockedAccountOperation(e *Event, op *Operation)
		OnTextMessage(e *Event, msg *Message)
		OnImageMessage(e *Event, msg *Message)
		OnVideoMessage(e *Event, msg *Message)
		OnAudioMessage(e *Event, msg *Message)
		OnLocationMessage(e *Event, msg *Message)
		OnStickerMessage(e *Event, msg *Message)
		OnContactMessage(e *Event, msg *Message)
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

func (h *YourEventHandler) OnAddedAsFriendOperation(ev *linebot.Event, op *linebot.Operation) {
	log.Infof("OnAddedAsFriendOperation: %v %v", ev, op)
}

func (h *YourEventHandler) OnBlockedAccountOperation(ev *linebot.Event, op *linebot.Operation) {
	log.Infof("OnBlockedAccountOperation: %v %v", ev, op)
}

func (h *YourEventHandler) OnTextMessage(ev *linebot.Event, msg *linebot.Message) {
	log.Infof("OnTextMesssage: %v %v", ev, msg)
}

func (h *YourEventHandler) OnImageMessage(ev *linebot.Event, msg *linebot.Message) {
	log.Infof("OnImageMesssage: %v %v", ev, msg)
}

func (h *YourEventHandler) OnVideoMessage(ev *linebot.Event, msg *linebot.Message) {
	log.Infof("OnVideoMesssage: %v %v", ev, msg)
}

func (h *YourEventHandler) OnAudioMessage(ev *linebot.Event, msg *linebot.Message) {
	log.Infof("OnAudioMesssage: %v %v", ev, msg)
}

func (h *YourEventHandler) OnLocationMessage(ev *linebot.Event, msg *linebot.Message) {
	log.Infof("OnLocationMesssage: %v %v", ev, msg)
}

func (h *YourEventHandler) OnStickerMessage(ev *linebot.Event, msg *linebot.Message) {
	log.Infof("OnStickerMesssage: %v %v", ev, msg)
}

func (h *YourEventHandler) OnContactMessage(ev *linebot.Event, msg *linebot.Message) {
	log.Infof("OnContactMesssage: %v %v", ev, msg)
}
```

以下のようにlinebot.NewServerメソッドに二つの引数を渡してbotServerの準備をします。

- channelSecret: あらかじめLINEのDeveloper Channelで発行された*channel secret*の文字列を指定
- yourEventHandler: 上で準備したEventHandlerインターフェースを実装したオブジェクト

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
  botServer := linebot.NewServer(channelSecret, evh)
  http.HandleFunc("/callback", botServer.HTTPHandler()) 
  http.ListenAndServe(address, nil)
  ...
}
```

上の例のように、botServer.HTTPHandler()で、http.HandleFuncを返します。
お使いに環境に合わせて
Routingの設定をするとよいでしょう。

GolangのWeb Application Frameworkであれば、多くの場合、
http.HandleFuncをその環境用のハンドラにラップする手段が用意されています。

以下はGinの例です。gin.WrapFを使ってhttp.HandleFuncを変換しています。

```
g := gin.Default()
g.POST("/", gin.WrapF(botServer.HTTPHandler()))
```

### Step 2: 非同期でイベントを受け取る 


上の例ではログ出力するだけなので問題ないのですが、
今後、機能を追加していく際に、内部のデータベースやマイクロサービス、
あるいは外部のサービスとの連動などが必要になってくるかもしれません。

そうすると一つのイベントハンドラ内での処理時間が長くなってくることが予想されます。
LINEからリクエストが来たら、10秒以内にレスポンスは返さなければなりませんし、
そもそもHTTPのリクエスト処理のプールを出来るだけ塞がないようにするほうがよいでしょう。

このために、AsyncEventDispatcherが用意されています。
EventHandlerを直接使わずに、AsyncEventDispatcherでラップします。
```
evh := your_event_handler.New()

queueSize = 10
evd := linebot.NewAsyncEventDispatcher(evh, queueSize)
evd.Run()

botServer := linebot.NewServer(channelSecret, evd)

g := gin.Default()
g.POST("/", gin.WrapF(botServer.HTTPHandler()))

```
AsyncEventDispatcherのRunメソッドを呼ぶと、独立したgoroutineが開始され、
その中のループでイベントハンドラが呼び出されます。

このためLineからのリクエストを受け取るHTTPサーバーの処理をつまらせることがありません。

(あくまで簡易的なものなので、本格的に非同期の分散処理が必要ならば、AmazonSQSやNSQなどの
メッセージキュープロダクトを利用するのがよいでしょう)

### Step 3: メッセージを送信する

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

### Step 4: ユーザーからのメッセージに対して返信する

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

func (h *YourEventHandler) OnTextMessage(ev *linebot.Event, msg *linebot.Message) {
	log.Infof("OnTextMesssage: %v %v", ev, msg)

  h.client.PostText(msg.From, msg.Text)
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

  evd := linebot.NewAsyncEventDispatcher(evh, queueSize)
  evd.Run()

  botServer := linebot.NewServer(channelSecret, evd)

  g := gin.Default()
  g.POST("/", gin.WrapF(botServer.HTTPHandler()))
...
}

```

### Step 5: クライアントも非同期処理に

クライアントの処理も、LINEへのHTTPリクエストで時間がかかることもあるかもしれません。
ここも独立したgoroutine内のループで処理することが可能です。

このためには、以下のようにClientWorkerを利用します。
こちらもAsyncEventDispatcherと同様に、queue sizeの指定と、Runメソッドでgoroutineを開始しておくのを忘れないようにしましょう。

```
func main() {
...
  queueSize = 10
  c := linebot.NewClientWorker(channelId, channelSecre, mid, queueSize)
  c.Run()

  evh := your_event_handler.New(c)

  evd := linebot.NewAsyncEventDispatcher(evh, queueSize)
  evd.Run()

  botServer := linebot.NewServer(channelSecret, evd)

  g := gin.Default()
  g.POST("/", gin.WrapF(botServer.HTTPHandler()))
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
