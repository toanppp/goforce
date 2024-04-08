package livechat_test

import (
	"context"
	"fmt"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/toanppp/goforce/livechat"
)

var (
	domain         = os.Getenv("LIVECHAT_DOMAIN")
	version        = os.Getenv("LIVECHAT_VERSION")
	organizationID = os.Getenv("LIVECHAT_ORGANIZATION_ID")
	deploymentID   = os.Getenv("LIVECHAT_DEPLOYMENT_ID")
	buttonID       = os.Getenv("LIVECHAT_BUTTON_ID")
	agentID        = os.Getenv("LIVECHAT_AGENT_ID")
	contactID      = os.Getenv("LIVECHAT_CONTAC_ID")
)

func TestLivechat(t *testing.T) {
	l := livechat.New(domain, version)
	ctx := context.Background()

	// 1. Create Session
	session, err := l.CreateSession(ctx)
	if err != nil {
		t.Fatalf("CreateSession: %+v", err)
	}
	fmt.Printf("session: %+v\n", session)

	header := livechat.Header{
		Version:    version,
		Affinity:   session.AffinityToken,
		SessionKey: session.Key,
		Sequence:   1,
	}
	// 2. Init Chasitor
	chasitorInitReq := livechat.ChasitorInit{
		OrganizationID:   organizationID,
		DeploymentID:     deploymentID,
		ButtonID:         buttonID,
		AgentID:          agentID,
		DoFallback:       true,
		SessionID:        session.ID,
		UserAgent:        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.5.6; rv:5.2) Gecko/20100101 Firefox/5.2.1",
		Language:         "en-US",
		ScreenResolution: "900x1080",
		VisitorName:      "Ryan P",
		PrechatDetails: []livechat.PrechatDetail{
			{
				Label:             "ContactId",
				Value:             contactID,
				TranscriptFields:  []string{},
				DisplayToAgent:    true,
				DoKnowledgeSearch: false,
			},
			{
				Label:             "Contact Name",
				Value:             contactID,
				TranscriptFields:  []string{},
				DisplayToAgent:    true,
				DoKnowledgeSearch: false,
			},
			{
				Label:             "Subject",
				Value:             "Case Title",
				TranscriptFields:  []string{},
				DisplayToAgent:    true,
				DoKnowledgeSearch: false,
			},
		},
		PrechatEntities: []livechat.PrechatEntity{
			{
				EntityName:        "Contact",
				ShowOnCreate:      true,
				LinkToEntityName:  "Contact",
				LinkToEntityField: "Id",
				SaveToTranscript:  "Contact",
				EntityFieldsMaps: []livechat.EntityFieldsMap{
					{
						FieldName:    "Id",
						Label:        "ContactId",
						DoFind:       true,
						IsExactMatch: true,
						DoCreate:     false,
					},
				},
			},
			{
				EntityName:       "Case",
				ShowOnCreate:     true,
				SaveToTranscript: "CaseId",
				EntityFieldsMaps: []livechat.EntityFieldsMap{
					{
						FieldName:    "ContactId",
						Label:        "Contact Name",
						DoFind:       false,
						IsExactMatch: false,
						DoCreate:     true,
					},
					{
						FieldName:    "Subject",
						Label:        "Subject",
						DoFind:       false,
						IsExactMatch: false,
						DoCreate:     true,
					},
				},
			},
		},
		ButtonOverrides: []string{
			agentID,
			buttonID,
		},
		ReceiveQueueUpdates: true,
		IsPost:              true,
	}
	if err := l.InitChasitor(ctx, header, chasitorInitReq); err != nil {
		t.Fatalf("InitChasitor: %+v", err)
	}

	time.Sleep(time.Second)

	// 3. List Messages
	messages, err := l.ListMessages(ctx, header)
	if err != nil {
		t.Fatalf("ListMessages: %+v", err)
	}

	if len(messages.Messages) == 0 {
		t.Fatal("ListMessages: empty messages")
	}

	messageTypes := make([]string, len(messages.Messages))
	for i, v := range messages.Messages {
		messageTypes[i] = v.Type
	}

	if !slices.Contains(messageTypes, livechat.MessageTypeChatRequestSuccess) {
		t.Fatal("ListMessages: chat request not success")
	}

	if !slices.Contains(messageTypes, livechat.MessageTypeChatEstablished) {
		t.Fatal("ListMessages: chat not established")
	}

	fmt.Printf("messages: %+v\n", messages)

	// 4. Send Message
	header.Sequence++
	sendMessageReq := livechat.SendMessageReq{
		Text: "Hello world",
	}
	if err := l.SendMessage(ctx, header, sendMessageReq); err != nil {
		t.Fatalf("SendMessage: %+v", err)
	}

	// 5. End Chat
	header.Sequence++
	endChatReq := livechat.EndChatReq{
		Reason: "client",
	}
	if err := l.EndChat(ctx, header, endChatReq); err != nil {
		t.Fatalf("SendMessage: %+v", err)
	}
}
