FROM cgr.dev/chainguard/go as build

WORKDIR /src/amp
COPY --from=common . ../libs/
COPY --from=utils ./cfn-events/ ../cfn-events/
COPY --from=utils ./rds-provision/ ../rds-provision/
COPY . .

RUN --mount=type=cache,target=/root/go/pkg \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o amp . && \
    cd ../cfn-events && \
    go build -o cfn-events . && \
    cd ../rds-provision && \
    go build -o provision .

FROM cgr.dev/chainguard/wolfi-base
RUN apk add -U --no-cache nodejs npm bash busybox aws-cli jq curl zip && \
    apk del --no-cache wolfi-base apk-tools
RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" && \
    unzip awscliv2.zip && \
    ./aws/install
RUN npm install -g aws-cdk
WORKDIR /app
COPY ./cdk.json ./
COPY ./scripts ./scripts
COPY --from=utils ./scripts/ ./scripts/
COPY --from=build /src/amp/amp .
COPY --from=build /src/cfn-events/cfn-events .
COPY --from=build /src/rds-provision/provision .
CMD [ "/app/scripts/apply.sh" ]