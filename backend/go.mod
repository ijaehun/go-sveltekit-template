module app

go 1.24

// 의존성 제로 정책: Go 1.22+ 표준 라이브러리(net/http 라우팅 패턴)만 사용한다.
// 외부 패키지가 필요해지면 `go get <pkg>` 후 go.sum도 함께 커밋할 것.
