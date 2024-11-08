# grpcFaultInjection

GRPC fault injection library using interceptors

This is a small library with client + server GRPC intercepts

The client injects metadata (http2 headers) to request that the server does
fault injection

This library is designed to allow the client to control the fault injection,
and is generally designed to allow testing of error handling code



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

| ClientFaultPercent | Description                                                  |
| ------------------ | ------------------------------------------------------------ |
| 10                 | 10% of the time the metadata(headers) are injected           |
| 100                | 100% of the time the metadata(headers) are injected = Always |


### ServerFaultPercent
The client can make requests with "faultpercent" and "faultcodes" metadata(headers)

The configuration ServerFaultPercent injects the "faultpercent" header which is
passed to the server.

| ServerFaultPercent | Description                                                     |
| ------------------ | --------------------------------------------------------------- |
| 50                 | 50% there is a 50% chance that the server will return a fault   |
| 90                 | 90% chance the server will always return a fault                |
| 100                | 100% of the time the server will always return a fault = Always |

The client adds the "faultpercent" header which is passed to the server:

| "faultpercent"     | Description                                                     |
| ------------------ | --------------------------------------------------------------- |
| 50                 | 50% there is a 50% chance that the server will return a fault   |
| 90                 | 90% chance the server will always return a fault                |
| 100                | 100% of the time the server will always return a fault = Always |

### ServerFaultCodes

The configuration ServerFaultCodes injects the "faultcodes" header which is
passed to the server.

If "faultcodes" is NOT supplied, any random valid GRPC status code, except zero (0), is returned

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

| ClientFaultPercent | ServerFaultPercent | Probability |
|--------------------|--------------------|-------------|
| 10                 | 10                 | 1%          |
| 50                 | 50                 | 25%         |
| 100                | 50                 | 50%         |
| 50                 | 100                | 50%         |
| 100                | 100                | 100%        |

## Examples

### ExampleA
clientfaultpercent = 100 ( client always inserts the headers )
failpercent = 100 ( the headers instruct the server to always inject the fault )
failcodes is not specified, so we get any random status code
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
clientfaultpercent = 100 ( client always inserts the headers )
failpercent = 100 ( the headers instruct the server to always inject the fault )
faultcodes = 14 (unavailable), so the server will only return code 14
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
clientfaultpercent = 100 ( client always inserts the headers )
failpercent = 100 ( the headers instruct the server to always inject the fault )
faultcodes = 10,12,14, so the server will randomly return one of the codes 10, 12, or 14
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

See also:

https://grpc.io/docs/guides/metadata/

https://github.com/grpc/grpc-go/blob/master/examples/features/metadata/client/main.go

## GRPC Interceptors

UnaryClientInterceptor

https://pkg.go.dev/google.golang.org/grpc?utm_source=godoc#UnaryClientInterceptor

UnaryServerInterceptor

https://pkg.go.dev/google.golang.org/grpc?utm_source=godoc#UnaryServerInterceptor

## Todo
Streaming version