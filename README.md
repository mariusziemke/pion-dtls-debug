# Pion DTLS issue example

When pion is the offerer and a renegotiation occurs from the answerer we run into a DTLS issue.

```
Uncaught (in promise) DOMException: Failed to execute 'setRemoteDescription' on 'RTCPeerConnection': Failed to set remote answer sdp: Failed to apply the description for m= section with mid='0': Failed to set SSL role for the transport.
```

### Steps to reproduce

This example reproduces the issue.

1. `go run main.go` starts a websocket listener on `:4444`
2. Open the `index.html` file in `/browser`
3. Open the console
4. Click on connect
5. Wait for the connection to be established
6. Click on renegotiate
7. You should now see the error message in the console log