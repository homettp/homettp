FROM alpine
ARG TARGETARCH

WORKDIR /app

RUN apk add --no-cache curl

COPY dist/.env ./
COPY dist/${TARGETARCH}/homettp ./

EXPOSE 4000

ENTRYPOINT ["/app/homettp"]
CMD ["web", "serve"]
