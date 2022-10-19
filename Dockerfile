FROM alpine
ARG TARGETARCH

WORKDIR /app

RUN apk add --no-cache libc6-compat curl

COPY dist/.env ./
COPY dist/${TARGETARCH}/homettp ./

EXPOSE 4000

ENTRYPOINT ["./homettp"]
CMD ["web", "serve"]
