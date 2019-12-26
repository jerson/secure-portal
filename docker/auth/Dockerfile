FROM jerson/go:1.13 AS builder

ENV WORKDIR /app
WORKDIR ${WORKDIR}

COPY go.mod go.sum Makefile ./
RUN make deps

USER root

COPY config.toml-dist config.toml
COPY . .

RUN make build-auth

FROM jerson/base:1.2

LABEL maintainer="jeral17@gmail.com"

ENV BUILDER_PATH /app
ENV WORKDIR /app
WORKDIR ${WORKDIR}

COPY --from=builder ${BUILDER_PATH}/config.toml .
COPY --from=builder ${BUILDER_PATH}/auth-server .

EXPOSE 80

ENTRYPOINT ["/app/auth-server"]