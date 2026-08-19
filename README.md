# Go + SvelteKit Template

빠른 웹 개발 시작용 풀스택 템플릿 — 새 스택 학습 겸용. `fastapi-nextjs-template`의 구조/배포 방식을 계승하고 기술 스택만 전부 교체한 버전.

## 🎯 사용법

```bash
# 1. 템플릿 복사 (또는 GitHub에서 새 repo로 fork)
cp -r go-sveltekit-template/ <your-new-project>/
cd <your-new-project>/

# 2. .env 파일 준비
cp .env.example .env
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env

# 3. git 초기화
rm -rf .git
git init
git add .
git commit -m "init: from go-sveltekit-template"

# 4. 개발 시작 (핫리로드: Go는 air, Svelte는 vite)
docker compose -f docker-compose.yml -f docker-compose.dev.yml up
# → http://localhost:5173
```

## 🛠️ 기술 스택

### Backend — Go
- **Go 1.24 표준 라이브러리만** (외부 프레임워크 없음, 의존성 제로)
- Go 1.22+ 내장 패턴 라우팅 (`"GET /api/health"`, `PathValue`)
- graceful shutdown, 미들웨어 체인(Logger/Recover/CORS), SSE 스트리밍 데모
- 프로덕션 이미지 ~15MB (정적 바이너리 + alpine)

### Frontend — SvelteKit
- **Svelte 5** (runes: `$state`, `$derived`, `$effect`) + TypeScript
- **Tailwind CSS v4** (설정 파일 없이 CSS import 방식)
- **adapter-node** — SSR 유지 (Next.js와 같은 능력치)
- 데모 페이지: API 호출 + 경로 변수 + SSE 실시간 시계

### Infrastructure
- **Caddy** — 리버스 프록시. 도메인만 주면 **HTTPS 자동** (certbot 불필요)
- **Docker Compose** — prod/dev 이원화
- **GitHub Actions** — 빌드 게이트 통과 시 EC2 SSH 자동 배포 (`DEPLOY.md`)
- **Kubernetes 학습 코스** — 로컬 클러스터로 자가복구/스케일링/롤링배포 실험 (`K8S.md`)

## 📁 구조

```
project/
├── README.md, CLAUDE.md, DEPLOY.md, K8S.md, LICENSE
├── .env.example
├── docker-compose.yml          # 프로덕션 (backend + frontend + caddy)
├── docker-compose.dev.yml      # 개발 override (핫리로드)
│
├── backend/                    # Go
│   ├── Dockerfile              # 프로덕션 (멀티스테이지 → ~15MB)
│   ├── Dockerfile.dev          # air 핫리로드
│   ├── go.mod                  # 의존성 제로
│   └── cmd/server/main.go      # 엔트리 (graceful shutdown)
│   └── internal/
│       ├── config/             # 환경변수 로드
│       ├── middleware/         # Logger, Recover, CORS (http.Handler 체인)
│       └── api/                # 라우터 + feature별 핸들러
│
├── frontend/                   # SvelteKit
│   ├── Dockerfile
│   ├── svelte.config.js        # adapter-node
│   ├── vite.config.ts          # dev proxy: /api → backend
│   └── src/
│       ├── routes/             # 페이지 (+page.svelte, +layout.svelte)
│       └── lib/
│           ├── api.ts          # fetch 래퍼 (항상 상대경로 /api)
│           └── components/     # Header, Footer
│
├── caddy/Caddyfile             # /api/* → backend, 나머지 → frontend
├── k8s/                        # 쿠버네티스 학습용 매니페스트 (K8S.md)
├── scripts/deploy.sh           # 수동 배포 (tar+scp)
└── .github/workflows/deploy.yml  # CI: 빌드 게이트 → SSH 배포
```

## 🚀 개발

### 로컬 직접 실행 (가장 빠른 반복)
```bash
# Backend (Go 설치 필요: https://go.dev/dl/)
cd backend
go run ./cmd/server            # :8080

# Frontend
cd frontend
npm install
npm run dev                    # :5173, /api는 proxy로 백엔드에 전달
```

### Docker (전체 스택)
```bash
# Dev (핫리로드)
docker compose -f docker-compose.yml -f docker-compose.dev.yml up

# Prod (caddy + 최적화 빌드)
docker compose up -d --build   # → http://localhost
```

## 🌐 배포

`DEPLOY.md` 참고. 요약: GitHub Secrets(EC2_HOST/EC2_USER/EC2_SSH_KEY) 등록 → main push → 빌드 검사 통과 시 자동 배포. 도메인이 있으면 `.env`에 `SITE_ADDRESS=example.com` 한 줄로 HTTPS 자동 적용.

## ☸️ 쿠버네티스 학습

`K8S.md` 참고. Docker Desktop 내장 클러스터에서 같은 이미지를 배포해보고, 파드 자가복구·스케일링·무중단 롤링 배포를 직접 실험하는 단계별 코스. (실배포는 Compose — k8s는 학습용)

## 💡 SPA로 바꾸고 싶다면 (Node 서버 제거)

SEO가 필요 없는 대시보드/관리도구라면 `@sveltejs/adapter-static`으로 교체하면 빌드 결과가 정적 파일이 되어 Caddy가 직접 서빙할 수 있다 (frontend 컨테이너 제거 가능). `svelte.config.js`에서 adapter를 바꾸고 `src/routes/+layout.ts`에 `export const prerender = true; export const ssr = false;` 추가.
