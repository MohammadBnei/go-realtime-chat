# Protobuf - Oneof

Assume you have a protobuf file named `event.proto` that defines an `Event` message with an `oneof` for event types:

```protobuf
syntax = "proto3";
package events;

message Event {
  oneof event_type {
    string created;
    string updated;
    string deleted;
  }
}
```

To generate the Golang code for this protobuf file, run the following command:

```sh
$ protoc --go_out=. event.proto
```

This will generate a `event.pb.go` file that includes all of the necessary code for marshalling and unmarshalling protobuf-encoded `Event` messages.

Now, you can use the generated `Event` struct in your gRPC server and client code. Here's an example of a gRPC server implementation that accepts an `Event` message:

```go
package main

import (
    "context"
    "log"
    "net"

    "google.golang.org/grpc"
    pb "path/to/your/protobuf/file"
)

type server struct {}

func (s *server) CreateEvent(ctx context.Context, req *pb.Event) (*pb.Empty, error) {
    switch req.Event.(type) {
    case *pb.Event_Created:
        log.Printf("Received Created event: %s", req.Created)
    case *pb.Event_Updated:
        log.Printf("Received Updated event: %s", req.Updated)
    case *pb.Event_Deleted:
        log.Printf("Received Deleted event: %s", req.Deleted)
    default:
        return nil, status.Errorf(codes.InvalidArgument, "Event type not recognized: %v", req)
    }

    // TODO: implement event handling

    return &pb.Empty{}, nil
}

func main() {
    // Create a gRPC server
    lis, err := net.Listen("tcp", ":8080")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }
    s := grpc.NewServer()

    // Register the server
    pb.RegisterEventServiceServer(s, &server{})

    // Start the server
    log.Println("Starting server...")
    if err := s.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
```

This implementation defines a gRPC server that can handle an `Event` message in the `CreateEvent` method. The implementation uses a type switch to handle the different types of events that could be encoded in the message.

Note that in the generated `Event` struct, the `XXX_event_type` field is used to indicate which field is present in the message. In the switch statement, we use a type assertion to determine which field is present and extract its value accordingly.

You can also use the `Event` struct in a gRPC client implementation to send events to the server. Here's an example of a client implementation that sends a `Created` event:

```go
package main

import (
    "context"
    "log"

    "google.golang.org/grpc"
    pb "path/to/your/protobuf/file"
)

func main() {
    // Set up a connection to the server
    conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to dial server: %v", err)
    }
    defer conn.Close()

    // Create a gRPC client
    c := pb.NewEventServiceClient(conn)

    // Send a Created event
    event := &pb.Event{
        Event: &pb.Event_Created{
            Created: "Some data for the Created event",
        },
    }
    _, err = c.CreateEvent(context.Background(), event)
    if err != nil {
        log.Fatalf("Failed to send event: %v", err)
    }
}
```

This implementation creates a gRPC client and sends a `Created` event to the server. The `Event` message is constructed by setting the `Created` field and leaving the other fields unset. When the message is marshalled, the `XXX_event_type` field will be set to indicate that the `Created` field is present.