FROM golang:alpine AS build-env
RUN apk --no-cache add build-base git bzr mercurial gcc
ADD . /src
RUN cd /src && go build -o jira-api-exporter ./cmd/jira-api-exporter

# final stage
FROM alpine
RUN apk --no-cache add ca-certificates tzdata
COPY --from=build-env /src/jira-api-exporter /usr/local/bin/
CMD ["jira-api-exporter"]
