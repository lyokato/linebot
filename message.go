package linebot

const (
	EventTypeReceivedMessage   = "138311609000106303"
	EventTypeSendingMessage    = "138311608800106203"
	EventTypeReceivedOperation = "138311609100106403"

	ToChannelForSendingMessage = 1383378250

	ContentTypeText     = 1
	ContentTypeImage    = 2
	ContentTypeVideo    = 3
	ContentTypeAudio    = 4
	ContentTypeLocation = 7
	ContentTypeSticker  = 8
	ContentTypeContact  = 10

	OPTypeAddedAsFriend  = 4
	OPTypeBlockedAccount = 8

	ToTypeUser = 1
)

type (
	OutboundContent struct {
		ContentType        int              `json:"contentType"`
		ToType             int              `json:"toType"`
		Text               string           `json:"text,omitempty"`
		OriginalContentUrl string           `json:"originalContentUrl,omitempty"`
		PreviewImageUrl    string           `json:"previewImageUrl,omitempty"`
		Metadata           *ContentMetadata `json:"contentMetadata,omitempty"`
		Location           *Location        `json:"location,omitempty"`
	}

	RequestJSON struct {
		Events []*Event `json:"result"`
	}

	Event struct {
		From        string                 `json:"from,omitempty"`
		FromChannel int                    `json:"fromChannel,omitempty"`
		To          []string               `json:"to,omitempty"`
		ToChannel   int                    `json:"toChannel,omitempty"`
		EventType   string                 `json:"eventType,omitempty"`
		Id          string                 `json:"id,omitempty"`
		Content     map[string]interface{} `json:"content,omitempty"`
	}

	Message struct {
		Id              string           `json:"id,omitempty"`
		ContentType     int              `json:"contentType,omitempty"`
		From            string           `json:"from,omitempty"`
		CreatedTime     int64            `json:"createdTime,omitempty"`
		To              []string         `json:"to,omitempty"`
		ToType          int              `json:"toType,omitempty"`
		ContentMetadata *ContentMetadata `json:"contentMetadata,omitempty"`
		Text            string           `json:"text,omitempty"`
		Location        *Location        `json:"location,omitempty"`
	}

	Location struct {
		Title     string  `json:"title"`
		Address   string  `json:"address"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	ContentMetadata struct {
		STKPKGID string `json:"STKPKGID,omitempty"`
		STKID    string `json:"STKID,omitempty"`
		STKVER   string `json:"STKVER,omitempty"`
		STKTXT   string `json:"STKTXT,omitempty"`

		AUDLEN string `json:"AUDLEN,omitempty"`

		MID         string `json:"mid,omitempty"`
		DisplayName string `json:"displayName,omitempty"`
	}

	Operation struct {
		Revision int      `json:"revision"`
		OPType   int      `json:"opType"`
		Params   []string `json:"params"`
	}
)
