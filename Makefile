clean:
	rm -rf main main.zip

build: clean
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main
	zip main.zip main

upload: build
	aws s3 cp main.zip s3://source.ap-southeast-2.zac.base2services.com/main.zip

deploy: upload
	aws cloudformation create-stack --stack-name main --template-body file://template.yaml --capabilities CAPABILITY_NAMED_IAM
