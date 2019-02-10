package main

import (
	"context"
)

func main() {
	StartLambda(&handler{})
}

type resourceProperties struct {
	UserNamePrefix string `json:"UserNamePrefix,omitempty"`
}

type responseData struct {
	UserName string `json:"UserName,omitempty"`
}

type handler struct{}

func (h *handler) Create(ctx context.Context, req Request) (physicalResourceID string, data interface{}, err error) {
	p := resourceProperties{}
	err = req.ResourceProperties.Unmarshal(&p)
	if err != nil {
		return "", nil, err
	}

	return "", &responseData{
		UserName: p.UserNamePrefix + "001",
	}, nil
}

func (h *handler) Update(ctx context.Context, req Request) (data interface{}, err error) {
	_, data, err = h.Create(ctx, req)
	return data, err
}

func (h *handler) Delete(ctx context.Context, req Request) error {
	return nil
}
