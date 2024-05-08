# sudo docker container ls -a (모든 도커 컨테이너 확인)
# docker run -it gotodo /bin/bash (도커를 실행하고 내부 shell 접근)

# 배포용 컨테이너에 포함시킬 바이너리를 생성하는 컨테이너

# Go 언어의 특정 버전인 1.18.2-bullseye를 사용하는 Docker 이미지를 기반 이미지로 사용 
# deploy-build라는 이름으로 명명
FROM golang:1.22-bullseye as deploy-builder

# 컨테이너 내에서 작업 디렉토리를 /app으로 설정
# 이후 모든 명령이 실해애될 기본 경로.
WORKDIR /app

# 호스트의 go.mod, go.sum 파일을 컨테이너의 현재 작업 디렉토리로 복사.
COPY go.mod go.sum ./
# Go 모듈에 정의된 종속성을 다운로드
RUN go mod download

# 현재 호스트 디렉토리의 모든 파일 및 폴더를 컨테이너의 현재 작업 디렉토리로 복사.
COPY . .
# 프로젝트 빌드 수행.
# -trimpath: 실행 가능한 파일에 컴파일된 파일 경로를 제거하여 이식성을 향상.
# -ldflags "-w -s": 실행 가능한 파일에 링크되는 정보를 최소화하여 실행 파일의 크기를 줄입니다.
#   -w: 모든 디버그 정보 제거
#   -s: 기호 테이블과 디버그 정보 제거.
# -o app: 출력 파일의 이름을 app으로 설정.
RUN go build -trimpath -buildvcs=false -ldflags "-w -s" -o app

# -----------------------------------------------

# 배포용 컨테이너
FROM debian:bullseye-slim as deploy

RUN apt-get update

# deploy-builder로 부터 /app/app(실행 파일)을 현재 디렉토리로 이동.
COPY --from=deploy-builder /app/app .

# 이 명령은 컨테이너가 시잘될 때 실행할 명령을 지정.
# 현재 디렉토리에 있는 app 실행 파일을 실행.
CMD ["./app"]

# -----------------------------------------------

# 로컬 개발 환경에서 사용하는 자동 새로고침 환경
FROM golang:1.22 as dev
WORKDIR /app
RUN go install github.com/cosmtrek/air@latest
CMD ["air"]

