package handler

import (
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

func (h *ExampleEventHandler) OnAddedAsFriendOperation(ev *linebot.Event, op *linebot.Operation) {
	log.Infof("OnAddedAsFriendOperation: %v %v", ev, op)
}

func (h *ExampleEventHandler) OnBlockedAccountOperation(ev *linebot.Event, op *linebot.Operation) {
	log.Infof("OnBlockedAccountOperation: %v %v", ev, op)
}

func (h *ExampleEventHandler) OnTextMessage(ev *linebot.Event, msg *linebot.Message) {
	log.Infof("OnTextMesssage: %v %v", ev, msg)

	h.worker.PostText(msg.From, "Nice Words!")
}

func (h *ExampleEventHandler) OnImageMessage(ev *linebot.Event, msg *linebot.Message) {
	log.Infof("OnImageMesssage: %v %v", ev, msg)

	h.worker.PostText(msg.From, "Nice Picture!")
}

func (h *ExampleEventHandler) OnVideoMessage(ev *linebot.Event, msg *linebot.Message) {
	log.Infof("OnVideoMesssage: %v %v", ev, msg)

	h.worker.PostText(msg.From, "Nice Movie")
}

func (h *ExampleEventHandler) OnAudioMessage(ev *linebot.Event, msg *linebot.Message) {
	log.Infof("OnAudioMesssage: %v %v", ev, msg)

	h.worker.PostText(msg.From, "Nice Audio!")
}

func (h *ExampleEventHandler) OnLocationMessage(ev *linebot.Event, msg *linebot.Message) {
	log.Infof("OnLocationMesssage: %v %v", ev, msg)

	h.worker.PostText(msg.From, "Nice Place!")
}

func (h *ExampleEventHandler) OnStickerMessage(ev *linebot.Event, msg *linebot.Message) {
	log.Infof("OnStickerMesssage: %v %v", ev, msg)

	h.worker.PostText(msg.From, "Nice Sticker!")
}

func (h *ExampleEventHandler) OnContactMessage(ev *linebot.Event, msg *linebot.Message) {
	log.Infof("OnContactMesssage: %v %v", ev, msg)

	h.worker.PostText(msg.From, "Nice Person!")
}
