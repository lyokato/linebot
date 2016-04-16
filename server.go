package linebot

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

const HTTPHeaderNameLineChannelSignature = "X-LINE-ChannelSignature"

type Server struct {
	logger Logger
}

func NewServer() *Server {
	return &Server{
		logger: &defaultLogger{},
	}
}

func (b *Server) SetLogger(l Logger) {
	b.logger = l
}

func (b *Server) HTTPHandler(channelSecret string, eventHandler EventHandler, queueSize int) http.HandlerFunc {

	d := newEventDispatcher(eventHandler, queueSize)
	go d.run()

	return func(w http.ResponseWriter, r *http.Request) {

		b.logger.Debug("start to check signature")

		body, _ := ioutil.ReadAll(r.Body)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		sig := r.Header.Get(HTTPHeaderNameLineChannelSignature)
		if !verify(sig, body, channelSecret) {
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
			d.handleEvent(ev)
		}
		w.WriteHeader(200)
		return
	}

}
