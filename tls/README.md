```shell
go run $GOROOT/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
```

1) First it generates a 2048-bit RSA key pair, 
   which is a cryptographically secure `public key` and `private key`.
   
2) It then stores the `private key` in a `key.pem` file, 
   and generates a self-signed TLS certificate for the host localhost 
   containing the `public key` — which it stores in a `cert.pem` file. 
   
3) Both the private key and certificate are PEM encoded, 
   which is the standard format used by most TLS implementations.

And that’s it! We’ve now got a `self-signed TLS certificate`