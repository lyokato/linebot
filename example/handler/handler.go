package handler

import (
	"fmt"

	"github.com/lyokato/linebot"

	log "github.com/Sirupsen/logrus"
)

type ExampleEventHandler struct {
	worker *linebot.ClientWorker
}

func New(worker *linebot.ClientWorker) *ExampleEventHandler {
	return &ExampleEventHandler{
		worker: worker,
	}
}

func (h *ExampleEventHandler) OnAddedAsFriendOperation(mids []string) {
	log.Infof("OnAddedAsFriendOperation: %v", mids)
}

func (h *ExampleEventHandler) OnBlockedAccountOperation(mids []string) {
	log.Infof("OnBlockedAccountOperation: %v", mids)
}

func (h *ExampleEventHandler) OnTextMessage(from, text string) {
	log.Infof("OnTextMesssage: %s %s", from, text)

	h.worker.PostText(from, text)
}

func (h *ExampleEventHandler) OnImageMessage(from string) {
	log.Infof("OnImageMesssage: %s", from)

	h.worker.PostImage(from, "http://navi.harinezumi.org/wp-content/uploads/kohari13.jpg", "http://navi.harinezumi.org/wp-content/themes/wp_temp_harinavi_v2.0/images/index4.jpg")
}

func (h *ExampleEventHandler) OnVideoMessage(from string) {
	log.Infof("OnVideoMesssage: %s", from)

	h.worker.PostText(from, "Nice Movie")
}

func (h *ExampleEventHandler) OnAudioMessage(from string) {
	log.Infof("OnAudioMesssage: %s", from)

	h.worker.PostText(from, "Nice Audio!")
}

func (h *ExampleEventHandler) OnLocationMessage(from, title, address string, latitude, longitude float64) {
	log.Infof("OnLocationMesssage: %s %s %s", from, title, address)

	h.worker.PostLocation(from, title, address, latitude, longitude)
}

func (h *ExampleEventHandler) OnStickerMessage(from, stickerPackageId, stickerId, stickerVersion, stickerText string) {
	log.Infof("OnStickerMesssage: %s", from)

	h.worker.PostSticker(from, stickerId, stickerPackageId, stickerVersion)
}

func (h *ExampleEventHandler) OnContactMessage(from, MID, displayName string) {
	log.Infof("OnContactMesssage: %s", from)

	h.worker.PostText(from, fmt.Sprintf("%s:%s", MID, displayName))
}
