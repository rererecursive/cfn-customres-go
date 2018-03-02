# cfn-customres-go

Simple utility for creating AWS CloudFormation Lambda-backed Custom Resource in Go

[![Build Status](https://travis-ci.org/nfukasawa/cfn-customres-go.svg?branch=master)](https://travis-ci.org/nfukasawa/cfn-customres-go)
[![Coverage Status](https://coveralls.io/repos/github/nfukasawa/cfn-customres-go/badge.svg?branch=master)](https://coveralls.io/github/nfukasawa/cfn-customres-go?branch=master)

## Install
```
go get github.com/nfukasawa/cfn-customres-go
```
## Run Sample
```shell
$ cd sample
$ BUCKET_NAME=your_bucket_name STACK_NAME=your_stack_name make deploy
```

## Usage

Import the packege.
```go
import (
	"context"

	"github.com/nfukasawa/cfn-customres-go/cusres"
)
```

Define the resource properties and response data structures if required.
```go

type resourceProperties struct {
	UserNamePrefix string `json:"UserNamePrefix,omitempty"`
}

type responseData struct {
	UserName string `json:"UserName,omitempty"`
}
```

Define the custom resourse handler.
```go

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

```

Run the Lambda function.
```go
func main() {
	cusres.StartLambda(&handler{})
}
```
