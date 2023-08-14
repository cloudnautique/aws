FROM cgr.dev/chainguard/go as build
WORKDIR /src/rds-provision
COPY . ./
COPY --from=cdk . ../cdk/
COPY --from=amp . ../amp/
RUN go build -o provision && \
    cd ../cdk && \
    go build -o cdk-context && \
    cd ../amp && \
    go build -o amp

FROM cgr.dev/chainguard/wolfi-base as app
COPY --from=build /src/rds-provision/provision /app/
COPY --from=build /src/cdk/cdk-context /app/
COPY --from=build /src/amp/amp /app/
RUN apk add -U --no-cache nodejs bash && \
    npm install -g aws-cdk
ENV CONFIG_FILE "/app/config.json"
ENV SERVICE_OUTFILE "/run/secrets/output"
ENV CDK_CONTEXT_OUTFILE "/dev/null"
USER 1000
WORKDIR "/app"
CMD ["/app/sculpt"]
