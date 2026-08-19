# Production build: 정적 바이너리 하나로 컴파일 → 초경량 이미지 (~15MB)
FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod ./
# go.sum이 생기면(외부 의존성 추가 시) 아래 주석 해제
# COPY go.sum ./
# RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /server ./cmd/server

FROM alpine:3.21
RUN adduser -D -u 1000 app
USER app
COPY --from=build /server /server
EXPOSE 8080
# busybox wget으로 헬스체크 (compose/k8s에서 사용)
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s \
  CMD wget -qO- http://localhost:8080/api/health || exit 1
ENTRYPOINT ["/server"]
