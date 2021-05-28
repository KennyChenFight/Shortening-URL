FROM gcr.io/distroless/base-debian10 AS base
ENV MIGRATION_FILE_DIR /migrations
COPY migrations $MIGRATION_DIR

FROM base AS shortening-url-server

ENV PORT 8080
EXPOSE $PORT

COPY bin/cmd/server /server
CMD ["/server"]

FROM base AS shortening-url-cron
COPY bin/cmd/cron /cron
CMD ["/cron"]