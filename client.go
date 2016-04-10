package linebot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	HTTPHeaderNameLineChannelID          = "X-Line-ChannelID"
	HTTPHeaderNameLineChannelSecret      = "X-Line-ChannelSecret"
	HTTPHeaderNameLineTrustedUserWithACL = "X-Line-Trusted-User-With-ACL"
	TrialEndpointHost                    = "trialbot-api.line.me"
	PostEventEndpointPath                = "v1/events"
)

type (
	Client interface {
		PostEvent(r *PostEventRequest)
		PostText(to, text string)
		PostImage(to, imageUrl, thumbnailUrl string)
		PostVideo(to, movieUrl, thumbnailUrl string)
		PostAudio(to, audioUrl string, playTimeMilliSeconds int)
		PostLocation(to, text, locationTitle string, latitude, longitude float64)
		PostSticker(to, stickerId, stickerPackageId, stickerVersion string)
	}

	client struct {
		channelId     int
		channelSecret string
		mid           string
		logger        Logger
	}

	PostEventRequest struct {
		To        []string    `json:"to"`
		ToChannel int         `json:"toChannel"`
		EventType string      `json:"eventType"`
		Content   interface{} `json:"content"`
	}
)

func postEventURL() string {
	return fmt.Sprintf("https://%s/%s", TrialEndpointHost, PostEventEndpointPath)
}

func NewClient(channelId int, channelSecret, mid string) Client {
	return &client{
		channelId:     channelId,
		channelSecret: channelSecret,
		mid:           mid,
		logger:        &defaultLogger{},
	}
}

func (c *client) PostEvent(r *PostEventRequest) {
	s, err := json.Marshal(r)
	if err != nil {
		c.logger.Infof("failed to marshal json: %s", err)
		return
	}
	url := postEventURL()
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(s))
	if err != nil {
		c.logger.Infof("failed to build post request: %s", err)
		return
	}

	cidStr := strconv.Itoa(c.channelId)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set(HTTPHeaderNameLineChannelID, cidStr)
	req.Header.Set(HTTPHeaderNameLineChannelSecret, c.channelSecret)
	req.Header.Set(HTTPHeaderNameLineTrustedUserWithACL, c.mid)

	client := &http.Client{}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		c.logger.Infof("http request failed : %s", err)
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		c.logger.Infof("http response status code is not OK : %d %s", resp.StatusCode, string(body))
		return
	}
}

func buildPostEventRequest(to string, content *OutboundContent) *PostEventRequest {
	return &PostEventRequest{
		To:        []string{to},
		ToChannel: ToChannelForSendingMessage,
		EventType: EventTypeSendingMessage,
		Content:   content,
	}
}

func buildTextMessageRequest(to, text string) *PostEventRequest {
	content := &OutboundContent{
		ContentType: ContentTypeText,
		ToType:      ToTypeUser,
		Text:        text,
	}
	return buildPostEventRequest(to, content)
}

func buildImageMessageRequest(to, imageUrl, thumbnailUrl string) *PostEventRequest {
	content := &OutboundContent{
		ContentType:        ContentTypeImage,
		ToType:             ToTypeUser,
		OriginalContentUrl: imageUrl,
		PreviewImageUrl:    thumbnailUrl,
	}
	return buildPostEventRequest(to, content)
}

func buildVideoMessageRequest(to, movieUrl, thumbnailUrl string) *PostEventRequest {
	content := &OutboundContent{
		ContentType:        ContentTypeVideo,
		ToType:             ToTypeUser,
		OriginalContentUrl: movieUrl,
		PreviewImageUrl:    thumbnailUrl,
	}
	return buildPostEventRequest(to, content)
}

func buildAudioMessageRequest(to, audioUrl string, playTimeMilliSeconds int) *PostEventRequest {
	content := &OutboundContent{
		ContentType:        ContentTypeAudio,
		ToType:             ToTypeUser,
		OriginalContentUrl: audioUrl,
		Metadata: &ContentMetadata{
			AUDLEN: strconv.Itoa(playTimeMilliSeconds),
		},
	}
	return buildPostEventRequest(to, content)
}

func buildLocationMessageRequest(to, text, locationTitle string, latitude, longitude float64) *PostEventRequest {
	content := &OutboundContent{
		ContentType: ContentTypeLocation,
		ToType:      ToTypeUser,
		Text:        text,
		Location: &Location{
			Title:     locationTitle,
			Latitude:  latitude,
			Longitude: longitude,
		},
	}
	return buildPostEventRequest(to, content)
}

func buildStickerMessageRequest(to, stickerId, stickerPackageId, stickerVersion string) *PostEventRequest {
	content := &OutboundContent{
		ContentType: ContentTypeSticker,
		ToType:      ToTypeUser,
		Metadata: &ContentMetadata{
			STKID:    stickerId,
			STKPKGID: stickerPackageId,
			STKVER:   stickerVersion,
		},
	}
	return buildPostEventRequest(to, content)
}

func (c *client) PostText(to, text string) {
	r := buildTextMessageRequest(to, text)
	c.PostEvent(r)
}

func (c *client) PostImage(to, imageUrl, thumbnailUrl string) {
	r := buildImageMessageRequest(to, imageUrl, thumbnailUrl)
	c.PostEvent(r)
}

func (c *client) PostVideo(to, movieUrl, thumbnailUrl string) {
	r := buildVideoMessageRequest(to, movieUrl, thumbnailUrl)
	c.PostEvent(r)
}

func (c *client) PostAudio(to, audioUrl string, playTimeMilliSeconds int) {
	r := buildAudioMessageRequest(to, audioUrl, playTimeMilliSeconds)
	c.PostEvent(r)
}

func (c *client) PostLocation(to, text, locationTitle string, latitude, longitude float64) {
	r := buildLocationMessageRequest(to, text, locationTitle, latitude, longitude)
	c.PostEvent(r)
}

func (c *client) PostSticker(to, stickerId, stickerPackageId, stickerVersion string) {
	r := buildStickerMessageRequest(to, stickerId, stickerPackageId, stickerVersion)
	c.PostEvent(r)
}
