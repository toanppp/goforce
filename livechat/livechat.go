package livechat

import (
	"context"
	"net/http"
)

const (
	pathCreateSession = "chat/rest/System/SessionId"
	pathInitChasitor  = "chat/rest/Chasitor/ChasitorInit"
	pathListMessages  = "chat/rest/System/Messages"
	pathSendMessage   = "chat/rest/Chasitor/ChatMessage"
	pathEndChat       = "chat/rest/Chasitor/ChatEnd"
)

type Livechat interface {
	CreateSession(ctx context.Context) (Session, error)
	InitChasitor(ctx context.Context, header Header, input ChasitorInit) error
	ListMessages(ctx context.Context, header Header) (Messages, error)
	SendMessage(ctx context.Context, header Header, input SendMessageReq) error
	EndChat(ctx context.Context, header Header, input EndChatReq) error
}

type livechat struct {
	domain     string
	version    string
	httpClient *http.Client
}

func New(domain, version string, opts ...Option) Livechat {
	l := &livechat{
		domain:     domain,
		version:    version,
		httpClient: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

type Option func(l *livechat)

func OptionHTTPClient(client *http.Client) func(*livechat) {
	return func(l *livechat) {
		l.httpClient = client
	}
}
