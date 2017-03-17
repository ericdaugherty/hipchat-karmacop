# hipchat-karmacop

This is a simple HipChat AddOn that checks for users who use the find and replace functionality in HipChat to secretly give or take Karma using the HipChat Karma plugin.

It is written in Go and setup to deploy to AWS Lamba.

Please refer to the https://github.com/eawsy/aws-lambda-go-shim project for a summary of how the Makefile works along with how to setup a Lamba Function to work properly with the shim.

This project utilizes S3 and assumes a bucket named "karmacoprooms". You can edit aws.go to change the bucket name.

You will also need to grant your Lambda IAM Role permissions to read and write to the S3 bucket.

To publish the Lambda function as web API you need to create an Amazon API Gateway and use a LAMBDA_PROXY Integration Request and reference your Lambda function. You then need to change the base URL in the init function in handler.go

Here is a blog post with more background on this project: http://blog.ericdaugherty.com/2017/03/simple-hipchat-addon-in-go.html