# DEPLOY.md — 실배포 가이드

EC2(또는 아무 리눅스 서버) 한 대에 Docker Compose로 배포한다. main push 시 GitHub Actions가 빌드 검사 → SSH 배포까지 자동으로 한다.

```
[git push main] → [GitHub Actions]
                    ① Build Gate: go vet+build / npm ci+build  ← 깨지면 여기서 중단
                    ② SSH 접속 → git pull → docker compose up -d --build
[서버] Caddy(:80/:443) → /api/* → Go(:8080), 나머지 → SvelteKit(:3000)
```

## 1. GitHub Secrets 등록

리포 Settings → Secrets and variables → Actions → New repository secret:

| Secret | 값 |
|---|---|
| `EC2_HOST` | 서버 퍼블릭 IP (또는 도메인) |
| `EC2_USER` | SSH 유저 (보통 `ubuntu`) |
| `EC2_SSH_KEY` | SSH 개인키 전체 내용 (`-----BEGIN ... KEY-----` 포함) |

`gh` CLI가 있으면:

```bash
gh secret set EC2_HOST --body "1.2.3.4"
gh secret set EC2_USER --body "ubuntu"
gh secret set EC2_SSH_KEY < ~/.ssh/your-key.pem
```

> Secret이 없으면 배포 스텝은 자동 스킵된다 (빌드 게이트만 돌고 초록불). 등록하는 순간부터 배포 활성화.

## 2. 서버가 리포를 pull 할 수 있게 (비공개 리포)

둘 중 하나:

**방법 A — Deploy Key (권장)**

```bash
# 서버에서
ssh-keygen -t ed25519 -C "deploy@server" -f ~/.ssh/id_ed25519 -N ""
cat ~/.ssh/id_ed25519.pub
# → 출력된 공개키를 GitHub 리포 Settings > Deploy keys 에 등록 (read-only)
#   또는: gh repo deploy-key add ~/.ssh/id_ed25519.pub

# remote를 SSH 형식으로
git remote set-url origin git@github.com:<owner>/<repo>.git
```

**방법 B — PAT (Personal Access Token)**

```bash
# 서버에서 (PAT는 repo read 권한이면 충분)
git remote set-url origin https://x-access-token:<PAT>@github.com/<owner>/<repo>.git
```

## 3. 서버 최초 셋업 (1회)

```bash
# Docker 설치 (Ubuntu)
curl -fsSL https://get.docker.com | sudo sh
sudo usermod -aG docker $USER   # 재로그인 필요

# clone — Actions 배포 스크립트가 ~/<리포이름> 을 전제로 한다
cd ~
git clone <repo-url>
cd <repo-name>

# .env 준비
cp .env.example .env
# 도메인이 있으면 .env에:
#   SITE_ADDRESS=example.com
#   PUBLIC_ORIGIN=https://example.com

# 최초 기동
docker compose up -d --build
docker compose ps         # 3개(backend/frontend/caddy) Up 확인
```

## 4. HTTPS (Caddy 자동)

certbot 불필요. 조건만 맞추면 Caddy가 인증서 발급/갱신을 알아서 한다:

1. 도메인의 A 레코드 → 서버 IP
2. 서버 보안그룹에서 80, 443 인바운드 오픈
3. `.env`에 `SITE_ADDRESS=example.com`, `PUBLIC_ORIGIN=https://example.com`
4. `docker compose up -d` (재기동)

첫 접속 시 자동으로 Let's Encrypt 인증서를 받는다. 인증서는 `caddy_data` 볼륨에 저장되므로 볼륨을 지우지 말 것 (재발급 반복 시 rate limit 걸림).

도메인 없이 IP만 쓰면 `SITE_ADDRESS`를 비워두면 된다 (:80 HTTP).

## 5. 운영 명령 모음

```bash
docker compose ps                      # 상태
docker compose logs -f backend         # Go 로그
docker compose logs -f caddy           # 프록시/인증서 로그
docker compose up -d --build           # 수동 배포 (git pull 후)
docker compose restart frontend        # 개별 재시작
docker image prune -f                  # 옛 이미지 정리
```

## 6. 트러블슈팅

- **Actions build 실패** → 로컬에서 같은 검사 재현: `cd frontend && npm run build`, `cd backend && go vet ./... && go build ./...`
- **Actions deploy 실패 (SSH)** → EC2_HOST/EC2_USER/EC2_SSH_KEY 값 확인, 서버 보안그룹 22번 포트, `~/<리포이름>` 경로에 clone 돼 있는지
- **서버에서 git pull 실패** → 2번(Deploy Key/PAT) 설정 확인
- **HTTPS 안 됨** → `docker compose logs caddy`에서 인증서 발급 에러 확인. 대부분 DNS 미전파 또는 80/443 미오픈
- **502 Bad Gateway** → `docker compose ps`로 backend/frontend가 살아있는지, `docker compose logs`로 크래시 원인 확인
