package linebot

type ClientWorker struct {
	requestQueue chan *PostEventRequest
	client       Client
}

func NewClientWorker(channelId int, channelSecret string, mid string, queueSize int) *ClientWorker {
	return &ClientWorker{
		requestQueue: make(chan *PostEventRequest, queueSize),
		client:       NewClient(channelId, channelSecret, mid),
	}
}

func (w *ClientWorker) PostText(to, text string) {
	r := buildTextMessageRequest(to, text)
	w.PostEvent(r)
}

func (w *ClientWorker) PostImage(to, imageUrl, thumbnailUrl string) {
	r := buildImageMessageRequest(to, imageUrl, thumbnailUrl)
	w.PostEvent(r)
}

func (w *ClientWorker) PostVideo(to, movieUrl, thumbnailUrl string) {
	r := buildVideoMessageRequest(to, movieUrl, thumbnailUrl)
	w.PostEvent(r)
}

func (w *ClientWorker) PostAudio(to, audioUrl string, playTimeMilliSeconds int) {
	r := buildAudioMessageRequest(to, audioUrl, playTimeMilliSeconds)
	w.PostEvent(r)
}

func (w *ClientWorker) PostLocation(to, locationTitle, address string, latitude, longitude float64) {
	r := buildLocationMessageRequest(to, locationTitle, address, latitude, longitude)
	w.PostEvent(r)
}

func (w *ClientWorker) PostSticker(to, stickerId, stickerPackageId, stickerVersion string) {
	r := buildStickerMessageRequest(to, stickerId, stickerPackageId, stickerVersion)
	w.PostEvent(r)
}

func (w *ClientWorker) PostEvent(msg *PostEventRequest) {
	w.requestQueue <- msg
}

func (w *ClientWorker) Run() {
	go func() {
		for {
			select {
			case msg, ok := <-w.requestQueue:
				if ok {
					w.client.PostEvent(msg)
				}
			}
		}
	}()
}
