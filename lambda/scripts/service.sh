#!/bin/bash

arn="$(jq -r '.[]| select(.OutputKey=="FunctionARN")|.OutputValue' outputs.json)"

cat >/run/secrets/output<<EOF
services: function: data: arn: "${arn}"
EOF