#!/bin/bash

arn="$(jq -r '.[]| select(.OutputKey=="ApiGatewayARN")|.OutputValue' outputs.json)"
url="$(jq -r '.[]| select(.OutputKey=="ApiGatewayURL")|.OutputValue' outputs.json)"

cat >/run/secrets/output<<EOF
services: gateway: data: {
    arn: "${arn}"
    url: "${url}"
}
EOF