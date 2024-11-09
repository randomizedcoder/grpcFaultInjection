# grpcFaultInjection

GRPC fault injection library using interceptors

This is a small library with client + server GRPC intercepts

The client injects metadata (http2 headers) to request that the server does
fault injection

This library is designed to allow the client to control the fault injection,
and is generally designed to allow testing of error handling code

Hopefully this library is easy to integrate into an existing GRPC eco-system,
and will provide value for failure mode testing

## Overview Diagram

<img src="./docs/Screenshot from 2024-11-08 18-49-32.png" alt="xtcp_sampling diagram" width="100%" height="100%"/>


## Example implementation
Example implementations are in the /cmd/client and /cmd/server directory

**Client**

https://github.com/randomizedcoder/grpcFaultInjection/blob/main/cmd/client/client.go

**Server**

https://github.com/randomizedcoder/grpcFaultInjection/blob/main/cmd/server/server.go

## Client usage

The client needs a configuration as follows:
```
type UnaryClientInterceptorConfig struct {
	ClientFaultPercent int
	ServerFaultPercent int
	ServerFaultCodes   string
}
```

### ClientFaultPercent
The client can be configured to randomly trigger the fault headers to be injected.

Examples

| ClientFaultPercent | Description                                                  |
| ------------------ | ------------------------------------------------------------ |
| 10                 | 10% of the time the metadata(headers) are injected           |
| 100                | 100% of the time the metadata(headers) are injected = Always |


### ServerFaultPercent
The client can make requests with "faultpercent" and "faultcodes" metadata(headers)

The configuration ServerFaultPercent injects the "faultpercent" header which is
passed to the server.

Examples

| ServerFaultPercent | Description                                                     |
| ------------------ | --------------------------------------------------------------- |
| 50                 | 50% there is a 50% chance that the server will return a fault   |
| 90                 | 90% chance the server will always return a fault                |
| 100                | 100% of the time the server will always return a fault = Always |

The client adds the "faultpercent" header which is passed to the server:

Examples

| "faultpercent"     | Description                                                     |
| ------------------ | --------------------------------------------------------------- |
| 50                 | 50% there is a 50% chance that the server will return a fault   |
| 90                 | 90% chance the server will always return a fault                |
| 100                | 100% of the time the server will always return a fault = Always |

### ServerFaultCodes

The configuration ServerFaultCodes injects the "faultcodes" header which is
passed to the server.

If "faultcodes" is NOT supplied, any random valid GRPC status code, except zero (0), is returned

Examples

| "faultcodes  "     | Description                                                         |
| ------------------ | ------------------------------------------------------------------- |
| 14                 | If the server injects the fault, the only return status is 14       |
| 10,12,14           | If the server injects the fault, possible status codes are 10,12,14 |
| <not set >         | If the server injects the fault, codes 1-16 are possible            |


Possible failcodes are:
https://github.com/grpc/grpc/blob/master/doc/statuscodes.md

## Config Matrix

Pleae keep in mind the ClientFaultPercent and ServerFaultPercent result in fault
injection probabilities like the following:

| ClientFaultPercent | ServerFaultPercent | Probability | Comment                                  |
|--------------------|--------------------|-------------|------------------------------------------|
| 100                | 0                  | 0%          | Doesn't do much                          |
| 0                  | 100                | 0%          |                                          |
|                    |                    |             |                                          |
| 100                | 10                 | 10%         | Recommended to use 100% on one end       |
| 100                | 50                 | 50%         |                                          |
| 10                 | 100                | 10%         |                                          |
| 50                 | 100                | 50%         |                                          |
|                    |                    |             |                                          |
| 10                 | 10                 | 1%          | Tricky to reason about                   |
| 50                 | 50                 | 25%         |                                          |
| 100                | 50                 | 50%         |                                          |
|                    |                    |             |                                          |
| 100                | 100                | 100%        | Unlikely to be successful!               |

( test_test.go tries to follow this table )

## Examples

### ExampleA

This is an example that will always return a GRPC error status, the status code will be random.

| Variable           | Value             | Description                                                       |
|--------------------|-------------------|-------------------------------------------------------------------|
| clientfaultpercent | 100               | The client always inserts the headers                             |
| faultpercent       | 100               | The client header instructs the server to always inject the fault |
| failcodes          | < not specified > | Random response code between 1-16 inclusive                       |
| Command            |                   | ./client --loops 5 -clientfaultpercent 100 -faultpercent 100      |


```
[das@t:~/Downloads/grpcFaultInjection/cmd/client]$ ./client --loops 5 -clientfaultpercent 100 -faultpercent 100
2024/11/08 11:29:39.002621 request success:0 fault:1
2024/11/08 11:29:39 i:0 UnaryEcho error: rpc error: code = Unavailable desc = intercept fault code:14 rp:53 success:54 fault:57
2024/11/08 11:29:39 i:0 UnaryEcho reply: <nil>
2024/11/08 11:29:39.008069 request success:0 fault:2
2024/11/08 11:29:39 i:1 UnaryEcho error: rpc error: code = ResourceExhausted desc = intercept fault code:8 rp:38 success:54 fault:58
2024/11/08 11:29:39 i:1 UnaryEcho reply: <nil>
2024/11/08 11:29:39.008260 request success:0 fault:3
2024/11/08 11:29:39 i:2 UnaryEcho error: rpc error: code = DeadlineExceeded desc = intercept fault code:4 rp:70 success:54 fault:59
2024/11/08 11:29:39 i:2 UnaryEcho reply: <nil>
2024/11/08 11:29:39.008447 request success:0 fault:4
2024/11/08 11:29:39 i:3 UnaryEcho error: rpc error: code = DeadlineExceeded desc = intercept fault code:4 rp:45 success:54 fault:60
2024/11/08 11:29:39 i:3 UnaryEcho reply: <nil>
2024/11/08 11:29:39.008584 request success:0 fault:5
2024/11/08 11:29:39 i:4 UnaryEcho error: rpc error: code = ResourceExhausted desc = intercept fault code:8 rp:28 success:54 fault:61
2024/11/08 11:29:39 i:4 UnaryEcho reply: <nil>
```

### ExampleB
This is an example which will always return code 14 = Unavailable

| Variable           | Value | Description                                                                 |
|--------------------|-------|-----------------------------------------------------------------------------|
| clientfaultpercent | 100   | The client always inserts the headers                                       |
| faultpercent       | 100   | The client header instructs the server to always inject the fault           |
| failcodes          | 14    | Response code is 14                                                         |
| Command            |       | ./client --loops 5 -clientfaultpercent 100 -faultpercent 100 -faultcodes 14 |

```
[das@t:~/Downloads/grpcFaultInjection/cmd/client]$ ./client --loops 5 -clientfaultpercent 100 -faultpercent 100 --faultcodes 14
2024/11/08 11:31:03.541534 request success:0 fault:1
2024/11/08 11:31:03 i:0 UnaryEcho error: rpc error: code = Unavailable desc = intercept fault code:14 rp:22 success:54 fault:62
2024/11/08 11:31:03 i:0 UnaryEcho reply: <nil>
2024/11/08 11:31:03.547073 request success:0 fault:2
2024/11/08 11:31:03 i:1 UnaryEcho error: rpc error: code = Unavailable desc = intercept fault code:14 rp:16 success:54 fault:63
2024/11/08 11:31:03 i:1 UnaryEcho reply: <nil>
2024/11/08 11:31:03.547288 request success:0 fault:3
2024/11/08 11:31:03 i:2 UnaryEcho error: rpc error: code = Unavailable desc = intercept fault code:14 rp:8 success:54 fault:64
2024/11/08 11:31:03 i:2 UnaryEcho reply: <nil>
2024/11/08 11:31:03.547460 request success:0 fault:4
2024/11/08 11:31:03 i:3 UnaryEcho error: rpc error: code = Unavailable desc = intercept fault code:14 rp:33 success:54 fault:65
2024/11/08 11:31:03 i:3 UnaryEcho reply: <nil>
2024/11/08 11:31:03.547620 request success:0 fault:5
2024/11/08 11:31:03 i:4 UnaryEcho error: rpc error: code = Unavailable desc = intercept fault code:14 rp:83 success:54 fault:66
2024/11/08 11:31:03 i:4 UnaryEcho reply: <nil>
```

### ExampleC

This is an example which will always return one of the error codes listed (10, 12, or 14)

| Variable           | Value    | Description                                                                       |
|--------------------|----------|-----------------------------------------------------------------------------------|
| clientfaultpercent | 100      | The client always inserts the headers                                             |
| faultpercent       | 100      | The client header instructs the server to always inject the fault                 |
| failcodes          | 10,12,14 | The server will randomly return one of the codes 10, 12, or 14                    |
| Command            |          | ./client --loops 5 -clientfaultpercent 100 -faultpercent 100 -faultcodes 10,12,14 |

```
[das@t:~/Downloads/grpcFaultInjection/cmd/client]$ ./client --loops 5 -clientfaultpercent 100 -faultpercent 100 --faultcodes 10,12,14
2024/11/08 11:32:29.397913 request success:0 fault:1
2024/11/08 11:32:29 i:0 UnaryEcho error: rpc error: code = Unavailable desc = intercept fault code:14 rp:4 success:54 fault:67
2024/11/08 11:32:29 i:0 UnaryEcho reply: <nil>
2024/11/08 11:32:29.402258 request success:0 fault:2
2024/11/08 11:32:29 i:1 UnaryEcho error: rpc error: code = Aborted desc = intercept fault code:10 rp:9 success:54 fault:68
2024/11/08 11:32:29 i:1 UnaryEcho reply: <nil>
2024/11/08 11:32:29.402484 request success:0 fault:3
2024/11/08 11:32:29 i:2 UnaryEcho error: rpc error: code = Aborted desc = intercept fault code:10 rp:32 success:54 fault:69
2024/11/08 11:32:29 i:2 UnaryEcho reply: <nil>
2024/11/08 11:32:29.402655 request success:0 fault:4
2024/11/08 11:32:29 i:3 UnaryEcho error: rpc error: code = Unimplemented desc = intercept fault code:12 rp:18 success:54 fault:70
2024/11/08 11:32:29 i:3 UnaryEcho reply: <nil>
2024/11/08 11:32:29.402804 request success:0 fault:5
2024/11/08 11:32:29 i:4 UnaryEcho error: rpc error: code = Unimplemented desc = intercept fault code:12 rp:66 success:54 fault:71
2024/11/08 11:32:29 i:4 UnaryEcho reply: <nil>
```

The client headers for this example look like this

<img src="./docs/Screenshot from 2024-11-08 11-36-03.png" alt="xtcp_sampling diagram" width="100%" height="100%"/>

The server headers for this example look like this.
Keep in mind although this is a HTTP 200, it's actually a grpc-status = 14

<img src="./docs/Screenshot from 2024-11-08 11-38-12.png" alt="xtcp_sampling diagram" width="100%" height="100%"/>

## GRPC Metadata

The GRPC library calls the HTTP2 headers "metadata".  I guess this isn't wrong, but it is a little confusing.


If we "follow the HTTP2 stream" in wireshark, we see the HTTP2 headers look like the following

```
:method: POST
:scheme: http
:path: /grpc.examples.echo.Echo/UnaryEcho
:authority: localhost:50052
content-type: application/grpc
user-agent: grpc-go/1.68.0
te: trailers
grpc-timeout: 995078u
faultpercent: 100                   <---- injected header
faultcodes: 10,12,14                <---- injected header

:status: 200                        <--- don't be tricked. this is a fault!
content-type: application/grpc
grpc-status: 10
grpc-message: intercept fault code:10 rp:96 success:0 fault:1    <--- fault
```


See also:

https://grpc.io/docs/guides/metadata/

https://github.com/grpc/grpc-go/blob/master/examples/features/metadata/client/main.go

## GRPC Interceptors

See also the GRPC interceptor documentation:

UnaryClientInterceptor

https://pkg.go.dev/google.golang.org/grpc?utm_source=godoc#UnaryClientInterceptor

UnaryServerInterceptor

https://pkg.go.dev/google.golang.org/grpc?utm_source=godoc#UnaryServerInterceptor


## Tests

Disclaimer:

I wasn't sure about the best way to write the tests via mocking, or whatever, so the tests are pretty expensive.
These tests are running up the GRPC server, and running real GPRC requests across the loopback interface.
This is probably not really required, and it's definitely slow, but you could argue it's pretty realistic.

The point of these tests is really to make sure somebody can confirm this library works, so I hope it does the job.

### Test Functions

All the functions are covered with tests.

### Test_test code

There are serveral tests, loosly following the "Config Matrix" section above.

The probabilitic nature of the tests makes it tricky not to be flakey.

https://github.com/randomizedcoder/grpcFaultInjection/blob/main/cmd/test_test/test_test.go

### Reliability of tests

To ensrue the tests aren't flakey, the test_test Makefile ( https://github.com/randomizedcoder/grpcFaultInjection/blob/main/cmd/test_test/Makefile ),
includes a couple of shortcuts

```
make loop             <--- bash script to loop over the tests 100 times to make sure the test don't fail
make hyperfine        <--- performace measure the tests
make hyperfineDebug   <--- performace measure the tests, in debug mode
```

If you want to measure the performance of the tests, try using "hyperfine"

( https://github.com/sharkdp/hyperfine )

On nixOS, just do "nix-shell hyperfine" "make hyperfine"
```
[das@t:~/Downloads/grpcFaultInjection/cmd/test_test]$ nix-shell -p hyperfine

[nix-shell:~/Downloads/grpcFaultInjection/cmd/test_test]$ make hyperfine
hyperfine \
	--ignore-failure \
	--runs 100 \
	'go test -v'
Benchmark 1: go test -v
  Time (mean ± σ):     521.9 ms ±  77.6 ms    [User: 820.7 ms, System: 283.6 ms]
  Range (min … max):   452.5 ms … 759.2 ms    100 runs

  Warning: Statistical outliers were detected. Consider re-running this benchmark on a quiet system without any interferences from other programs. It might help to use the '--warmup' or '--prepare' options.
```


## Todo

- Streaming inteerceptors
- Mocked tests rather than real GRPC client+server?
- Updates based on feedback