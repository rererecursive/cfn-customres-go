package main

import (
	"context"

	"github.com/nfukasawa/cfn-customres-go/cusres"
)

func main() {
	cusres.StartLambda(&handler{})
}

type resourceProperties struct {
	UserNamePrefix string `json:"UserNamePrefix,omitempty"`
}

type responseData struct {
	UserName string `json:"UserName,omitempty"`
}

type handler struct{}

func (h *handler) Create(ctx context.Context, req cusres.Request) (physicalResourceID string, data interface{}, err error) {
	p := resourceProperties{}
	err = req.ResourceProperties.Unmarshal(&p)
	if err != nil {
		return "", nil, err
	}

	return "", &responseData{
		UserName: p.UserNamePrefix + "001",
	}, nil
}

func (h *handler) Update(ctx context.Context, req cusres.Request) (data interface{}, err error) {
	_, data, err = h.Create(ctx, req)
	return data, err
}

func (h *handler) Delete(ctx context.Context, req cusres.Request) error {
	return nil
}
