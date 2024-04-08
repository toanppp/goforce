package livechat

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Session struct {
	Key               string `json:"key"`
	ID                string `json:"id"`
	ClientPollTimeout int64  `json:"clientPollTimeout"`
	AffinityToken     string `json:"affinityToken"`
}

func (l *livechat) CreateSession(ctx context.Context) (Session, error) {
	url := fmt.Sprintf("%s/%s", l.domain, pathCreateSession)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Session{}, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	header := Header{
		Version: l.version,
	}
	for k, v := range header.Values() {
		req.Header.Add(k, v)
	}

	resp, err := l.httpClient.Do(req)
	if err != nil {
		return Session{}, fmt.Errorf("httpClient.Do: %w", err)
	}

	if resp.StatusCode < http.StatusOK || http.StatusMultipleChoices <= resp.StatusCode {
		pResp, _ := io.ReadAll(resp.Body)
		return Session{}, fmt.Errorf("request failed: %s", string(pResp))
	}

	pResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return Session{}, fmt.Errorf("io.ReadAll: %w", err)
	}

	var output Session
	if err := json.Unmarshal(pResp, &output); err != nil {
		return Session{}, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return output, nil
}
