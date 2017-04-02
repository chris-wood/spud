echo "Installing package dependencies..."
go get -u golang.org/x/crypto/nacl/box
go get -u golang.org/x/crypto/nacl/secretbox
go get -u golang.org/x/crypto/curve25519
go get -u golang.org/x/crypto/poly1305
go get -u golang.org/x/crypto/salsa20/salsa
echo "Done!"
