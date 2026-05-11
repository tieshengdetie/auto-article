*   **RPC 服务写在根目录的proto 文件夹下**

    *   go install google.golang.org/protobuf/cmd/protoc-gen-go&#x20;
    *   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc&#x20;

<!---->

*   protoc --go\_out=. --go-grpc\_out=. .\proto\account\account.proto

