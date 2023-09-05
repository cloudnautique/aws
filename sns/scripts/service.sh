#!/bin/bash

STACK_NAME="${ACORN_EXTERNAL_ID}"
# Check if outputs.json exists
if [ ! -f outputs.json ]; then
    echo "No outputs.json found. Exiting."
    exit 1
fi


# Render Output
arn="$(jq -r '.[] | select(.OutputKey=="TopicARN")|.OutputValue' outputs.json )"

cat > /run/secrets/output <<EOF
services: {
    admin: {
        default: true
        address: "${address}"
        consumer: permissions: rules: [{
            apiGroup: "aws.acorn.io"
		    verbs: [
			    "sns:*",
		    ]
		    resources: ["${arn}"]
        }]
        data: {
            arn: "${arn}"
            name: "${name}"
        }
    }
    publisher: {
        address: "${address}"
        consumer: permissions: rules: [{
            apiGroup: "aws.acorn.io"
            verbs: [
                "sns:Publish"
            ]
            resources: ["${arn}"]
        }]
        data: {
            arn: "${arn}"
        }
    }
    subscriber: {
        address: "${address}"
        consumer: permissions: rules: [{
            apiGroup: "aws.acorn.io"
            verbs: [
                "sns:Subscribe",
            ]
            resources: ["${arn}"]
        }]
        data: {
            arn: "${arn}"
        }
    }
}
EOF
