#!/bin/bash

if [ ! -f outputs.json ]; then
   echo "outputs.json file not found!"
   exit 1
fi

# Render Output
CLUSTER_NAME=$(jq -r '.[] | select(.OutputKey=="clustername")|.OutputValue' outputs.json)
ADDRESS=$(jq -r '.[] | select(.OutputKey=="address")|.OutputValue' outputs.json)
PORT=$(jq -r '.[]| select(.OutputKey=="port")|.OutputValue' outputs.json)
TOKEN_ARN=$(jq -r '.[]| select(.OutputKey=="tokenarn")|.OutputValue' outputs.json)

TOKEN="$(aws --output json secretsmanager get-secret-value --secret-id "${TOKEN_ARN}" --query 'SecretString' | jq -r .|jq -r ."${CLUSTER_NAME}-token")"

cat > /run/secrets/output<<EOF
services: redis: {
  default: true
  address: "${ADDRESS}"
  ports: [${PORT}]
  data: {
    clusterName: "${CLUSTER_NAME}"
    address: "${ADDRESS}"
    port: "${PORT}"
  }
}

secrets: "admin": {
  type: "token"
  data: {
    token: "${TOKEN}"
  }
}
EOF
