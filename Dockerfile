ARG FROM_IMAGE=amd64/golang:1.20-alpine
ARG PROD_IMAGE=scratch

FROM ${FROM_IMAGE} as base
ENV WORKDIR=/build

WORKDIR ${WORKDIR}

COPY ./go.mod ./go.sum ./
RUN go mod download && go mod verify

COPY ./ ./

FROM base as test
ARG USER_UID=user

# Add user to not run tests as root and avoid permissions errors
RUN addgroup ${USER_UID} && adduser -D -G ${USER_UID} ${USER_UID} && chown -R ${USER_UID}:${USER_UID} $WORKDIR
RUN apk add make

USER ${USER_UID}

FROM base as builder
ARG ARCH=amd64
ARG APP_VERSION
RUN GOARCH=${ARCH} go build -ldflags="-w -s ${APP_VERSION:+-X github.com/franciscolkdo/guntar/cmd.Version=${APP_VERSION}}" -o /guntar

FROM ${PROD_IMAGE} as prod

COPY --from=builder /guntar /guntar

ARG APP_VERSION
ARG COMMIT_ID

ENV APP_VERSION=${APP_VERSION}
ENV COMMIT_ID=${COMMIT_ID}
# Enable colors in container
ENV COLORTERM=truecolor

ENTRYPOINT [ "/guntar" ]