FROM gcr.io/distroless/base-debian11
ARG TARGETARCH

WORKDIR /app

COPY dist/.env ./
COPY dist/${TARGETARCH}/homettp ./

EXPOSE 4000

ENTRYPOINT ["./homettp"]
CMD ["web", "serve"]
