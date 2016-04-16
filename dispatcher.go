package linebot

import (
	"encoding/json"
	"fmt"
)

type eventDispatcher struct {
	handler    EventHandler
	eventQueue chan *Event
}

func newEventDispatcher(handler EventHandler, queueSize int) *eventDispatcher {
	return &eventDispatcher{
		handler:    handler,
		eventQueue: make(chan *Event, queueSize),
	}
}

func (d *eventDispatcher) run() {
	for {
		select {
		case e, ok := <-d.eventQueue:
			if ok {
				d.dispatchEvent(e)
			}
		}
	}
}

func (d *eventDispatcher) handleEvent(e *Event) {
	d.eventQueue <- e
}

func (d *eventDispatcher) dispatchEvent(e *Event) {
	switch e.EventType {
	case EventTypeReceivedMessage:
		d.handleMessage(e)
	case EventTypeReceivedOperation:
		d.handleOperation(e)
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

func (d *eventDispatcher) handleMessage(e *Event) {
	msg, err := convertToMessage(e.Content)
	if err != nil {
		// TODO logging
		return
	}
	d.dispatchMessage(msg)
}

func (d *eventDispatcher) dispatchMessage(msg *Message) {
	switch msg.ContentType {
	case ContentTypeText:
		d.handler.OnTextMessage(msg.From, msg.Text)
	case ContentTypeImage:
		d.handler.OnImageMessage(msg.From)
	case ContentTypeVideo:
		d.handler.OnVideoMessage(msg.From)
	case ContentTypeAudio:
		d.handler.OnAudioMessage(msg.From)
	case ContentTypeLocation:
		d.handler.OnLocationMessage(msg.From, msg.Location.Title, msg.Location.Address,
			msg.Location.Latitude, msg.Location.Longitude)
	case ContentTypeSticker:
		d.handler.OnStickerMessage(msg.From, msg.ContentMetadata.STKPKGID,
			msg.ContentMetadata.STKID, msg.ContentMetadata.STKVER, msg.ContentMetadata.STKTXT)
	case ContentTypeContact:
		d.handler.OnContactMessage(msg.From, msg.ContentMetadata.MID,
			msg.ContentMetadata.DisplayName)
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

func (d *eventDispatcher) handleOperation(e *Event) {
	op, err := convertToOperation(e.Content)
	if err != nil {
		// TODO logging
		return
	}
	d.dispatchOperation(op)
}

func (d *eventDispatcher) dispatchOperation(op *Operation) {
	switch op.OPType {
	case OPTypeAddedAsFriend:
		d.handler.OnAddedAsFriendOperation(op.Params)
	case OPTypeBlockedAccount:
		d.handler.OnBlockedAccountOperation(op.Params)
	}
}
