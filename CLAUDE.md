# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

(Edit this for your specific project)

- **What**: Go + SvelteKit fullstack web app
- **Domain**: 
- **Status**: 

## Stack

- **Backend**: Go 1.24, 표준 라이브러리만 (의존성 제로, go.sum 없음)
- **Frontend**: SvelteKit + Svelte 5 (runes) + TypeScript + Tailwind CSS v4
- **Proxy**: Caddy (자동 HTTPS)
- **Infra**: Docker Compose (prod/dev 이원화) + GitHub Actions
- **k8s/**: 로컬 학습용 (실배포 아님, `K8S.md`)

## Architecture

```
project/
├── backend/            # Go (표준 라이브러리)
│   ├── cmd/server/     # main.go — 엔트리, graceful shutdown
│   └── internal/
│       ├── config/     # 환경변수 로드
│       ├── middleware/ # Logger/Recover/CORS (http.Handler 체인)
│       └── api/        # router.go + feature별 핸들러 파일
├── frontend/           # SvelteKit
│   └── src/
│       ├── routes/     # 파일 기반 라우팅 (+page.svelte, +layout.svelte)
│       └── lib/
│           ├── api.ts        # fetch 래퍼 (항상 상대경로 /api 사용)
│           └── components/   # 재사용 컴포넌트
├── caddy/Caddyfile     # /api/* → backend:8080, 나머지 → frontend:3000
├── k8s/                # 쿠버네티스 학습용 매니페스트
├── docker-compose.yml      # 프로덕션
└── docker-compose.dev.yml  # 개발 override (핫리로드)
```

## Development Commands

### Local (recommended for fast iteration)
```bash
# Backend (Go 필요)
cd backend
go run ./cmd/server        # :8080

# Frontend
cd frontend
npm install
npm run dev                # :5173 (vite proxy가 /api → localhost:8080)
```

### Docker (full stack)
```bash
# Dev (핫리로드: Go는 air, Svelte는 vite) → http://localhost:5173
docker compose -f docker-compose.yml -f docker-compose.dev.yml up

# Prod (caddy + 최적화 빌드) → http://localhost
docker compose up -d --build
```

## Conventions

### Backend (Go)
- 라우트는 `internal/api/<feature>.go`에 핸들러 작성 → `internal/api/router.go`에 등록
- Go 1.22+ 패턴 라우팅 사용: `mux.HandleFunc("GET /api/xxx/{id}", ...)`, 경로 변수는 `r.PathValue("id")`
- JSON 응답은 `internal/api/respond.go`의 `writeJSON`/`writeError` 헬퍼 사용
- 설정은 `internal/config`에서 환경변수로만 (파일 로드 없음)
- **의존성 제로 정책**: 표준 라이브러리로 해결이 원칙. 외부 패키지가 꼭 필요하면 `go get` 후 go.sum 커밋 + backend/Dockerfile의 `COPY go.sum` / `go mod download` 주석 해제

### Frontend (SvelteKit)
- 페이지는 `src/routes/<경로>/+page.svelte` (파일 기반 라우팅)
- **Svelte 5 runes 문법 사용**: `$state`, `$derived`, `$effect`, `$props`. 옛 문법(`export let`, `$:`) 쓰지 말 것
- API 호출은 `$lib/api.ts`의 `api()` 래퍼로, **항상 상대경로** `/api/...` (dev는 vite proxy, prod는 Caddy가 라우팅)
- 재사용 컴포넌트는 `src/lib/components/`, `$lib` alias로 import
- Tailwind로 스타일링 (v4 — tailwind.config 파일 없음, `app.css`의 `@import 'tailwindcss'`가 전부)

### Environment variables
- `.env.example` always tracked (template) / `.env`는 gitignored
- Backend: `internal/config/config.go`가 환경변수 직접 읽음 (compose가 주입)
- Frontend: 브라우저 노출 변수는 `PUBLIC_` 접두사 + `$env/static/public`

### Production deployment
- **자동 배포**: `.github/workflows/deploy.yml` — main push 시 빌드 검사 후 서버에 SSH로 배포. 셋업은 `DEPLOY.md` 참고.
- **수동 배포**: `scripts/deploy.sh` (tar+scp) 또는 서버에서 `git pull && docker compose up -d --build`
- 프로덕션은 `docker-compose.yml` 사용 (`.dev.yml` 아님). 프론트는 **adapter-node 빌드**(`node build`)로 실행 — 절대 `npm run dev`로 운영하지 말 것
- Caddy routes `/api/*` → backend, rest → frontend. 도메인 있으면 `.env`에 `SITE_ADDRESS=example.com` → HTTPS 자동

## Notes for Claude Code

이 리포는 `go-sveltekit-template` 기반이다. 작업 시 아래 가이드라인을 따른다.

### 일반
- Don't commit secrets — `.env*` files (except `.env.example`) are gitignored. 빌드 산출물(`frontend/build/`, `.svelte-kit/`, `backend/tmp/`, `node_modules/`)도 커밋 금지.
- 새 feature 추가 시: 백엔드 핸들러(`internal/api/<feature>.go` + router 등록)와 프론트 호출부(`$lib/api.ts` 재사용) 양쪽을 함께 만든다

### 배포가 깨지지 않게 (중요)
- **`svelte.config.js`의 adapter-node를 함부로 바꾸지 말 것** — 프로덕션 Dockerfile(`node build`)이 이걸 전제로 동작한다. adapter-static 전환은 Dockerfile/compose/Caddyfile 수정까지 세트로 해야 한다 (README 참고).
- **프로덕션 `docker-compose.yml`에 소스 bind-mount(`./frontend:/app` 등)를 추가하지 말 것** — 빌드된 이미지를 덮어써서 망가진다. 핫리로드 마운트는 `docker-compose.dev.yml`에만.
- 프론트 Dockerfile을 `npm run dev`로 바꾸지 말 것. 빌드 타입 에러가 나면 dev로 우회하지 말고 **타입 에러 자체를 고친다**.
- backend에 외부 패키지를 추가했다면 go.sum 커밋 + Dockerfile의 go.sum 관련 주석 해제를 잊지 말 것 (안 하면 도커 빌드 깨짐).
- `k8s/`는 학습용이다. 실배포 파이프라인(Actions/compose)과 무관하므로 배포 이슈 디버깅에 끌어들이지 말 것.

### 변경 후 셀프 체크
- 프론트 변경 시: push 전 `cd frontend && npm run build` 통과 확인 (CI 빌드 게이트와 동일한 검사)
- 백엔드 변경 시: `cd backend && go vet ./... && go build ./...` (로컬에 Go 없으면 `docker build ./backend`로 대체)

### CI/CD
- 배포 파이프라인·서버 인증(Deploy Key / PAT)·Secrets 설정은 `DEPLOY.md`에 정리돼 있다. 배포 관련 질문/작업은 거기를 먼저 참고한다.

### CI/CD 셋업 체크리스트 (이 리포에 자동배포 붙일 때 따를 순서)

먼저 사용자에게 받을 값: **서버 퍼블릭 IP / SSH 키(경로 또는 .pem 내용) / 인증 방식(Deploy Key 또는 PAT)**.
(없으면 진행 못 하는 부분이므로 먼저 물어본다. 명령 상세는 `DEPLOY.md`.)

1. **리포 파일 정비** (Claude 단독 가능)
   - `.github/workflows/deploy.yml` 존재 확인 (없으면 템플릿에서 복사)
   - `frontend/svelte.config.js`가 adapter-node인지
   - 프로덕션 `docker-compose.yml`에 소스 bind-mount 없는지
   - 빌드 검증: `cd frontend && npm run build`, `cd backend && go build ./...` (또는 docker build)
2. **GitHub Secrets 등록** — `gh` 로그인돼 있으면 `gh secret set`으로 직접: `EC2_HOST` / `EC2_USER`(ubuntu) / `EC2_SSH_KEY`(개인키 전체)
3. **서버 인증** (사용자가 서버에서 실행하거나, 권한 주면 Claude가 SSH로)
   - Deploy Key: 서버에서 키 생성 → 공개키를 GitHub deploy key로 등록(읽기전용) → remote를 SSH로 변경
   - 또는 PAT: `git remote set-url origin https://x-access-token:<PAT>@github.com/<owner>/<repo>.git`
4. **서버 최초 준비**: Docker 설치 → `~/<repo>`로 clone → `.env` 채우기 (도메인 있으면 `SITE_ADDRESS`) → `docker compose up -d --build` 1회
5. **확인**: push(또는 Actions 수동 실행) → build → deploy 초록불

> 주의: 2·3·4는 서버 IP/키/서버 접근이 있어야 하는 외부 작업이라 "리포만 보고 100% 자동"은 안 된다. 1번(리포 파일)까지는 자동으로 끝내고, 나머지는 사용자 입력을 받아 진행한다.
