# grpcFaultInjection

GRPC fault injection library using interceptors

This is a small library with client + server GRPC intercepts

The client injects metadata (http2 headers) to request that the server does
fault injection

This library is designed to allow the client to control the fault injection,
and is generally designed to allow testing of error handling code

Hopefully this library is easy to integrate into an existing GRPC eco-system,
and will provide value for failure mode testing

Please note there are x2 differnet modes:
- Modulus
- Percent
The modulus mode makes if reliable testing, while percentage is random probability,
which can lead to flaky tests.

## Overview Diagram

<img src="./docs/Screenshot from 2024-11-12 10-51-06.png" alt="xtcp_sampling diagram" width="100%" height="100%"/>


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
	ClientFaultModulus int
	ClientFaultPercent int
	ServerFaultModulus int
	ServerFaultPercent int
	ServerFaultCodes   string
}
```
Modulus is preferred over percentage.  If modulus is zero (0) then it tries to use the percentage.

### ClientFaultModulus
The client can be configured to insert the fault headers based on a modulus of the number
of requests the client has made.

Examples

| ClientFaultModulus | Description                                                          |
| ------------------ | -------------------------------------------------------------------- |
| 0                  | Does NOT use modulus, so you need to configure percentage            |
| 1                  | 1/1 = 100% of the time the metadata(headers) are injected = Always   |
| 2                  | 1/2 of the time the metadata(headers) are injected                   |
| 3                  | 1/3 of the time the metadata(headers) are injected                   |
|...                 |                                                                      |
| 100                | 1/100 of the time the metadata(headers) are injected                 |


### ClientFaultPercent
Alternatively, the client can be configured to randomly trigger the fault headers to be injected.

Examples

| ClientFaultPercent | Description                                                  |
| ------------------ | ------------------------------------------------------------ |
| 10                 | 10% of the time the metadata(headers) are injected           |
| 100                | 100% of the time the metadata(headers) are injected = Always |

### ServerFaultModulus
The server is controlled by the headers being passed to it.  The client uses the "ServerFaultModulus"
variable to tell the client what values to pass in the "faultmodulus" header

Examples

| ServerFaultModulus | Description                                                          |
| ------------------ | -------------------------------------------------------------------- |
| 0                  | Does NOT use modulus, so you need to configure percentage            |
| 1                  | 1/1 = 100% of the time the server will respond with a fault = Always |
| 2                  | 1/2 of the time the server will inject the fault                     |
| 3                  | 1/3 of the time the server will inject the fault                     |
|...                 |                                                                      |
| 100                | 1/100 of the time the server will inject the fault                   |

### ServerFaultPercent
The client can be configured with "ServerFaultPercent", which instructs the client
to insert metadata into the request "faultpercent", which the GPRC server uses to
randomly insert faults at this percentage. e.g. 10 means server has a 10% chance of injecting a
fault.

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

The configuration ServerFaultCodes configures the client to injects the "faultcodes" header,
which is passed to the server, and then the server uses this configuration to control the
fault response codes it will return.

If "faultcodes" is NOT supplied, any random valid GRPC status code, except zero (0), is returned

Examples

| "faultcodes  "     | Description                                                         |
| ------------------ | ------------------------------------------------------------------- |
| 14                 | If the server injects the fault, the only return status is 14       |
| 10,12,14           | If the server injects the fault, possible status codes are 10,12,14 |
| <not set >         | If the server injects the fault, codes 1-16 are possible            |


Possible failcodes are:
https://github.com/grpc/grpc/blob/master/doc/statuscodes.md

## Config Matrix - Modulus

Please keep in mind the ClientFaultModulus and ServerFaultModulus result in fault
injection behaviour like the following:

| ClientFaultModulus | ServerFaultModulus | Inject        | Comment                                  |
|--------------------|--------------------|---------------|------------------------------------------|
| 1                  | 0                  | 0             | Doesn't do much                          |
| 0                  | 1                  | 0             |                                          |
|                    |                    |               |                                          |
| 1                  | 10                 | 1/10          | Recommended to use 1 on one end          |
| 1                  | 2                  | 1/2           |                                          |
| 10                 | 1                  | 1/10          |                                          |
| 2                  | 1                  | 1/2           |                                          |
|                    |                    |               |                                          |
| 10                 | 10                 | 1/(10*10=100) | Tricky to reason about                   |
| 2                  | 2                  | 1/(2*2=4)     |                                          |
| 1                  | 2                  | 1/2           |                                          |
|                    |                    |               |                                          |
| 1                  | 1                  | 1             | Unlikely to be successful!               |

## Config Matrix - Percentage

Please keep in mind the ClientFaultPercent and ServerFaultPercent result in fault
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

## Modulus Examples

### ExampleA

Always inject errors, with any random error code.
```
./client -clientfaultmodulus 1 -faultmodulus 1
```

```
[nix-shell:~/Downloads/grpcFaultInjection/cmd/client]$ ./client -clientfaultmodulus 1 -faultmodulus 1
2024/11/12 10:24:54.198618 fault request success:0 fault:1
2024/11/12 10:24:54 i:0 UnaryEcho error: rpc error: code = DeadlineExceeded desc = intercept fault code:4 counter:317 success:158 fault:159
2024/11/12 10:24:54.290040 fault request success:0 fault:2
2024/11/12 10:24:54 i:1 UnaryEcho error: rpc error: code = DeadlineExceeded desc = intercept fault code:4 counter:318 success:158 fault:160
2024/11/12 10:24:54.290285 fault request success:0 fault:3
2024/11/12 10:24:54 i:2 UnaryEcho error: rpc error: code = DeadlineExceeded desc = intercept fault code:4 counter:319 success:158 fault:161
2024/11/12 10:24:54.290547 fault request success:0 fault:4
2024/11/12 10:24:54 i:3 UnaryEcho error: rpc error: code = ResourceExhausted desc = intercept fault code:8 counter:320 success:158 fault:162
2024/11/12 10:24:54.290809 fault request success:0 fault:5
2024/11/12 10:24:54 i:4 UnaryEcho error: rpc error: code = DeadlineExceeded desc = intercept fault code:4 counter:321 success:158 fault:163
2024/11/12 10:24:54.291004 fault request success:0 fault:6
2024/11/12 10:24:54 i:5 UnaryEcho error: rpc error: code = DeadlineExceeded desc = intercept fault code:4 counter:322 success:158 fault:164
2024/11/12 10:24:54.291197 fault request success:0 fault:7
2024/11/12 10:24:54 i:6 UnaryEcho error: rpc error: code = ResourceExhausted desc = intercept fault code:8 counter:323 success:158 fault:165
2024/11/12 10:24:54.291463 fault request success:0 fault:8
2024/11/12 10:24:54 i:7 UnaryEcho error: rpc error: code = Unavailable desc = intercept fault code:14 counter:324 success:158 fault:166
2024/11/12 10:24:54.291693 fault request success:0 fault:9
2024/11/12 10:24:54 i:8 UnaryEcho error: rpc error: code = Unavailable desc = intercept fault code:14 counter:325 success:158 fault:167
2024/11/12 10:24:54.291876 fault request success:0 fault:10
2024/11/12 10:24:54 i:9 UnaryEcho error: rpc error: code = DeadlineExceeded desc = intercept fault code:4 counter:326 success:158 fault:168
```

### ExampleB

Server will inject faults 1/2 of the time.
```
./client -clientfaultmodulus 1 -faultmodulus 2
```


```
[nix-shell:~/Downloads/grpcFaultInjection/cmd/client]$ ./client -clientfaultmodulus 1 -faultmodulus 2
2024/11/12 10:25:00.502643 fault request success:0 fault:1
2024/11/12 10:25:00 i:0 UnaryEcho reply: message:"Try and Success"
2024/11/12 10:25:00.507951 fault request success:0 fault:2
2024/11/12 10:25:00 i:1 UnaryEcho error: rpc error: code = DeadlineExceeded desc = intercept fault code:4 counter:328 success:159 fault:169
2024/11/12 10:25:00.508160 fault request success:0 fault:3
2024/11/12 10:25:00 i:2 UnaryEcho reply: message:"Try and Success"
2024/11/12 10:25:00.508379 fault request success:0 fault:4
2024/11/12 10:25:00 i:3 UnaryEcho error: rpc error: code = ResourceExhausted desc = intercept fault code:8 counter:330 success:160 fault:170
2024/11/12 10:25:00.508550 fault request success:0 fault:5
2024/11/12 10:25:00 i:4 UnaryEcho reply: message:"Try and Success"
2024/11/12 10:25:00.508735 fault request success:0 fault:6
2024/11/12 10:25:00 i:5 UnaryEcho error: rpc error: code = DeadlineExceeded desc = intercept fault code:4 counter:332 success:161 fault:171
2024/11/12 10:25:00.508927 fault request success:0 fault:7
2024/11/12 10:25:00 i:6 UnaryEcho reply: message:"Try and Success"
2024/11/12 10:25:00.509115 fault request success:0 fault:8
2024/11/12 10:25:00 i:7 UnaryEcho error: rpc error: code = ResourceExhausted desc = intercept fault code:8 counter:334 success:162 fault:172
2024/11/12 10:25:00.509297 fault request success:0 fault:9
2024/11/12 10:25:00 i:8 UnaryEcho reply: message:"Try and Success"
2024/11/12 10:25:00.509446 fault request success:0 fault:10
2024/11/12 10:25:00 i:9 UnaryEcho error: rpc error: code = Unavailable desc = intercept fault code:14 counter:336 success:163 fault:173
```

### ExampleC

Server will inject fault 1/4 of the time.

```
./client -clientfaultmodulus 1 -faultmodulus 4
```


```
[nix-shell:~/Downloads/grpcFaultInjection/cmd/client]$ ./client -clientfaultmodulus 1 -faultmodulus 4
2024/11/12 10:25:11.998466 fault request success:0 fault:1
2024/11/12 10:25:12 i:0 UnaryEcho reply: message:"Try and Success"
2024/11/12 10:25:12.003191 fault request success:0 fault:2
2024/11/12 10:25:12 i:1 UnaryEcho error: rpc error: code = ResourceExhausted desc = intercept fault code:8 counter:348 success:171 fault:177
2024/11/12 10:25:12.003461 fault request success:0 fault:3
2024/11/12 10:25:12 i:2 UnaryEcho reply: message:"Try and Success"
2024/11/12 10:25:12.003659 fault request success:0 fault:4
2024/11/12 10:25:12 i:3 UnaryEcho reply: message:"Try and Success"
2024/11/12 10:25:12.003812 fault request success:0 fault:5
2024/11/12 10:25:12 i:4 UnaryEcho reply: message:"Try and Success"
2024/11/12 10:25:12.003951 fault request success:0 fault:6
2024/11/12 10:25:12 i:5 UnaryEcho error: rpc error: code = DeadlineExceeded desc = intercept fault code:4 counter:352 success:174 fault:178
2024/11/12 10:25:12.004095 fault request success:0 fault:7
2024/11/12 10:25:12 i:6 UnaryEcho reply: message:"Try and Success"
2024/11/12 10:25:12.004242 fault request success:0 fault:8
2024/11/12 10:25:12 i:7 UnaryEcho reply: message:"Try and Success"
2024/11/12 10:25:12.004402 fault request success:0 fault:9
2024/11/12 10:25:12 i:8 UnaryEcho reply: message:"Try and Success"
2024/11/12 10:25:12.004549 fault request success:0 fault:10
2024/11/12 10:25:12 i:9 UnaryEcho error: rpc error: code = DeadlineExceeded desc = intercept fault code:4 counter:356 success:177 fault:179
```


## Percent Examples

### ExampleA

This is an example that will always return a GRPC error status, the status code will be random.

```
./client --loops 5 -clientfaultpercent 100 -faultpercent 100
```

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

```
./client --loops 5 -clientfaultpercent 100 -faultpercent 100 --faultcodes 14
```

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

```
./client --loops 5 -clientfaultpercent 100 -faultpercent 100 --faultcodes 10,12,14
```

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
grpc-timeout: 995078u                                       <--- ( ctx timeout )
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