# etcd Go client(v3) examples
etcd 설치부터 Go client 테스트까지 몇가지 어려운 점이 있었기에 간략히 정리하고자 합니다.
- etcd Official 홈페이지 및 GitHub의 정보 동기화가 되지 않은 부분이 있어 etcd 설치, 예제 실행에 어려움을 겪음
- coreos 기반으로 개발해 오다가 ectd-io로 이전 됨에 따라 "예전 정보"와 "최근 정보"의 구분이 필요함
- etcd Go library (client v3) 예제는 현재 적은 것으로 파악됨


## etcd standalone cluster 설치 가이드 (on Windows 10)
저는 Windows 10에 환경구성 및 스터디를 진행하였습니다.
- Windows 10에 WSL2 설치 및 사용 ([참고: WSL2 설치 및 사용 방법](https://www.44bits.io/ko/post/wsl2-install-and-basic-usage))
- Windows Terminal 실행 후 Ubuntu 탭 Open
![image](https://user-images.githubusercontent.com/7975459/111126003-9167b600-85b5-11eb-8fe2-4cc3096bb553.png)

  - Home 디렉토리로 이동
  ```bash
  cd ~
  ```

- Golang 설치 ([참고](https://github.com/cloud-barista/cb-coffeehouse/tree/master/scripts/golang))
```bash
wget https://raw.githubusercontent.com/cloud-barista/cb-coffeehouse/master/scripts/golang/go-installation.sh
source go-installation.sh
```

- Download and build etcd ([참고](https://etcd.io/docs/v3.4.0/dl-build/))
  - 방법 1
  ```bash
  git clone https://github.com/etcd-io/etcd.git
  cd etcd
  ./build
  ```
  - 방법 2
  ```bash
  go get -v go.etcd.io/etcd
  go get -v go.etcd.io/etcd/etcdctlcd etcd
  ```


## etcd Go client examples running on a standalone cluster
- example1-basic.go: 기본 예제
- example2-watch.go: Watch 기능을 활용한 예제
- example3-concurrency-by-distributed-lock.go: Adder와 Subtractor가 덧셈과 뺄셈을 반복하는 Race condition 예제로 Distributed lock (분산 락)을 활용하여 동시성 보장을 확인하는 예제
- example4-atomic-compare-and-swap.go: Atomic Compare-And-Swap (CAS) 기능을 활용하여 동시성 보장을 확인하는 예제

## etcd multi-member cluster 구성 가이드 (On multi-cloud)
Single Point Of Failre (SPOF) 회피를 위해 etcd cluster를 구성하였습니다.

[TBD]
구성 방법은 추후 추가 하겠습니다 😅

## etcd Go client examples running on a multi-member cluster
- example4-host1-watcher.go: etcd cluster의 "phoo" key값을 Watch하는 예제
- example4-host2-adder.go: etcd cluster의 "phoo" key값을 1씩 더하는 예제
- example4-host2-adder-with-lock.go: etcd cluster의 "phoo" key값을 1씩 더하는 예제로 Distributed lock이 적용되어 있음
- example4-host2-adder-with-case.go: etcd cluster의 "phoo" key값을 1씩 더하는 예제로 Compare-and-swap(CAS)이 적용되어 있음
- example4-host3-subtractor.go: etcd cluster의 "phoo" key값을 1씩 빼는 예제
- example4-host3-subtractor-with-lock.go: etcd cluster의 "phoo" key값을 1씩 빼는 예제로 Distributed lock이 적용되어 있음
- example4-host3-subtractor-with-cas.go: etcd cluster의 "phoo" key값을 1씩 빼는 예제로 Compare-and-swap(CAS)이 적용되어 있음
