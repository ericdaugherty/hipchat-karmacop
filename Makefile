all: build pack update

build:
	GOOS=linux go build -o main

pack:
	zip deployment.zip main

clean:
	@rm -rf main deployment.zip

update: build pack
		@aws lambda update-function-code                                           \
		  --function-name karma-Function-R4IDE2LPTA4U                                                 \
		  --zip-file fileb://deployment.zip

.PHONY:  all build pack clean update
