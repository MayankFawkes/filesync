FROM golang:1.18 as build

# Create appuser.
# See https://stackoverflow.com/a/55757473/12429735
ENV USER=appuser
ENV NUID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${NUID}" \
    "${USER}"



RUN go env -w GO111MODULE=off && apt-get update && apt-get install -y ca-certificates
RUN go get github.com/mayankfawkes/filesync

# Build
WORKDIR /go/src/github.com/mayankfawkes/filesync
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/filesync main.go

###############################################################################
# final stage
FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
USER appuser:appuser

ARG APPLICATION="filesync"
ARG DESCRIPTION="Sync all your SSL and other small files around the servers with docker-compose."
ARG PACKAGE="mayankfawkes/filesync"

LABEL maintainer="Mayank Gupta <mkgupta74d@gmail.com>" \
    org.opencontainers.image.ref.name="${PACKAGE}" \
    org.opencontainers.image.authors="Mayank Gupta <@MayankFawkes>" \
    org.opencontainers.image.documentation="https://github.com/${PACKAGE}/README.md" \
    org.opencontainers.image.description="${DESCRIPTION}" \
    org.opencontainers.image.licenses="MIT" \
    org.opencontainers.image.source="https://github.com/${PACKAGE}"

COPY --from=build /go/bin/${APPLICATION} /filesync


ARG git_sha="development"
ENV GIT_SHA=$git_sha

ARG app_version="unknown"
ENV APP_VERSION=$app_version

ENV GIN_MODE=release

ENTRYPOINT ["/filesync"]