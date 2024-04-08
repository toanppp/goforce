package livechat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ChasitorInit struct {
	OrganizationID      string          `json:"organizationId"`
	DeploymentID        string          `json:"deploymentId"`
	ButtonID            string          `json:"buttonId"`
	AgentID             string          `json:"agentId"`
	DoFallback          bool            `json:"doFallback"`
	SessionID           string          `json:"sessionId"`
	UserAgent           string          `json:"userAgent"`
	Language            string          `json:"language"`
	ScreenResolution    string          `json:"screenResolution"`
	VisitorName         string          `json:"visitorName"`
	PrechatDetails      []PrechatDetail `json:"prechatDetails"`
	PrechatEntities     []PrechatEntity `json:"prechatEntities"`
	ButtonOverrides     []string        `json:"buttonOverrides"`
	ReceiveQueueUpdates bool            `json:"receiveQueueUpdates"`
	IsPost              bool            `json:"isPost"`
}

type PrechatDetail struct {
	Label             string   `json:"label"`
	Value             string   `json:"value"`
	TranscriptFields  []string `json:"transcriptFields"`
	DisplayToAgent    bool     `json:"displayToAgent"`
	DoKnowledgeSearch bool     `json:"doKnowledgeSearch"`
}

type PrechatEntity struct {
	EntityName        string            `json:"entityName"`
	ShowOnCreate      bool              `json:"showOnCreate"`
	LinkToEntityName  string            `json:"linkToEntityName,omitempty"`
	LinkToEntityField string            `json:"linkToEntityField,omitempty"`
	SaveToTranscript  string            `json:"saveToTranscript"`
	EntityFieldsMaps  []EntityFieldsMap `json:"entityFieldsMaps"`
}

type EntityFieldsMap struct {
	FieldName    string `json:"fieldName"`
	Label        string `json:"label"`
	DoFind       bool   `json:"doFind"`
	IsExactMatch bool   `json:"isExactMatch"`
	DoCreate     bool   `json:"doCreate"`
}

type EndChatReq struct {
	Reason string `json:"reason"`
}

func (l *livechat) InitChasitor(ctx context.Context, header Header, input ChasitorInit) error {
	url := fmt.Sprintf("%s/%s", l.domain, pathInitChasitor)
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

func (l *livechat) EndChat(ctx context.Context, header Header, input EndChatReq) error {
	url := fmt.Sprintf("%s/%s", l.domain, pathEndChat)
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
