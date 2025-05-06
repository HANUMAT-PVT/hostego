package messageclient

type Message struct {
	ID     string
	Values map[string]interface{}
}

type MessagesWithError struct {
	Messages []Message
	Err      error
}

type SendRequest struct {
	Stream       string
	Body         string
	Attributes   []Attribute
	StreamLength int64
}

type Attribute struct {
	Key   string
	Value string
	Type  string
}
