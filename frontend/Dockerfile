# base: 의존성 설치까지 (dev에서는 이 스테이지를 target으로 사용)
FROM node:22-alpine AS base
WORKDIR /app
COPY package*.json ./
RUN npm ci

# build: SvelteKit 프로덕션 빌드 (adapter-node → build/ 디렉토리)
FROM base AS build
COPY . .
RUN npm run build

# 최종 이미지: 빌드 결과만 담아 node로 실행.
# 런타임 dependencies가 없으므로(전부 devDependencies) node_modules 불필요.
# 런타임 dep을 추가하게 되면: npm ci --omit=dev 한 node_modules 복사 필요.
FROM node:22-alpine
WORKDIR /app
ENV NODE_ENV=production
COPY --from=build /app/build ./build
COPY --from=build /app/package.json ./
EXPOSE 3000
ENV PORT=3000
CMD ["node", "build"]
