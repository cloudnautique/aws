#!/bin/bash

set -e

if [ ! -f outputs.json ]; then
    echo "outputs.json file not found!"
    exit 1
fi

url=$(jq -r '.[] | select(.OutputKey=="AMPEndpointURL")|.OutputValue' outputs.json)
arn=$(jq -r '.[]| select(.OutputKey=="AMPWorkspaceArn")|.OutputValue' outputs.json)
proto="${url%%://*}"
no_proto="${url#*://}"
address="${no_proto%%/*}"
uri="${no_proto#*$address}"

cat > /run/secrets/output<<EOF
services: amp: {
    default: true
    address: "${address}"
    data: {
        arn: "${arn}"
        url: "${url}"
        proto: "${proto}"
        uri: "${uri}"
    }
}
EOF
