package main

import (
	"github.com/acorn-io/services/aws/libs/common"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecr"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/sirupsen/logrus"
)

type LambdaStackProps struct {
	StackProps          awscdk.StackProps
	UserTags            map[string]string `json:"tags"`
	EcrRepoArn          string            `json:"ecrRepoArn"`
	EcrRepo             string            `json:"ecrRepo"`
	EcrRepoTag          string            `json:"ecrRepoTag"`
	LambdaArchitecture  string            `json:"lambdaArchitecture"`
	Timeout             int               `json:"functionTimeout"`
	AdditionalPolicyArn string            `json:"additionalPolicyArn"`
	EnvironmentVars     map[string]string `json:"runtimeEnvironmentVars"`
}

func NewLambdaStack(scope constructs.Construct, id string, props *LambdaStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	dockerFunctionProps := getAwsDockerImageFunctionProps(stack, props)
	function := awslambda.NewDockerImageFunction(stack, jsii.String("lambdaContainerFunction"), dockerFunctionProps)

	if props.AdditionalPolicyArn != "" {
		function.Role().AddManagedPolicy(awsiam.ManagedPolicy_FromManagedPolicyArn(stack, jsii.String("lambdaManagedPolicy"), &props.AdditionalPolicyArn))
	}

	awscdk.NewCfnOutput(stack, jsii.String("FunctionARN"), &awscdk.CfnOutputProps{
		Value: function.FunctionArn(),
	})

	return stack
}

func getAwsDockerImageFunctionProps(scope constructs.Construct, props *LambdaStackProps) *awslambda.DockerImageFunctionProps {
	returns := &awslambda.DockerImageFunctionProps{
		Architecture: getArchitecture(props.LambdaArchitecture),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(props.Timeout)),
		Environment:  common.MapStringToMapStringPtr(props.EnvironmentVars),
	}

	repo := getRepo(scope, props)
	dockerImageFunction := awslambda.DockerImageCode_FromEcr(repo, &awslambda.EcrImageCodeProps{
		Tag: jsii.String(props.EcrRepoTag),
	})
	returns.Code = dockerImageFunction

	return returns
}

func getArchitecture(a string) awslambda.Architecture {
	switch a {
	case "ARM_64":
		return awslambda.Architecture_ARM_64()
	case "X86_64":
		return awslambda.Architecture_X86_64()
	default:
		return awslambda.Architecture_ARM_64()
	}
}

func getRepo(scope constructs.Construct, lambdaProps *LambdaStackProps) awsecr.IRepository {
	var repo awsecr.IRepository

	if lambdaProps.EcrRepoArn != "" {
		repo = awsecr.Repository_FromRepositoryArn(scope, jsii.String("sourceEcrRepo"), jsii.String(lambdaProps.EcrRepoArn))
	} else {
		repo = awsecr.Repository_FromRepositoryName(scope, jsii.String("sourceEcrRepo"), jsii.String(lambdaProps.EcrRepo))
	}

	return repo
}

func main() {
	defer jsii.Close()

	app := common.NewAcornTaggedApp(nil)
	stackProps := &LambdaStackProps{
		StackProps: *common.NewAWSCDKStackProps(),
	}

	if err := common.NewConfig(&stackProps); err != nil {
		logrus.Fatal(err)
	}

	common.AppendScopedTags(app, stackProps.UserTags)

	NewLambdaStack(app, "lambdaStack", stackProps)

	app.Synth(nil)
}
