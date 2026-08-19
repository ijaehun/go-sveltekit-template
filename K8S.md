# K8S.md — 쿠버네티스 로컬 학습 코스

> ⚠️ 이 폴더의 k8s 설정은 **실배포용이 아니라 학습용**이다. 실제 배포는 `DEPLOY.md`(Docker Compose)를 따른다.
> 목표: 같은 앱을 Compose와 쿠버네티스 양쪽에 배포해보면서 차이를 몸으로 느끼는 것.

## 쿠버네티스가 뭐였더라 (1분 요약)

Docker Compose는 **서버 한 대** 안에서 컨테이너를 띄우는 도구다. 쿠버네티스는 **여러 서버를 하나처럼 묶어** 컨테이너의 배치·복구·확장을 자동화하는 관리자다.

핵심 개념은 **선언형(desired state)**: "backend 파드는 항상 3개 떠 있어야 해"라고 *선언*하면, 쿠버네티스가 알아서 그 상태를 *유지*한다. 파드가 죽으면 되살리고, 부족하면 새로 띄운다. Compose가 "띄워라"(명령)라면 k8s는 "유지해라"(선언)다.

용어 최소 세트:
- **파드(Pod)**: 컨테이너를 감싼 최소 실행 단위 (대부분 파드 1개 = 컨테이너 1개)
- **Deployment**: "이 이미지로 파드 N개 유지해라"는 선언. 자가복구·롤링배포 담당
- **Service**: 파드들 앞에 붙는 고정 내부 주소 + 로드밸런서 (파드는 죽고 살면서 IP가 바뀌니까)

---

## 0. 준비: 로컬 클러스터 켜기

Docker Desktop → **Settings → Kubernetes → Enable Kubernetes** 체크 → Apply & Restart.
(몇 분 걸린다. 하단에 초록색 k8s 아이콘이 뜨면 완료)

```bash
kubectl cluster-info     # 클러스터가 응답하면 OK
kubectl get nodes        # docker-desktop 노드 1개가 Ready면 준비 끝
```

> `kubectl`은 Docker Desktop이 같이 설치해준다. 안 되면 Docker Desktop 재시작.

## 1. 이미지 빌드

쿠버네티스는 "이미지"를 배포한다. 먼저 로컬에서 빌드:

```bash
# 프로젝트 루트에서
docker build -t go-svelte-backend:v1 ./backend
docker build -t go-svelte-frontend:v1 ./frontend
```

Docker Desktop의 k8s는 로컬 docker 이미지를 그대로 쓸 수 있다 (`imagePullPolicy: IfNotPresent`). 실무에서는 이 자리에 레지스트리(ECR, Docker Hub) push/pull이 들어간다.

## 2. 배포

```bash
kubectl apply -f k8s/
```

**지금 일어난 일**: `backend.yaml`/`frontend.yaml`의 Deployment 선언을 읽고, 쿠버네티스가 "원하는 상태(각 1개 파드)"와 "현재 상태(0개)"의 차이를 스스로 메꿨다 — 파드를 만들어서.

```bash
kubectl get pods                 # STATUS가 Running이 될 때까지 잠깐 기다림
kubectl get deployments          # READY 1/1
kubectl get services             # backend(8080), frontend(3000) 내부 주소 생김
kubectl logs deploy/backend      # Go 서버 로그 확인
```

## 3. 접속 (port-forward)

Service는 클러스터 **내부** 주소라 브라우저에서 바로 못 들어간다. 터널을 뚫자:

```bash
# 터미널 2개에서 각각
kubectl port-forward svc/frontend 3000:3000
kubectl port-forward svc/backend 8080:8080
```

→ http://localhost:3000 접속. 페이지는 뜨는데 **API 데모는 실패한다**. 왜?
프론트가 `/api/...`를 호출하면 3000번(프론트 터널)으로 가는데, 거기엔 Go가 없다. Compose에서 이 교통정리를 하던 게 Caddy였고, k8s에서 같은 역할을 하는 게 **Ingress**다(4번, 옵션). 지금은 백엔드만 직접 확인해보자:

```bash
curl http://localhost:8080/api/health     # {"status":"ok"}
curl http://localhost:8080/api/hello
```

## 4. [옵션] Ingress로 제대로 라우팅

`/api → backend, 나머지 → frontend` 라우팅까지 원하면:

```bash
# ingress-nginx 컨트롤러 설치 (클러스터당 1회)
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/cloud/deploy.yaml
# 잠깐 기다린 후 (ingress-nginx 파드가 Running 되면)
kubectl apply -f k8s/ingress.yaml
```

→ http://localhost 접속하면 Compose+Caddy 때와 똑같이 전부 동작한다.

---

## 실험 3종 — 쿠버네티스의 존재 이유를 눈으로 확인

### 실험 A: 자가복구 (파드를 죽여보자)

```bash
kubectl get pods                          # backend 파드 이름 확인
kubectl delete pod <backend-파드-이름>     # 고의로 죽임
kubectl get pods                          # 바로 다시 조회
```

**결과**: 몇 초 만에 **새 파드가 자동으로 떠 있다** (이름은 바뀜). Deployment가 "1개 유지"라는 선언과 현실(0개)의 차이를 감지하고 즉시 메꾼 것. Compose였다면 컨테이너는 죽은 채로 있다 (restart 정책은 프로세스 죽음만 잡고, 이런 수준의 관리는 못 한다). **새벽에 서버 프로세스가 죽어도 사람이 안 일어나도 되는 이유.**

### 실험 B: 스케일링 (1개 → 3개)

```bash
kubectl scale deployment backend --replicas=3
kubectl get pods -w        # 파드가 3개로 늘어나는 걸 실시간으로 봄 (Ctrl+C로 종료)
```

**결과**: backend 파드 3개. Service(backend:8080)로 오는 요청은 3개 파드에 **자동으로 분배**된다(로드밸런싱). 트래픽 몰릴 때 명령 한 줄(실무에선 자동)로 확장하는 것. `k8s/backend.yaml`의 `replicas: 3`으로 바꾸고 `kubectl apply -f k8s/`를 해도 같다 — 이게 선언형 방식이다.

```bash
kubectl scale deployment backend --replicas=1    # 원상복구
```

### 실험 C: 무중단 롤링 배포

코드를 조금 바꿔보자. 예: `backend/internal/api/hello.go`의 `"Hello from Go!"`를 `"Hello from Go v2!"`로 수정 후:

```bash
docker build -t go-svelte-backend:v2 ./backend      # 새 버전 이미지
kubectl set image deployment/backend backend=go-svelte-backend:v2
kubectl rollout status deployment/backend            # 교체 과정 관찰
```

**결과**: 새 파드(v2)를 먼저 띄우고 → readinessProbe(/api/health) 통과를 확인한 뒤 → 옛 파드(v1)를 내린다. **서비스가 한순간도 안 끊긴다.** Compose의 `up -d --build`는 잠깐 끊기는 것과 대비되는 지점.

실수했다면 롤백도 한 줄:

```bash
kubectl rollout undo deployment/backend
```

---

## 5. 정리

```bash
kubectl delete -f k8s/       # 배포한 것 전부 삭제
# Docker Desktop 설정에서 Kubernetes 체크 해제하면 클러스터 자체도 꺼짐 (리소스 절약)
```

## 다음 단계 (더 배우고 싶으면)

- **ConfigMap / Secret**: 환경변수·비밀값을 k8s 리소스로 관리 (.env의 k8s 버전)
- **HPA**(Horizontal Pod Autoscaler): CPU 사용률 기반 replicas 자동 조절
- **k3s**: 진짜 서버(EC2 한 대)에 올리는 경량 쿠버네티스 — 이 매니페스트가 거의 그대로 동작한다
- **Helm**: 매니페스트 묶음의 패키지 매니저

> 감 잡는 요령: 뭐든 `kubectl get <리소스>` → `kubectl describe <리소스> <이름>` → `kubectl logs <파드>` 순서로 들여다보면 된다.
