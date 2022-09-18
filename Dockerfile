FROM gcr.io/distroless/base-debian11:debug
ARG TARGETARCH

WORKDIR /app

COPY dist/.env ./
COPY dist/${TARGETARCH}/homettp ./

RUN /busybox/sh -c ln -s /busybox/sh /bin/sh

EXPOSE 4000

ENTRYPOINT ["./homettp"]
CMD ["web", "serve"]
