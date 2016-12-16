package esic

type CCNxESIC struct {
    writeEncKey []byte
    writeMacKey []byte
    readEncKey []byte
    readMacKey []byte
}

func (n *CCNxESIC) Read() []byte {
    // XXX: do the handshake establishment here...
}

func (n *CCNxESIC) Write(data []byte) []byte {
    // XXX: put in interest and send
}
