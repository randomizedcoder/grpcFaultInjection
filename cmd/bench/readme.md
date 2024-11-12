
Quick benchmark to compare atomic counter verse mutex

TLDR is use atomic!!

```
[nix-shell:~/Downloads/grpcFaultInjection/cmd/bench]$ go test -bench=.
goos: linux
goarch: amd64
pkg: randomizedcoder/grpcFaultInjection/cmd/bench
cpu: Intel(R) Core(TM) i9-10885H CPU @ 2.40GHz
BenchmarkAtomicUint64-16    	68541187	       14.76 ns/op   <---- faster
BenchmarkMutexUint64-16     	19613998	       59.95 ns/op
PASS
ok  	randomizedcoder/grpcFaultInjection/cmd/bench	2.271s
```