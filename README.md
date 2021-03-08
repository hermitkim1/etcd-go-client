# etcd Go client(v3) examples
etcd 설치부터 Go client 테스트까지 몇가지 어려운 점이 있었기에 간략히 정리하고자 합니다.
- etcd Official 홈페이지 및 GitHub의 정보 동기화가 되지 않은 부분이 있어 etcd 설치, 예제 실행에 어려움을 겪음
- coreos 기반으로 개발해 오다가 ectd-io로 이전 됨에 따라 예전 정보와 최근 정보의 구분이 필요함
- etcd Go library (client v3) 예제는 현재 적은 것으로 파악됨


## etcd 설치
- Windows10에 WSL2를 활성화
- Windows Terminal를 통해 Ubuntu 20.04 활용
- Single etcd 구성(Localhost)

비고: Single Point Of Failre (SPOF) 회피를 위해 추후 etcd cluster를 구축 예정

## Go client examples
- example1-basic.go: 기본 예제
- example2-watch.go: Watch 기능을 활용한 예제
- example3-concurrency-by-distributed-lock.go: Adder와 Subtractor가 덧셈과 뺄셈을 반복하는 Racing 예제로 Distributed lock (분산 락)을 활용하여 동시성 보장 예제

향후 etcd cluster 구성 및 활용에 대한 내용 추가 예정
