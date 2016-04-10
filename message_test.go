package linebot

import (
	"encoding/json"
	"testing"
)

func TestParseRequest2(t *testing.T) {
	reqStr := `{
"result":[{
	"content":{
		"toType":1,
		"createdTime":1460182428120,
		"from":"u88f76bdbedd1beba7367e5b21343b410",
		"location":null,
		"id":"4147327103651",
		"to":["ufe9e45c6bbcfd2fa7d53a1eb26ae76dc"],
		"text":"あああ",
		"contentMetadata":{"AT_RECV_MODE":"2"},
		"deliveredTime":0,
		"contentType":1,
		"seq":null 
	},
	"createdTime":1460182428138,
	"eventType":"138311609000106303",
	"from":"u206d25c2ea6bd87c17655609a1c37cb8",
	"fromChannel":1341301815,
	"id":"WB1519-3344403275",
	"to":["ufe9e45c6bbcfd2fa7d53a1eb26ae76dc"],
	"toChannel":1461568637}]}`

	var lrj RequestJSON
	err := unmarshalJSON([]byte(reqStr), &lrj)
	if err != nil {
		t.Errorf("failed to unmarshal request %s", err)
	}

	events := lrj.Events
	if len(events) != 1 {
		t.Errorf("invalid count of events")
	}

	ev := events[0]
	msgJSON, err := json.Marshal(ev.Content)
	if err != nil {
		t.Errorf("failed to re-marshal content %s", err)
	}
	var msg Message
	err = unmarshalJSON(msgJSON, &msg)
	if err != nil {
		t.Errorf("failed to parse content %s %v", err, ev.Content)
	}
}

func TestParseRequest(t *testing.T) {
	reqStr := `{
"result":[{
	"from":"u206d25c2ea6bd87c17655609a1c37cb8",
	"fromChannel":1341301815,
	"to":["u0cc15697597f61dd8b01cea8b027050e"],
	"toChannel":1441301333,
	"eventType":"138311609000106303",
	"id":"ABCDEF-12345678901",
	"content":{
		"params":[
			"u0f3bfc598b061eba02183bfc5280886a",
			null,
			null
		],
		"revision":2469,
		"opType":4
	}
}]
}`
	var lrj RequestJSON
	err := json.Unmarshal([]byte(reqStr), &lrj)
	if err != nil {
		t.Errorf("failed to unmarshal request %s", err)
	}

	events := lrj.Events
	if len(events) != 1 {
		t.Errorf("invalid count of events")
	}

	ev := events[0]

	if ev.From != "u206d25c2ea6bd87c17655609a1c37cb8" {
		t.Errorf("invalid From value")
	}

	if ev.Content == nil {
		t.Errorf("content not found")
		return
	}

	opJSON, err := json.Marshal(ev.Content)
	if err != nil {
		t.Errorf("failed rebuild content JSON")
		return
	}

	var op Operation
	err = json.Unmarshal(opJSON, &op)
	if err != nil {
		t.Errorf("failed to parse operation JSON")
		return
	}

	if op.Revision != 2469 {
		t.Errorf("invalid op revision")
	}

}
