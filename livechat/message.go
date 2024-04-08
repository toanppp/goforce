package livechat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	MessageTypeChatRequestSuccess   = "ChatRequestSuccess"
	MessageTypeChatRequestFail      = "ChatRequestFail"
	MessageTypeChatEstablished      = "ChatEstablished"
	MessageTypeChatTransferred      = "ChatTransferred"
	MessageTypeChatMessage          = "ChatMessage"
	MessageTypeChatEnded            = "ChatEnded"
	MessageTypeQueueUpdate          = "QueueUpdate"
	MessageTypeAgentTyping          = "AgentTyping"
	MessageTypeAgentNotTyping       = "AgentNotTyping"
	MessageTypeAgentDisconnect      = "AgentDisconnect"
	MessageTypeFileTransfer         = "FileTransfer"
	MessageTypeNewVisitorBreadcrumb = "NewVisitorBreadcrumb"
	MessageTypeCustomEvent          = "CustomEvent"

	ChatRequestFailReasonUnavailable = "Unavailable"

	FieldCaseID    = "CaseId"
	FieldContactID = "ContactId"
)

type Messages struct {
	Messages []Message[any] `json:"messages"`
	Sequence int64          `json:"sequence"`
	Offset   int64          `json:"offset"`
}

func (m *Messages) AssertMessageType() {
	for i := range m.Messages {
		if m.Messages[i].Message == nil {
			continue
		}

		if _, ok := m.Messages[i].Message.(map[string]any); !ok {
			continue
		}

		switch m.Messages[i].Type {
		case MessageTypeChatRequestSuccess:
			m.Messages[i].Message = assertJSON[ChatRequestSuccess](m.Messages[i].Message)
		case MessageTypeChatRequestFail:
			m.Messages[i].Message = assertJSON[ChatRequestFail](m.Messages[i].Message)
		case MessageTypeChatEstablished:
			m.Messages[i].Message = assertJSON[ChatEstablished](m.Messages[i].Message)
		case MessageTypeChatTransferred:
			m.Messages[i].Message = assertJSON[ChatTransferred](m.Messages[i].Message)
		case MessageTypeChatMessage:
			m.Messages[i].Message = assertJSON[ChatMessage](m.Messages[i].Message)
		case MessageTypeChatEnded:
			m.Messages[i].Message = assertJSON[ChatEnded](m.Messages[i].Message)
		case MessageTypeQueueUpdate:
			m.Messages[i].Message = assertJSON[QueueUpdate](m.Messages[i].Message)
		case MessageTypeAgentTyping, MessageTypeAgentNotTyping, MessageTypeAgentDisconnect:
			m.Messages[i].Message = assertJSON[Agent](m.Messages[i].Message)
		case MessageTypeFileTransfer:
			m.Messages[i].Message = assertJSON[FileTransfer](m.Messages[i].Message)
		case MessageTypeNewVisitorBreadcrumb:
			m.Messages[i].Message = assertJSON[NewVisitorBreadcrumb](m.Messages[i].Message)
		case MessageTypeCustomEvent:
			m.Messages[i].Message = assertJSON[CustomEvent](m.Messages[i].Message)
		default:
			continue
		}
	}
}

func assertJSON[T any](input any) any {
	p, err := json.Marshal(input)
	if err != nil {
		return input
	}

	var output T
	if err := json.Unmarshal(p, &output); err != nil {
		return input
	}

	return output
}

type Message[T any] struct {
	Type    string `json:"type"`
	Message T      `json:"message"`
}

type ChatRequestSuccess struct {
	ConnectionTimeout     int64               `json:"connectionTimeout"`
	EstimatedWaitTime     int64               `json:"estimatedWaitTime"`
	SensitiveDataRules    []SensitiveDataRule `json:"sensitiveDataRules"`
	TranscriptSaveEnabled bool                `json:"transcriptSaveEnabled"`
	URL                   string              `json:"url"`
	QueuePosition         int64               `json:"queuePosition"`
	CustomDetails         []CustomerDetail    `json:"customDetails"`
	VisitorID             string              `json:"visitorId"`
	GeoLocation           GeoLocation         `json:"geoLocation"`
}

type SensitiveDataRule struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CustomerDetail struct {
	Label            string   `json:"label"`
	Value            string   `json:"value"`
	TranscriptFields []string `json:"transcriptFields"`
	DisplayToAgent   bool     `json:"displayToAgent"`
}

type GeoLocation struct {
	Organization string  `json:"organization"`
	Region       string  `json:"region"`
	City         string  `json:"city"`
	CountryName  string  `json:"countryName"`
	Latitude     float64 `json:"latitude"`
	CountryCode  string  `json:"countryCode"`
	Longitude    float64 `json:"longitude"`
}

type ChatRequestFail struct {
	Reason      string `json:"reason"`
	PostChatURL string `json:"postChatUrl"`
}

type ChatEstablished struct {
	Name                string              `json:"name"`
	UserID              string              `json:"userId"`
	SneakPeekEnabled    bool                `json:"sneakPeekEnabled"`
	ChasitorIdleTimeout ChasitorIdleTimeout `json:"chasitorIdleTimeout"`
}

type ChasitorIdleTimeout struct {
	IsEnabled   bool  `json:"isEnabled"`
	WarningTime int64 `json:"warningTime"`
	Timeout     int64 `json:"timeout"`
}

type ChatTransferred struct {
	Name             string `json:"name"`
	UserID           string `json:"userId"`
	SneakPeekEnabled bool   `json:"sneakPeekEnabled"`
}

type ChatMessage struct {
	Text     string   `json:"text"`
	Name     string   `json:"name"`
	Schedule Schedule `json:"schedule"`
	AgentID  string   `json:"agentId"`
}

type Schedule struct {
	ResponseDelayMilliseconds float64 `json:"responseDelayMilliseconds"`
}

type ChatEnded struct {
	AttachedRecords []AttachedRecord `json:"attachedRecords"`
	Reason          string           `json:"reason"`
}

type AttachedRecord struct {
	FieldValue string `json:"fieldValue"`
	FieldName  string `json:"fieldName"`
}

type QueueUpdate struct {
	EstimatedWaitTime int64 `json:"estimatedWaitTime"`
	Position          int64 `json:"position"`
}

type Agent struct {
	Name    string `json:"name"`
	AgentID string `json:"agentId"`
}

type FileTransfer struct {
	UploadServletURL string `json:"uploadServletUrl"`
	FileToken        string `json:"fileToken"`
	CdmServletURL    string `json:"cdmServletUrl"`
	Type             string `json:"type"`
}

type CustomEvent struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type NewVisitorBreadcrumb struct {
	Location string `json:"location"`
}

type SendMessageReq struct {
	Text string `json:"text"`
}

func (l *livechat) ListMessages(ctx context.Context, header Header) (Messages, error) {
	url := fmt.Sprintf("%s/%s", l.domain, pathListMessages)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Messages{}, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	header.Version = l.version
	for k, v := range header.Values() {
		req.Header.Add(k, v)
	}

	resp, err := l.httpClient.Do(req)
	if err != nil {
		return Messages{}, fmt.Errorf("httpClient.Do: %w", err)
	}

	if resp.StatusCode < http.StatusOK || http.StatusMultipleChoices <= resp.StatusCode {
		pResp, _ := io.ReadAll(resp.Body)
		return Messages{}, fmt.Errorf("request failed: %s", string(pResp))
	}

	pResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return Messages{}, fmt.Errorf("io.ReadAll: %w", err)
	}

	var output Messages
	if err := json.Unmarshal(pResp, &output); err != nil {
		return Messages{}, fmt.Errorf("json.Unmarshal: %w", err)
	}

	output.AssertMessageType()
	return output, nil
}

func (l *livechat) SendMessage(ctx context.Context, header Header, input SendMessageReq) error {
	url := fmt.Sprintf("%s/%s", l.domain, pathSendMessage)
	pReq, _ := json.Marshal(input)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(pReq))
	if err != nil {
		return fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	header.Version = l.version
	for k, v := range header.Values() {
		req.Header.Add(k, v)
	}

	resp, err := l.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("httpClient.Do: %w", err)
	}

	if resp.StatusCode < http.StatusOK || http.StatusMultipleChoices <= resp.StatusCode {
		pResp, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed: %s", string(pResp))
	}

	return nil
}
