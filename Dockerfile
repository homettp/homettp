FROM alpine
ARG TARGETARCH

WORKDIR /app

RUN apk add --no-cache curl

COPY dist/.env ./
COPY dist/${TARGETARCH}/homettp ./

ENV PATH="/app:${PATH}"

EXPOSE 4000

ENTRYPOINT ["homettp"]
CMD ["web", "serve"]
