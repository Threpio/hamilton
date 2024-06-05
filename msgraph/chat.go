package msgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

const (
	ShortTypeChat               string = "chat"
	ShortTypeConversationMember string = "aadUserConversationMember"
	TypeChat                    string = "#microsoft.graph.chat"
	TypeConversationMember      string = "#microsoft.graph.aadUserConversationMember"
)

type ChatClient struct {
	BaseClient Client
}

func NewChatClient() *ChatClient {
	return &ChatClient{
		BaseClient: NewClient(Version10),
	}
}

// Create creates a new chat.
func (c *ChatClient) Create(ctx context.Context, chat Chat) (*Chat, int, error) {
	body, err := json.Marshal(chat)
	if err != nil {
		return nil, 0, fmt.Errorf("json.Marshal(): %v", err)
	}

	retryFunc := func(resp *http.Response, o *odata.OData) bool {
		if resp != nil && o != nil && o.Error != nil {
			if resp.StatusCode == http.StatusForbidden {
				return o.Error.Match("One or more members cannot be added to the thread roster")
			}
		}
		return false
	}

	resp, status, _, err := c.BaseClient.Post(ctx, PostHttpRequestInput{
		Body:                   body,
		ConsistencyFailureFunc: retryFunc,
		OData: odata.Query{
			Metadata: odata.MetadataFull,
		},
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: Uri{
			Entity: "/chats",
		},
	})

	if err != nil {
		return nil, status, fmt.Errorf("ChatsClient.BaseClient.Post(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var newChat Chat
	if err := json.Unmarshal(respBody, &newChat); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &newChat, status, nil
}

// Get retrieves a chat.
func (c *ChatClient) Get(ctx context.Context, id string, query odata.Query) (*Chat, int, error) {
	query.Metadata = odata.MetadataFull

	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		OData:                  query,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: Uri{
			Entity: fmt.Sprintf("/chats/%s", id),
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("ChatsClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var chat Chat
	if err := json.Unmarshal(respBody, &chat); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &chat, status, nil
}

// List returns a list of chats as Chat objects.
// To return just a lost of IDs then place the query to be Odata.Query{Select: "id"}.
func (c *ChatClient) List(ctx context.Context, userID string, query odata.Query) (*[]Chat, int, error) {
	query.Metadata = odata.MetadataFull

	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		OData:                  query,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: Uri{
			Entity: fmt.Sprintf("/users/%s/chats", userID),
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("ChatsClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var chatList struct {
		Value []Chat `json:"value"`
	}
	if err := json.Unmarshal(respBody, &chatList); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &chatList.Value, status, nil

}

// Update updates a chat.
func (c *ChatClient) Update(ctx context.Context, chat Chat) (int, error) {
	body, err := json.Marshal(chat)
	if err != nil {
		return 0, fmt.Errorf("json.Marshal(): %v", err)
	}

	_, status, _, err := c.BaseClient.Patch(ctx, PatchHttpRequestInput{
		Body:                   body,
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusNoContent},
		Uri: Uri{
			Entity: fmt.Sprintf("/chats/%s", *chat.ID),
		},
	})
	if err != nil {
		return status, fmt.Errorf("ChatsClient.BaseClient.Patch(): %v", err)
	}

	return status, nil
}

// Delete deletes a chat.
func (c *ChatClient) Delete(ctx context.Context, chatId string) (int, error) {
	_, status, _, err := c.BaseClient.Delete(ctx, DeleteHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusNoContent},
		Uri: Uri{
			Entity: fmt.Sprintf("/chats/%s", chatId),
		},
	})
	if err != nil {
		return status, fmt.Errorf("ChatsClient.BaseClient.Delete(): %v", err)
	}

	return status, nil
}