# Nettest

This package provides a simple way to get an open TCP/UDP port that can then be used for testing.

## Example Usage
```go
port := nettest.GetTCP()

address := net.JoinHostPort("0.0.0.0", strconv.Itoa(port))

listener, err := net.Listen("tcp", address)
if err != nil {
    return err
}

rpc.Accept(listener)
```