# go-grpc-example

演示Golang中基础的gRPC使用

## 导包

```
go get google.golang.org/grpc
```

## 四种RPC方式

### Unary RPC

客户端发送一次请求，然后服务端返回一次响应，这个就是Unary RPC。

例如下面的例子，客户端每发送一个 `CreateCellphoneRequest`就会得到服务端的 `CreateCellphoneResponse`响应。

```protobuf
rpc CreateCellphone(CreateCellphoneRequest) returns (CreateCellphoneResponse);
```

### Server streaming RPC

客户端发送一次请求，服务器不再以一次性的方式返回响应，而是以数据流的方式返回响应（a stream of messages），这个就是Server streaming RPC。

在proto文件内定义rpc服务的时候，在返回值类型前面加上 `stream`关键字，就完成了这种类型的定义。

在这种模式下，服务端涉及的方法为`Send`，用来发送数据流；客户端涉及的方法为`Recv`，用来接收数据流。

例如下面的例子，客户端接收流式响应，每次接收都得到一个 `Cellphone`类型的数据

```protobuf
rpc SearchCellphone(FilterCondition) returns (stream Cellphone);
```

### [Client streaming RPC](https://pkg.go.dev/google.golang.org/grpc#ClientStream)

客户端发送请求的时候以数据流的方式发送请求数据（a stream of messages)，服务端接收数据流，然后返回一次性的响应数据。

这种模式下，客户端涉及的方法为 `Send`用来发送流数据，`CloseAndRecv`用来接收响应数据；服务端涉及的方法为 `Recv`接收请求数据流，`SendAndClose`用来发送响应数据。

例如下面的例子，客户端发送字节流数据给服务端，服务端处理完成后，返回一个响应。

```protobuf
rpc UploadCellphoneCover(stream UploadCellphoneCoverRequest) returns (UploadCellphoneCoverResponse);
```

**注意事项:** 

1. 在goland中的gRPC实现中，客户端调用 `Send`方法发送数据的时候，并不会阻塞等待服务端调用Recv？（从测试来看好像这样)。服务端有可能随时通过返回错误来关闭流，所以客户端每次 `Send`或者最后 `CloseAndRecv`之后都要判断错误。
2. 如果客户端在 `Send`过程中发生了在客户端这一侧的错误，那么直接返回对应的 `error`；但是如果是服务端返回了错误，那么在客户端再使用 `Send`的时候，就会笼统地返回 `io.EOF`错误，更为具体的错误需要在客户端侧调用`RecvMsg(nil)`来获得。

### Bidirectional streaming RPC

客户端和服务端之间双向的数据交互都是数据流。

## gRPC的响应状态使用

gRPC中表示响应的状态和状态码要用到这两个package

```
google.golang.org/grpc/codes
google.golang.org/grpc/status
```

通过 `status.Error(c codes.Code, msg string) error`这个接口可以返回响应错误信息

gRPC定义了一系列内置的错误码，可以在codes中设置，常见的错误码比如 `codes.OK`, `codes.Canceled`, `codes.InvalidArgument`等。

## context.Context在gRPC中的使用

### 超时控制

使用 `context.WithTimeout`或者 `context.WithDeadline`

### 调用取消

使用 `context.WithCancel`

## gPRC中时间类型的使用

在proto文件中 `import "google/protobuf/timestamp.proto"`，这个是protobuf的内置message类型，用来表示时间戳。

生成的go代码中，类型为 `timestamppb.Timestamp`，其有一个 `AsTime`方法可以转化为go的内置 `time.Time`类型。
