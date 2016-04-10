package linebot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const HTTPHeaderNameLineChannelSignature = "X-LINE-ChannelSignature"

type Server struct {
	channelSecret string
	eventHandler  EventHandler
	logger        Logger
}

func NewServer(channelSecret string, eventHandler EventHandler) *Server {
	return &Server{
		channelSecret: channelSecret,
		eventHandler:  eventHandler,
		logger:        &defaultLogger{},
	}
}

func (b *Server) SetLogger(l Logger) {
	b.logger = l
}

func (b *Server) HTTPHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		b.logger.Debug("start to check signature")

		body, _ := ioutil.ReadAll(r.Body)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		sig := r.Header.Get(HTTPHeaderNameLineChannelSignature)
		if !verify(sig, body, b.channelSecret) {
			b.logger.Debugf("invalid signature: %s", string(sig))
			w.WriteHeader(400)
			return
		}
		var rj RequestJSON
		err := unmarshalJSON(body, &rj)
		if err != nil {
			b.logger.Debugf("failed to parse JSON: %s %v", err, body)
			w.WriteHeader(400)
			return
		}
		for _, ev := range rj.Events {
			b.logger.Debugf("try to handle event")
			b.handleEvent(ev)
		}
		w.WriteHeader(200)
		return
	}
}

func (b *Server) handleEvent(e *Event) error {
	switch e.EventType {
	case EventTypeReceivedMessage:
		b.logger.Debugf("try to handle message")
		return b.handleMessage(e)
	case EventTypeReceivedOperation:
		b.logger.Debugf("try to handle operation")
		return b.handleOperation(e)
	default:
		return fmt.Errorf("Unknown event type: %s", e.EventType)
	}
}

func convertToMessage(content map[string]interface{}) (*Message, error) {
	msgJSON, err := json.Marshal(content)
	if err != nil {
		return nil, fmt.Errorf("failed rebuild message JSON: %s", err)
	}
	var msg Message
	err = unmarshalJSON(msgJSON, &msg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse message JSON: %s", err)
	}
	return &msg, nil
}

func (b *Server) handleMessage(e *Event) error {
	msg, err := convertToMessage(e.Content)
	if err != nil {
		b.logger.Debugf("failed to convert message: %s", err)
		return err
	}
	switch msg.ContentType {
	case ContentTypeText:
		b.eventHandler.OnTextMessage(e, msg)
		return nil
	case ContentTypeImage:
		b.eventHandler.OnImageMessage(e, msg)
		return nil
	case ContentTypeVideo:
		b.eventHandler.OnVideoMessage(e, msg)
		return nil
	case ContentTypeAudio:
		b.eventHandler.OnAudioMessage(e, msg)
		return nil
	case ContentTypeLocation:
		b.eventHandler.OnLocationMessage(e, msg)
		return nil
	case ContentTypeSticker:
		b.eventHandler.OnStickerMessage(e, msg)
		return nil
	case ContentTypeContact:
		b.eventHandler.OnContactMessage(e, msg)
		return nil
	default:
		return fmt.Errorf("unknown content type: %s", msg.ContentType)
	}
}

func convertToOperation(content map[string]interface{}) (*Operation, error) {
	opJSON, err := json.Marshal(content)
	if err != nil {
		return nil, fmt.Errorf("failed rebuild content JSON: %s", err)
	}
	var op Operation
	err = unmarshalJSON(opJSON, &op)
	if err != nil {
		return nil, fmt.Errorf("failed to parse operation JSON: %s", err)
	}
	return &op, nil
}

func (b *Server) handleOperation(e *Event) error {
	op, err := convertToOperation(e.Content)
	if err != nil {
		b.logger.Debugf("failed to convert operation: %s", err)
		return err
	}
	switch op.OPType {
	case OPTypeAddedAsFriend:
		b.eventHandler.OnAddedAsFriendOperation(e, op)
		return nil
	case OPTypeBlockedAccount:
		b.eventHandler.OnBlockedAccountOperation(e, op)
		return nil
	default:
		return fmt.Errorf("unknown operation type: %s", op.OPType)
	}
}
