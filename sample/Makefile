clean:
	rm -rf _build

build: clean
	mkdir -p _build
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o cfn-customres-go-sample
	zip _build/cfn-customres-go-sample cfn-customres-go-sample
	rm cfn-customres-go-sample

upload: build
	aws s3 cp _build/cfn-customres-go-sample.zip s3://$(BUCKET_NAME)/cfn-customres-go-sample.zip

template:
	sed -e "s/{{BUCKET_NAME}}/$(BUCKET_NAME)/g" template.yaml > _build/template.yaml

deploy: upload template
	aws cloudformation create-stack --stack-name $(STACK_NAME) --template-body file://_build/template.yaml --capabilities CAPABILITY_NAMED_IAM