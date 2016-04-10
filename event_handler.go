package linebot

const (
	HandlerEventTypeTextMessage     = "text_message"
	HandlerEventTypeImageMessage    = "image_message"
	HandlerEventTypeVideoMessage    = "video_message"
	HandlerEventTypeAudioMessage    = "audio_message"
	HandlerEventTypeLocationMessage = "location_message"
	HandlerEventTypeStickerMessage  = "sticker_message"
	HandlerEventTypeContactMessage  = "contact_message"

	HandlerEventTypeAddedAsFriendOperation  = "added_as_friend_operation"
	HandlerEventTypeBlockedAccountOperation = "blocked_account"
)

type (
	EventHandler interface {
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

	messageEvent struct {
		typ     string
		event   *Event
		message *Message
	}

	operationEvent struct {
		typ       string
		event     *Event
		operation *Operation
	}

	AsyncEventDispatcher struct {
		handler        EventHandler
		messageQueue   chan *messageEvent
		operationQueue chan *operationEvent
	}
)

func NewAsyncEventDispatcher(handler EventHandler, queueSize int) *AsyncEventDispatcher {
	return &AsyncEventDispatcher{
		handler:        handler,
		messageQueue:   make(chan *messageEvent, queueSize),
		operationQueue: make(chan *operationEvent, queueSize),
	}
}

func (d *AsyncEventDispatcher) Run() {
	go func() {
		for {
			select {
			case msg, ok := <-d.messageQueue:
				if ok {
					d.dispatchMessage(msg)
				}
			case op, ok := <-d.operationQueue:
				if ok {
					d.dispatchOperation(op)
				}
			}
		}
	}()
}

func (d *AsyncEventDispatcher) dispatchMessage(msg *messageEvent) {
	switch msg.typ {
	case HandlerEventTypeTextMessage:
		d.handler.OnTextMessage(msg.event, msg.message)
	case HandlerEventTypeImageMessage:
		d.handler.OnImageMessage(msg.event, msg.message)
	case HandlerEventTypeVideoMessage:
		d.handler.OnVideoMessage(msg.event, msg.message)
	case HandlerEventTypeAudioMessage:
		d.handler.OnAudioMessage(msg.event, msg.message)
	case HandlerEventTypeLocationMessage:
		d.handler.OnLocationMessage(msg.event, msg.message)
	case HandlerEventTypeStickerMessage:
		d.handler.OnStickerMessage(msg.event, msg.message)
	case HandlerEventTypeContactMessage:
		d.handler.OnContactMessage(msg.event, msg.message)
	}
}

func (d *AsyncEventDispatcher) dispatchOperation(msg *operationEvent) {
	switch msg.typ {
	case HandlerEventTypeAddedAsFriendOperation:
		d.handler.OnAddedAsFriendOperation(msg.event, msg.operation)
	case HandlerEventTypeBlockedAccountOperation:
		d.handler.OnBlockedAccountOperation(msg.event, msg.operation)
	}
}

func newMessageEvent(typ string, ev *Event, msg *Message) *messageEvent {
	return &messageEvent{
		typ:     typ,
		event:   ev,
		message: msg,
	}
}

func newOperationEvent(typ string, ev *Event, op *Operation) *operationEvent {
	return &operationEvent{
		typ:       typ,
		event:     ev,
		operation: op,
	}
}

func (d *AsyncEventDispatcher) OnAddedAsFriendOperation(ev *Event, op *Operation) {
	d.operationQueue <- newOperationEvent(HandlerEventTypeAddedAsFriendOperation, ev, op)
}

func (d *AsyncEventDispatcher) OnBlockedAccountOperation(ev *Event, op *Operation) {
	d.operationQueue <- newOperationEvent(HandlerEventTypeBlockedAccountOperation, ev, op)
}

func (d *AsyncEventDispatcher) OnTextMessage(ev *Event, msg *Message) {
	d.messageQueue <- newMessageEvent(HandlerEventTypeTextMessage, ev, msg)
}

func (d *AsyncEventDispatcher) OnImageMessage(ev *Event, msg *Message) {
	d.messageQueue <- newMessageEvent(HandlerEventTypeImageMessage, ev, msg)
}

func (d *AsyncEventDispatcher) OnVideoMessage(ev *Event, msg *Message) {
	d.messageQueue <- newMessageEvent(HandlerEventTypeVideoMessage, ev, msg)
}

func (d *AsyncEventDispatcher) OnAudioMessage(ev *Event, msg *Message) {
	d.messageQueue <- newMessageEvent(HandlerEventTypeAudioMessage, ev, msg)
}

func (d *AsyncEventDispatcher) OnLocationMessage(ev *Event, msg *Message) {
	d.messageQueue <- newMessageEvent(HandlerEventTypeLocationMessage, ev, msg)
}

func (d *AsyncEventDispatcher) OnStickerMessage(ev *Event, msg *Message) {
	d.messageQueue <- newMessageEvent(HandlerEventTypeStickerMessage, ev, msg)
}

func (d *AsyncEventDispatcher) OnContactMessage(ev *Event, msg *Message) {
	d.messageQueue <- newMessageEvent(HandlerEventTypeContactMessage, ev, msg)
}
