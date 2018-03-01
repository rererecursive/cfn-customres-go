package cusres_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/nfukasawa/cfn-customres-go/cusres"
)

func TestCustomResource(t *testing.T) {

	type something struct {
		Name  string     `json:"name,omitempty"`
		Child *something `json:"child,omitempty"`
	}

	cases := []testCase{
		{
			cusres.RequestTypeCreate,
			&something{Name: "foo"},
			handler{
				t: t,
				createHandler: func(ctx context.Context, req cusres.Request) (string, interface{}, error) {
					x := something{}
					req.ResourceProperties.Unmarshal(&x)
					if x.Name != "foo" {
						t.Fatalf("wants %v, but got %v", "foo", x.Name)
					}
					return "", &something{Child: &something{Name: "bar"}}, nil
				},
			},
			cusres.Response{
				Status:            cusres.ResponseStatusSuccess,
				StackID:           stackID,
				RequestID:         requestID,
				LogicalResourceID: logicalResourceID,
				NoEcho:            false,
				Data:              map[string]interface{}{"child.name": "bar"},
			},
			http.StatusOK,
			false,
		},

		{
			cusres.RequestTypeCreate,
			&something{Name: "foo"},
			handler{
				t: t,
				createHandler: func(ctx context.Context, req cusres.Request) (string, interface{}, error) {
					return "xyz", nil, nil
				},
			},
			cusres.Response{
				Status:             cusres.ResponseStatusSuccess,
				StackID:            stackID,
				RequestID:          requestID,
				LogicalResourceID:  logicalResourceID,
				PhysicalResourceID: "xyz",
				NoEcho:             false,
				Data:               nil,
			},
			http.StatusOK,
			false,
		},

		{
			cusres.RequestTypeCreate,
			&something{Name: "foo"},
			handler{
				t: t,
				createHandler: func(ctx context.Context, req cusres.Request) (string, interface{}, error) {
					return "", &something{Name: "bar"}, fmt.Errorf("err")
				},
			},
			cusres.Response{
				Status:            cusres.ResponseStatusFailed,
				Reason:            "err",
				StackID:           stackID,
				RequestID:         requestID,
				LogicalResourceID: logicalResourceID,
				NoEcho:            false,
				Data:              nil,
			},
			http.StatusOK,
			false,
		},

		{
			cusres.RequestTypeUpdate,
			&something{Name: "foo"},
			handler{
				t: t,
				updateHandler: func(ctx context.Context, req cusres.Request) (interface{}, error) {
					return &something{Child: &something{Name: "bar"}}, nil
				},
			},
			cusres.Response{
				Status:            cusres.ResponseStatusSuccess,
				StackID:           stackID,
				RequestID:         requestID,
				LogicalResourceID: logicalResourceID,
				NoEcho:            false,
				Data:              map[string]interface{}{"child.name": "bar"},
			},
			http.StatusOK,
			false,
		},

		{
			cusres.RequestTypeDelete,
			&something{Name: "foo"},
			handler{
				t: t,
				deleteHandler: func(ctx context.Context, req cusres.Request) error {
					return nil
				},
			},
			cusres.Response{
				Status:            cusres.ResponseStatusSuccess,
				StackID:           stackID,
				RequestID:         requestID,
				LogicalResourceID: logicalResourceID,
				NoEcho:            false,
				Data:              nil,
			},
			http.StatusOK,
			false,
		},

		{
			cusres.RequestTypeDelete,
			&something{Name: "foo"},
			handler{
				t: t,
				deleteHandler: func(ctx context.Context, req cusres.Request) error {
					return nil
				},
			},
			cusres.Response{
				Status:            cusres.ResponseStatusSuccess,
				StackID:           stackID,
				RequestID:         requestID,
				LogicalResourceID: logicalResourceID,
				NoEcho:            false,
				Data:              nil,
			},
			http.StatusInternalServerError,
			true,
		},
	}

	for _, c := range cases {
		testCustomResource(t, c)
	}
}

type testCase struct {
	requestType   cusres.RequestType
	resourceProps interface{}

	handler handler

	wantedCustomResource cusres.Response

	cfnHTTPStatus     int
	lambdaErrorwanted bool
}

func testCustomResource(t *testing.T, c testCase) {
	ctx := context.Background()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		dec := json.NewDecoder(r.Body)
		res := cusres.Response{}
		dec.Decode(&res)

		if res.PhysicalResourceID == "" {
			t.Fatalf("PhysicalResourceID should not be empty")
		}
		if c.wantedCustomResource.PhysicalResourceID != "" && c.wantedCustomResource.PhysicalResourceID != res.PhysicalResourceID {
			t.Fatalf("PhysicalResourceID: wants %v, but got %v", c.wantedCustomResource.PhysicalResourceID, res.PhysicalResourceID)
		}
		c.wantedCustomResource.PhysicalResourceID = ""
		res.PhysicalResourceID = ""
		if !reflect.DeepEqual(c.wantedCustomResource, res) {
			t.Fatalf("CFn response: wants %v, but got %v", c.wantedCustomResource, res)
		}
		w.WriteHeader(c.cfnHTTPStatus)
	}))
	defer srv.Close()

	b, _ := json.Marshal(c.resourceProps)
	props := cusres.Properties{}
	json.Unmarshal(b, &props)

	phyID := ""
	if c.requestType != cusres.RequestTypeCreate {
		phyID = "phyid"
	}
	err := cusres.NewLambdaHandler(c.handler)(ctx, cusres.Request{
		RequestType:        c.requestType,
		ResponseURL:        srv.URL,
		StackID:            stackID,
		RequestID:          requestID,
		ResourceType:       resourceType,
		LogicalResourceID:  logicalResourceID,
		PhysicalResourceID: phyID,
		ResourceProperties: props,
	})
	if !c.lambdaErrorwanted && err != nil {
		t.Fatalf("Lambda: wants no error, but got error %v", err)
	}
	if c.lambdaErrorwanted && err == nil {
		t.Fatalf("Lambda: wants error, but got no error")
	}
}

const (
	stackID           = "arn:aws:cloudformation:region:account:stack/stack-name/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx"
	requestID         = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx"
	resourceType      = "Custom::TestCustomResource"
	logicalResourceID = "TestCustomResource"
)

type handler struct {
	t             *testing.T
	createHandler func(ctx context.Context, req cusres.Request) (string, interface{}, error)
	updateHandler func(ctx context.Context, req cusres.Request) (interface{}, error)
	deleteHandler func(ctx context.Context, req cusres.Request) error
}

func (h handler) Create(ctx context.Context, req cusres.Request) (physicalResourceID string, data interface{}, err error) {
	if h.createHandler == nil {
		h.t.Fatalf("Create handler should not be called")
	}
	return h.createHandler(ctx, req)
}
func (h handler) Update(ctx context.Context, req cusres.Request) (data interface{}, err error) {
	if h.updateHandler == nil {
		h.t.Fatalf("Update handler should not be called")
	}
	return h.updateHandler(ctx, req)
}
func (h handler) Delete(ctx context.Context, req cusres.Request) error {
	if h.deleteHandler == nil {
		h.t.Fatalf("Delete handler should not be called")
	}
	return h.deleteHandler(ctx, req)
}
