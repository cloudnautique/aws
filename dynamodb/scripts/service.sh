#!/bin/bash

STACK_NAME="${ACORN_EXTERNAL_ID}"
# Check if outputs.json exists
if [ ! -f outputs.json ]; then
    echo "No outputs.json found. Exiting."
    exit 1
fi

# Render Output
arn="$(jq -r '.[] | select(.OutputKey=="TableARN")|.OutputValue' outputs.json )"

cat > /run/secrets/output <<EOF
services: {
    admin: {
        consumer: permissions: rules: [{
		    apiGroup: "aws.acorn.io"
		    verbs: ["dynamodb:*"]
		    resources: ["${arn}"]
	    }]
        data: {
            arn: "${arn}"
        }
    }
    write: {
        consumer: permissions: rules: [{
		    apiGroup: "aws.acorn.io"
		    verbs: [
		    	"dynamodb:BatchWriteItem",
		    	"dynamodb:PutItem",
		    	"dynamodb:UpdateItem",
		    	"dynamodb:DeleteItem",
		    	"dynamodb:DescribeTable",
		    ]
		    resources: ["${arn}"]
	    }]
        data: {
            arn: "${arn}"
        }
    }
}
EOF