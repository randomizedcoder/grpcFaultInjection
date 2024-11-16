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

Please note that if you multiple interceptors use ChainUnaryInterceptor
https://pkg.go.dev/google.golang.org/grpc#ChainUnaryInterceptor
```
	s := grpc.NewServer(
		// Use ChainUnaryInterceptor if you have multiple interceptors
		grpc.ChainUnaryInterceptor(
			//authInterceptor,
			unaryServerFaultInjector.UnaryServerFaultInjector(*debugLevel),
		),
	)
```

https://github.com/randomizedcoder/grpcFaultInjection/blob/main/cmd/server/server.go

## Client usage

The configuration requires configuring both the client side and server side fault injection values.

Select between modes "Modulus" or "Percent", and then enter a integer value:
- Modules 1-10,000
- Percent 1-100

The client needs a configuration as follows:
```
const (
	Modulus Mode = iota
	Percent Mode = 1
)

type ModeValue struct {
	Mode  Mode
	Value int
}

type UnaryClientInterceptorConfig struct {
	Client ModeValue
	Server ModeValue
	Codes  string
}
```
Modulus is preferred over percentage.  If modulus is zero (0) then it tries to use the percentage.

Link to code:
https://github.com/randomizedcoder/grpcFaultInjection/blob/main/pkg/unaryClientFaultInjector/unaryClientFaultInjector_config.go#L19

### Client Modulus Mode
The client can be configured to insert the fault headers based on a modulus of the number
of requests the client has made.

Examples
Client.Mode = Modulus

| Client.Value       | Description                                                          |
| ------------------ | -------------------------------------------------------------------- |
| 0                  | Does NOT use modulus, so you need to configure percentage            |
| 1                  | 1/1 = 100% of the time the metadata(headers) are injected = Always   |
| 2                  | 1/2 of the time the metadata(headers) are injected                   |
| 3                  | 1/3 of the time the metadata(headers) are injected                   |
|...                 |                                                                      |
| 100                | 1/100 of the time the metadata(headers) are injected                 |


### Client Percent Mode
Alternatively, the client can be configured to randomly trigger the fault headers to be injected.

Client.Mode = Percent

Examples

| Client.Value       | Description                                                  |
| ------------------ | ------------------------------------------------------------ |
| 10                 | 10% of the time the metadata(headers) are injected           |
| 100                | 100% of the time the metadata(headers) are injected = Always |

### Server Modulus Mode
The server is controlled by the headers being passed to it.  The client uses the "ServerFaultModulus"
variable to tell the client what values to pass in the "faultmodulus" header

Sever.Mode = Modules

Examples

| Server.Value       | Description                                                          |
| ------------------ | -------------------------------------------------------------------- |
| 0                  | Does NOT use modulus, so you need to configure percentage            |
| 1                  | 1/1 = 100% of the time the server will respond with a fault = Always |
| 2                  | 1/2 of the time the server will inject the fault                     |
| 3                  | 1/3 of the time the server will inject the fault                     |
|...                 |                                                                      |
| 100                | 1/100 of the time the server will inject the fault                   |

### Server Percent Mode
The client can be configured with "ServerFaultPercent", which instructs the client
to insert metadata into the request "faultpercent", which the GPRC server uses to
randomly insert faults at this percentage. e.g. 10 means server has a 10% chance of injecting a
fault.

The configuration ServerFaultPercent injects the "faultpercent" header which is
passed to the server.

Sever.Mode = Percent

Examples

| Server.Value       | Description                                                     |
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

Please keep in mind the Client Modulus and Server Modulus value result in fault
injection behaviour like the following:

| Client Modulus Value | Server Modulus Value | Inject        | Comment                                  |
|----------------------|----------------------|---------------|------------------------------------------|
| 1                    | 0                    | 0             | Doesn't do much                          |
| 0                    | 1                    | 0             |                                          |
|                      |                      |               |                                          |
| 1                    | 10                   | 1/10          | Recommended to use 1 on one end          |
| 1                    | 2                    | 1/2           |                                          |
| 10                   | 1                    | 1/10          |                                          |
| 2                    | 1                    | 1/2           |                                          |
|                      |                      |               |                                          |
| 10                   | 10                   | 1/(10*10=100) | Tricky to reason about                   |
| 2                    | 2                    | 1/(2*2=4)     |                                          |
| 1                    | 2                    | 1/2           |                                          |
|                      |                      |               |                                          |
| 1                    | 1                    | 1             | Unlikely to be successful!               |

## Config Matrix - Percentage

Please keep in mind the Client Percent and Server Percent values result in fault
injection probabilities like the following:

| Client Percent Value | Server Percent Value | Probability | Comment                                  |
|----------------------|----------------------|-------------|------------------------------------------|
| 100                  | 0                    | 0%          | Doesn't do much                          |
| 0                    | 100                  | 0%          |                                          |
|                      |                      |             |                                          |
| 100                  | 10                   | 10%         | Recommended to use 100% on one end       |
| 100                  | 50                   | 50%         |                                          |
| 10                   | 100                  | 10%         |                                          |
| 50                   | 100                  | 50%         |                                          |
|                      |                      |             |                                          |
| 10                   | 10                   | 1%          | Tricky to reason about                   |
| 50                   | 50                   | 25%         |                                          |
| 100                  | 50                   | 50%         |                                          |
|                      |                      |             |                                          |
| 100                  | 100                  | 100%        | Unlikely to be successful!               |

( test_test.go tries to follow this table )

## Modulus Examples

### ExampleA

Always inject errors, with any random error code.
```
./client \
	-clientmode Modulus \
	-clientvalue 1 \
	-servermode Modulus \
	-servervalue 1 \
	-loops 5
```

```
[das@t:~/Downloads/grpcFaultInjection/cmd/client]$ ./client -clientmode Modulus -clientvalue 1 -servermode Modulus -servervalue 1 -loops 5
2024/11/12 16:37:12.245962 fault request success:0 fault:1
2024/11/12 16:37:12 i:0 UnaryEcho error: rpc error: code = Unimplemented desc = intercept fault code:12 counter:101 success:75 fault:26
2024/11/12 16:37:12.251152 fault request success:0 fault:2
2024/11/12 16:37:12 i:1 UnaryEcho error: rpc error: code = Aborted desc = intercept fault code:10 counter:102 success:75 fault:27
2024/11/12 16:37:12.251388 fault request success:0 fault:3
2024/11/12 16:37:12 i:2 UnaryEcho error: rpc error: code = Aborted desc = intercept fault code:10 counter:103 success:75 fault:28
2024/11/12 16:37:12.251559 fault request success:0 fault:4
2024/11/12 16:37:12 i:3 UnaryEcho error: rpc error: code = Unimplemented desc = intercept fault code:12 counter:104 success:75 fault:29
2024/11/12 16:37:12.251716 fault request success:0 fault:5
2024/11/12 16:37:12 i:4 UnaryEcho error: rpc error: code = Aborted desc = intercept fault code:10 counter:105 success:75 fault:30
```

### ExampleB

Server will inject faults 1/2 of the time.
```
./client \
	-clientmode Modulus \
	-clientvalue 1 \
	-servermode Modulus \
	-servervalue 2 \
	-loops 6
```


```
[das@t:~/Downloads/grpcFaultInjection/cmd/client]$ ./client -clientmode Modulus -clientvalue 1 -servermode Modulus -servervalue 2 -loops 6
2024/11/12 16:37:56.942287 fault request success:0 fault:1
2024/11/12 16:37:56 i:0 UnaryEcho reply: message:"Try and Success"
2024/11/12 16:37:56.946410 fault request success:0 fault:2
2024/11/12 16:37:56 i:1 UnaryEcho error: rpc error: code = Aborted desc = intercept fault code:10 counter:112 success:78 fault:34
2024/11/12 16:37:56.946641 fault request success:0 fault:3
2024/11/12 16:37:56 i:2 UnaryEcho reply: message:"Try and Success"
2024/11/12 16:37:56.946836 fault request success:0 fault:4
2024/11/12 16:37:56 i:3 UnaryEcho error: rpc error: code = Unavailable desc = intercept fault code:14 counter:114 success:79 fault:35
2024/11/12 16:37:56.947003 fault request success:0 fault:5
2024/11/12 16:37:56 i:4 UnaryEcho reply: message:"Try and Success"
2024/11/12 16:37:56.947148 fault request success:0 fault:6
2024/11/12 16:37:56 i:5 UnaryEcho error: rpc error: code = Unavailable desc = intercept fault code:14 counter:116 success:80 fault:36
```

### ExampleC

Server will inject fault 1/4 of the time.

```
./client \
	-clientmode Modulus \
	-clientvalue 1 \
	-servermode Modulus \
	-servervalue 4 \
	-loops 8
```

```
[das@t:~/Downloads/grpcFaultInjection/cmd/client]$ ./client -clientmode Modulus -clientvalue 1 -servermode Modulus -servervalue 4 -loops 8
2024/11/12 16:38:21.822487 fault request success:0 fault:1
2024/11/12 16:38:21 i:0 UnaryEcho reply: message:"Try and Success"
2024/11/12 16:38:21.827774 fault request success:0 fault:2
2024/11/12 16:38:21 i:1 UnaryEcho reply: message:"Try and Success"
2024/11/12 16:38:21.827962 fault request success:0 fault:3
2024/11/12 16:38:21 i:2 UnaryEcho reply: message:"Try and Success"
2024/11/12 16:38:21.828162 fault request success:0 fault:4
2024/11/12 16:38:21 i:3 UnaryEcho error: rpc error: code = Unimplemented desc = intercept fault code:12 counter:120 success:83 fault:37
2024/11/12 16:38:21.828336 fault request success:0 fault:5
2024/11/12 16:38:21 i:4 UnaryEcho reply: message:"Try and Success"
2024/11/12 16:38:21.828503 fault request success:0 fault:6
2024/11/12 16:38:21 i:5 UnaryEcho reply: message:"Try and Success"
2024/11/12 16:38:21.828653 fault request success:0 fault:7
2024/11/12 16:38:21 i:6 UnaryEcho reply: message:"Try and Success"
2024/11/12 16:38:21.828794 fault request success:0 fault:8
2024/11/12 16:38:21 i:7 UnaryEcho error: rpc error: code = Aborted desc = intercept fault code:10 counter:124 success:86 fault:38
```


## Percent Examples

### ExampleA

This is an example that will always return a GRPC error status, the status code will be random.

```
./client \
	-clientmode Percent \
	-clientvalue 100 \
	-servermode Percent \
	-servervalue 100 \
	-loops 5
```

```
[das@t:~/Downloads/grpcFaultInjection/cmd/client]$ ./client -clientmode Percent -clientvalue 100 -servermode Percent -servervalue 100 -loops 5
2024/11/12 16:39:11.534763 fault request success:0 fault:1
2024/11/12 16:39:11 i:0 UnaryEcho error: rpc error: code = Aborted desc = intercept fault code:10 counter:125 success:86 fault:39
2024/11/12 16:39:11.539314 fault request success:0 fault:2
2024/11/12 16:39:11 i:1 UnaryEcho error: rpc error: code = Unavailable desc = intercept fault code:14 counter:126 success:86 fault:40
2024/11/12 16:39:11.539540 fault request success:0 fault:3
2024/11/12 16:39:11 i:2 UnaryEcho error: rpc error: code = Unimplemented desc = intercept fault code:12 counter:127 success:86 fault:41
2024/11/12 16:39:11.539700 fault request success:0 fault:4
2024/11/12 16:39:11 i:3 UnaryEcho error: rpc error: code = Unavailable desc = intercept fault code:14 counter:128 success:86 fault:42
2024/11/12 16:39:11.539864 fault request success:0 fault:5
2024/11/12 16:39:11 i:4 UnaryEcho error: rpc error: code = Aborted desc = intercept fault code:10 counter:129 success:86 fault:43
```

### ExampleB
This is an example which will always return code 14 = Unavailable

```
./client \
	-clientmode Percent \
	-clientvalue 100 \
	-servermode Percent \
	-servervalue 100 \
	-loops 5 \
	-codes 14
```

```
[das@t:~/Downloads/grpcFaultInjection/cmd/client]$ ./client -clientmode Percent -clientvalue 100 -servermode Percent -servervalue 100 -loops 5 -codes 14
2024/11/12 16:39:48.790843 fault request success:0 fault:1
2024/11/12 16:39:48 i:0 UnaryEcho error: rpc error: code = Unavailable desc = intercept fault code:14 counter:130 success:86 fault:44
2024/11/12 16:39:48.797675 fault request success:0 fault:2
2024/11/12 16:39:48 i:1 UnaryEcho error: rpc error: code = Unavailable desc = intercept fault code:14 counter:131 success:86 fault:45
2024/11/12 16:39:48.797884 fault request success:0 fault:3
2024/11/12 16:39:48 i:2 UnaryEcho error: rpc error: code = Unavailable desc = intercept fault code:14 counter:132 success:86 fault:46
2024/11/12 16:39:48.798095 fault request success:0 fault:4
2024/11/12 16:39:48 i:3 UnaryEcho error: rpc error: code = Unavailable desc = intercept fault code:14 counter:133 success:86 fault:47
2024/11/12 16:39:48.798298 fault request success:0 fault:5
2024/11/12 16:39:48 i:4 UnaryEcho error: rpc error: code = Unavailable desc = intercept fault code:14 counter:134 success:86 fault:48
```

### ExampleC

This is an example which will always return one of the error codes listed (10, 12, or 14)

```
./client \
	-clientmode Percent \
	-clientvalue 100 \
	-servermode Percent \
	-servervalue 100 \
	-loops 6 \
	-codes 10,12,14
```

```
[das@t:~/Downloads/grpcFaultInjection/cmd/client]$ ./client -clientmode Percent -clientvalue 100 -servermode Percent -servervalue 100 -loops 6 -codes 10,12,14
2024/11/12 16:40:26.895077 fault request success:0 fault:1
2024/11/12 16:40:26 i:0 UnaryEcho error: rpc error: code = Unimplemented desc = intercept fault code:12 counter:135 success:86 fault:49
2024/11/12 16:40:26.899571 fault request success:0 fault:2
2024/11/12 16:40:26 i:1 UnaryEcho error: rpc error: code = Unimplemented desc = intercept fault code:12 counter:136 success:86 fault:50
2024/11/12 16:40:26.899781 fault request success:0 fault:3
2024/11/12 16:40:26 i:2 UnaryEcho error: rpc error: code = Aborted desc = intercept fault code:10 counter:137 success:86 fault:51
2024/11/12 16:40:26.899986 fault request success:0 fault:4
2024/11/12 16:40:26 i:3 UnaryEcho error: rpc error: code = Aborted desc = intercept fault code:10 counter:138 success:86 fault:52
2024/11/12 16:40:26.900140 fault request success:0 fault:5
2024/11/12 16:40:26 i:4 UnaryEcho error: rpc error: code = Unimplemented desc = intercept fault code:12 counter:139 success:86 fault:53
2024/11/12 16:40:26.900317 fault request success:0 fault:6
2024/11/12 16:40:26 i:5 UnaryEcho error: rpc error: code = Aborted desc = intercept fault code:10 counter:140 success:86 fault:54
```

The client headers for this example look like this

<img src="./docs/Screenshot from 2024-11-08 11-36-03.png" alt="xtcp_sampling diagram" width="100%" height="100%"/>

The server headers for this example look like this.
Keep in mind although this is a HTTP 200, it's actually a grpc-status = 14

<img src="./docs/Screenshot from 2024-11-08 11-38-12.png" alt="xtcp_sampling diagram" width="100%" height="100%"/>


### ExampleD

The client and server do not need to use the same mode

```
./client \
	-clientmode Modulus \
	-clientvalue 2 \
	-servermode Percent \
	-servervalue 100 \
	-loops 10 \
	-codes 10,12,14 \
	-debugLevel 111
```


```
[das@t:~/Downloads/grpcFaultInjection/cmd/client]$ ./client -clientmode Modulus -clientvalue 2 -servermode Percent -servervalue 100 -loops 10 -codes 10,12,14 -debugLevel 111
2024/11/12 16:47:46.209636 no fault request success:1 fault:0 ~= 0
2024/11/12 16:47:46 i:0 UnaryEcho reply: message:"Try and Success"
2024/11/12 16:47:46.223497 UnaryClientFaultInjector counter:2
2024/11/12 16:47:46.223503 fault request success:1 fault:1 ~= 1
2024/11/12 16:47:46.223507 md:map[faultcodes:[10,12,14] faultpercent:[100]]
2024/11/12 16:47:46 i:1 UnaryEcho error: rpc error: code = Unavailable desc = intercept fault code:14 counter:194 success:96 fault:98
2024/11/12 16:47:46.223733 no fault request success:2 fault:1 ~= 0.5
2024/11/12 16:47:46 i:2 UnaryEcho reply: message:"Try and Success"
2024/11/12 16:47:46.223909 UnaryClientFaultInjector counter:4
2024/11/12 16:47:46.223913 fault request success:2 fault:2 ~= 1
2024/11/12 16:47:46.223916 md:map[faultcodes:[10,12,14] faultpercent:[100]]
2024/11/12 16:47:46 i:3 UnaryEcho error: rpc error: code = Unavailable desc = intercept fault code:14 counter:196 success:97 fault:99
2024/11/12 16:47:46.224069 no fault request success:3 fault:2 ~= 0.667
2024/11/12 16:47:46 i:4 UnaryEcho reply: message:"Try and Success"
2024/11/12 16:47:46.224222 UnaryClientFaultInjector counter:6
2024/11/12 16:47:46.224225 fault request success:3 fault:3 ~= 1
2024/11/12 16:47:46.224228 md:map[faultcodes:[10,12,14] faultpercent:[100]]
2024/11/12 16:47:46 i:5 UnaryEcho error: rpc error: code = Unimplemented desc = intercept fault code:12 counter:198 success:98 fault:100
2024/11/12 16:47:46.224381 no fault request success:4 fault:3 ~= 0.75
2024/11/12 16:47:46 i:6 UnaryEcho reply: message:"Try and Success"
2024/11/12 16:47:46.224514 UnaryClientFaultInjector counter:8
2024/11/12 16:47:46.224517 fault request success:4 fault:4 ~= 1
2024/11/12 16:47:46.224519 md:map[faultcodes:[10,12,14] faultpercent:[100]]
2024/11/12 16:47:46 i:7 UnaryEcho error: rpc error: code = Aborted desc = intercept fault code:10 counter:200 success:99 fault:101
2024/11/12 16:47:46.224694 no fault request success:5 fault:4 ~= 0.8
2024/11/12 16:47:46 i:8 UnaryEcho reply: message:"Try and Success"
2024/11/12 16:47:46.224826 UnaryClientFaultInjector counter:10
2024/11/12 16:47:46.224829 fault request success:5 fault:5 ~= 1
2024/11/12 16:47:46.224832 md:map[faultcodes:[10,12,14] faultpercent:[100]]
2024/11/12 16:47:46 i:9 UnaryEcho error: rpc error: code = Unavailable desc = intercept fault code:14 counter:202 success:100 fault:102
2024/11/12 16:47:46 Complete.  success:5 fault:5
```

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

```
[nix-shell:~/Downloads/grpcFaultInjection/cmd/test_test]$ go test -run TestComprehensive -v
=== RUN   TestComprehensive
    test_test.go:64: listen on address localhost:50053
    test_test.go:320: run tests
    test_test.go:324: tt.Name:1/1 client, 1/1 server fault, loops 100, = 100%
=== RUN   TestComprehensive/1/1_client,_1/1_server_fault,_loops_100,_=_100%
    test_test.go:395: tt.Name:1/1 client, 1/1 server fault, loops 100, = 100% success:0 minSuccess:0 = good
    test_test.go:403: tt.Name:1/1 client, 1/1 server fault, loops 100, = 100% success:0 maxSuccess:0 = good
    test_test.go:411: tt.Name:1/1 client, 1/1 server fault, loops 100, = 100% fault:100 minFault:100 = good
    test_test.go:419: tt.Name:1/1 client, 1/1 server fault, loops 100, = 100% fault:100 minFault:100 = good
```

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