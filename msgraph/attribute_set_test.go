package msgraph_test

import (
	"net/http"
	"testing"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
)

const attributeSetId = "test"

func TestAttributeSetClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	c.AttributeSetClient.BaseClient.DisableRetries = true
	_, status, err := c.AttributeSetClient.Get(c.Context, attributeSetId, odata.Query{})
	c.AttributeSetClient.BaseClient.DisableRetries = false
	if err != nil {
		if status != http.StatusNotFound {
			t.Fatalf("AttributeSetClient.Get(): unable to retrieve attribute set for testing: %v", err)
		}

		testAttributeSetClientCreate(
			t,
			c,
			msgraph.AttributeSet{
				Description: utils.StringPtr("test attribute set"),
				ID:          utils.StringPtr(attributeSetId),
			},
		)
	}

	testAttributeSetClientGet(t, c, attributeSetId)
	testAttributeSetClientUpdate(
		t,
		c,
		msgraph.AttributeSet{
			ID:          utils.StringPtr(attributeSetId),
			Description: utils.StringPtr("updated test description"),
		},
	)

	testAttributeSetClientList(t, c)
}

func testAttributeSetClientCreate(t *testing.T, c *test.Test, csad msgraph.AttributeSet) *msgraph.AttributeSet {

	attributeSet, status, err := c.AttributeSetClient.Create(c.Context, csad)
	if err != nil {
		t.Fatalf("AttributeSetClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AttributeSetClient.Create(): invalid status:%d", status)
	}
	if attributeSet == nil {
		t.Fatalf("AttributeSet.Create(): attributeSet was nil")
	}
	if attributeSet.ID == nil {
		t.Fatalf("AttributeSetClient.Create(): attributeSet.ID was nil")
	}

	return attributeSet
}

func testAttributeSetClientGet(t *testing.T, c *test.Test, id string) *msgraph.AttributeSet {
	attributeSet, status, err := c.AttributeSetClient.Get(c.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("AttributeSetClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AttributeSet.Client.Get(): invalid status: %d", status)
	}
	if attributeSet == nil {
		t.Fatalf("AttributeSetClient.Get(): attributeSet was nil")
	}

	return attributeSet
}

func testAttributeSetClientList(t *testing.T, c *test.Test) *[]msgraph.AttributeSet {
	attributeSets, _, err := c.AttributeSetClient.List(
		c.Context,
		odata.Query{Top: 10},
	)
	if err != nil {
		t.Fatalf("AttributeSetClient.List(): %v", err)
	}
	if attributeSets == nil {
		t.Fatalf("AttributeSetClient.List(): attributeSets was nil")
	}

	return attributeSets
}

func testAttributeSetClientUpdate(t *testing.T, c *test.Test, csad msgraph.AttributeSet) {
	status, err := c.AttributeSetClient.Update(c.Context, csad)
	if err != nil {
		t.Fatalf("AttributeSetClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AttributeSetClient.Update(): invalid status: %d", status)
	}
}
