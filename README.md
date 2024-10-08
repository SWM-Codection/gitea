<p align="center">
    <img height="360px" src="https://github.com/user-attachments/assets/c0e1a377-b180-40af-b02d-405389f7c296"/> 
</p>

<br/> 

[![](https://github.com/go-gitea/gitea/actions/workflows/release-nightly.yml/badge.svg?branch=main)](https://github.com/go-gitea/gitea/actions/workflows/release-nightly.yml?query=branch%3Amain "Release Nightly")
[![](https://img.shields.io/discord/322538954119184384.svg?logo=discord&logoColor=white&label=Discord&color=5865F2)](https://discord.gg/Gitea "Join the Discord chat at https://discord.gg/Gitea")
[![](https://goreportcard.com/badge/code.gitea.io/gitea)](https://goreportcard.com/report/code.gitea.io/gitea "Go Report Card")
[![](https://pkg.go.dev/badge/code.gitea.io/gitea?status.svg)](https://pkg.go.dev/code.gitea.io/gitea "GoDoc")
[![](https://img.shields.io/github/release/go-gitea/gitea.svg)](https://github.com/go-gitea/gitea/releases/latest "GitHub release")
[![](https://www.codetriage.com/go-gitea/gitea/badges/users.svg)](https://www.codetriage.com/go-gitea/gitea "Help Contribute to Open Source")
[![](https://opencollective.com/gitea/tiers/backers/badge.svg?label=backers&color=brightgreen)](https://opencollective.com/gitea "Become a backer/sponsor of gitea")
[![](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT "License: MIT")
[![Contribute with Gitpod](https://img.shields.io/badge/Contribute%20with-Gitpod-908a85?logo=gitpod)](https://gitpod.io/#https://github.com/go-gitea/gitea)
[![](https://badges.crowdin.net/gitea/localized.svg)](https://crowdin.com/project/gitea "Crowdin")
[![](https://badgen.net/https/api.tickgit.com/badgen/github.com/go-gitea/gitea/main)](https://www.tickgit.com/browse?repo=github.com/go-gitea/gitea&branch=main "TODOs")


## 목차 
1. [프로젝트 개요](##프로젝트-개요) 
    - [프로젝트 소개](###프로젝트-소개)
    - [시스템 구성도](###시스템-구성도)
    - [주요 기능](###주요-기능)
    - [개발 환경](###개발-환경)

2. [개발 결과물](##개발-결과물) 
    - [백엔드 아키텍처](###백엔드-아키텍처)
    - [API 서버 개발 및 API 문서화/테스트](###API-서버-개발-및-API-문서화/테스트)
    - [API 서버 리팩토링](###API-서버-리팩토링)
    - [API 서버 CI/CD 파이프라인](###API-서버-CI/CD-파이프라인)
    
3. [수행 방법 및 프로젝트 관리](##수행-방법-및-프로젝트-관리) 
    - [개발 프로세스](###개발-프로세스)
    - [KPT 회고 및 데일리 스크럼](###KPT-회고-및-데일리-스크럼)
    - [형상 관리 프로세스](###형상-관리-프로세스)
    - [오픈 소스 컨트리뷰트](###오픈-소스-컨트리뷰트)




## 프로젝트 개요 

### 프로젝트 소개
![](https://github.com/user-attachments/assets/909346cc-22d6-46a0-b4e6-a1e0c7730c9c)

Codection은 코드리뷰 효율성의 향상을 돕는 소프트웨어로서, 오픈 소스 프로젝트인 Gitea를 기반으로 동작하고 있습니다. 



### 시스템 구성도
<center>

![system-archiecture](https://github.com/user-attachments/assets/a31c67a3-3fa2-4eb6-abe6-db4d9f9b889e)

</center>

<center>

![application-architecture](https://github.com/user-attachments/assets/bef5d209-5538-4d96-b74f-d11123c1e04c)

</center>

### 주요 기능

<center>

![main-feature](https://github.com/user-attachments/assets/a6c3bbd0-80ee-4f1a-bb73-8821609713fe)

</center>


### 개발 환경

- Frontend: VueJS, Go template, Javascript
- Backend: Spring Boot, go-chi, Github Actions

## 개발 결과물 

### 백엔드 아키텍처
WIP

### API 서버 개발 및 API 문서화/테스트

<center>

![api-documentation](https://github.com/user-attachments/assets/7acffc5a-e7bf-4169-ae4b-a52e08f8f898)

</center>


OpenAPI 3.0 Spec을 준수하는 Swagger를 사용하여 API 문서화를 진행하였습니다.

### API 서버 리팩토링
WIP

### API 서버 CI/CD 파이프라인

<center>

![ci](https://github.com/user-attachments/assets/70019a62-5669-4cdc-892e-694c1a8240e2)

</center>

Github Actions 를 이용하여, CICD 파이프라인을 구축하였습니다.  
또한 Branch Protection Rule 을 통해 반드시 CICD 파이프라인을 통과하고, 코드 리뷰가 이루어진 Pull Request에 대해서만 Merge 를 허용하도록 정책을 설정하여, 보다 안전하고 효과적인 통핣 및 배포 프로세스를 구축할 수 있었습니다. 

## 수행 방법 및 프로젝트 관리

### 개발 프로세스 
<center>

![jira-kanban](https://github.com/user-attachments/assets/5397cd52-7082-47a6-9e5b-e3db872733df)

</center>

저희 팀에서는 Jira 의 칸반 기능을 이용하여, 팀원 간 역할 분담을 명확히 하여, 효율적인 협업 프로세스를 이룰 수 있었습니다. 
또한 Jira 에서 발급한 티켓을 기반으로, 풀 리퀘스트의 연동을 진행하였습니다.

### KPT 회고 및 데일리 스크럼

<p align="center">

<img src="https://github.com/user-attachments/assets/4e4072a5-f340-4ff0-b80e-e203e7005238"  height="1200"/>


</p>

하루의 시작을 데일리 스크럼을 통해 작업 일정을 공유하였고 주간 KPT 회고를 진행하여 작업 효율성 향상을 위한 방법론을 찾아갔습니다.

### 형상 관리 프로세스

<center>

![gitflow](https://github.com/user-attachments/assets/f89b2815-0add-4011-bc02-a7362a7f2438)  
git-flow 전략

</center>

저희 프로젝트 Codection에서는 효율적인 협업 방식을 위해 Git Flow 전략을 적극적으로 도입하였습니다. Git Flow 전략을 통해 팀 전체의 작업 흐름을 명확히 하고, 릴리스와 빠른 버그 수정을 이룰 수 있었습니다. 이를 통해 코드 관리의 복잡도를 줄이고, 협업 속도와 품질을 높일 수 있었습니다. 

### 오픈 소스 컨트리뷰트
<center>

![issue-32080](https://github.com/user-attachments/assets/18ef0089-1f16-4e51-aa4e-b88776e17bfc)

</center>

저희 팀에서는 프로젝트를 진행하던 중 기반이 되는 gitea에서의 버그를 발견하고 이를 제보 및 해결 방안을 제시하여, 
gitea 1.23 마일스톤에 등록하게끔 하는 성과를 이루었습니다. 

<center>

![pr-32081](https://github.com/user-attachments/assets/cb5a8c8a-6880-4640-bc24-6134beb8d081)

</center>

또한 gitea 메인테이너와 효율적인 의사소통을 통해, 기존의 제안을 개선 및 보강하였습니다. 


## 라이센스
이 프로젝트는 MIT 라이센스를 따릅니다.  
전체 라이센스 원문을 보기 위해서는 [라이센스 파일을](https://github.com/swm-codection/gitea/blob/main/LICENSE) 참고하세요
