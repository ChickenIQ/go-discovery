package discovery

type Member struct {
	Key       string `json:"key"`
	Metadata  string `json:"metadata"`
	Signature string `json:"signature"`
}

type Body struct {
	Data      string `json:"data"`
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature"`
}

type Entry struct {
	Member Member `json:"member"`
	Body   Body   `json:"body"`
}
