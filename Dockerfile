FROM public.ecr.aws/docker/library/golang:1.23.3-alpine AS build-stage

WORKDIR /src

RUN mkdir -p /app/bin && \
    apk add upx

ADD go.mod go.sum ./

RUN go mod download

ADD . .

RUN \
    go build -ldflags='-w -s' -o /app/bin/service ./cmd/service && \
    upx /app/bin/service

FROM public.ecr.aws/amazonlinux/amazonlinux:2023-minimal AS base-image

WORKDIR /app

COPY --from=build-stage /app /app

ENV PORT=8000

EXPOSE 8000

CMD [ "/app/bin/service" ]
