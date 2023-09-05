package main

import (
	"github.com/acorn-io/services/aws/libs/common"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/sirupsen/logrus"
)

type ApiGatewayStackProps struct {
	StackProps  awscdk.StackProps
	FunctionArn string            `json:"functionArn"`
	UserTags    map[string]string `json:"tags"`
}

func NewApiGatewayStack(scope constructs.Construct, id string, props *ApiGatewayStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// This should be just HTTP... need to look at v2
	api := awsapigateway.NewRestApi(stack, jsii.String("api"), &awsapigateway.RestApiProps{
		RestApiName: jsii.String("api"),
		Deploy:      jsii.Bool(true),
	})

	lambda := awslambda.Function_FromFunctionArn(stack, jsii.String("lambda"), jsii.String(props.FunctionArn))
	apiIntegration := awsapigateway.NewLambdaIntegration(lambda, &awsapigateway.LambdaIntegrationOptions{
		Proxy: jsii.Bool(true),
	})

	api.Root().AddMethod(jsii.String("ANY"), apiIntegration, nil)

	awscdk.NewCfnOutput(stack, jsii.String("ApiGatewayURL"), &awscdk.CfnOutputProps{
		Value: api.Url(),
	})

	awscdk.NewCfnOutput(stack, jsii.String("ApiGatewayARN"), &awscdk.CfnOutputProps{
		Value: api.RestApiId(),
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := common.NewAcornTaggedApp(nil)
	stackProps := &ApiGatewayStackProps{
		StackProps: *common.NewAWSCDKStackProps(),
	}

	if err := common.NewConfig(&stackProps); err != nil {
		logrus.Fatal(err)
	}

	common.AppendScopedTags(app, stackProps.UserTags)

	NewApiGatewayStack(app, "apiGatewayStack", stackProps)

	app.Synth(nil)
}
