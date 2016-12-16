package esic

type CCNxESIC struct {
    writeEncKey []byte
    writeMacKey []byte
    readEncKey []byte
    readMacKey []byte
}

func (n *CCNxESIC) Read() []byte {
    // XXX: do the handshake establishment here...
    return nil
}

func (n *CCNxESIC) Write(data []byte) int {
    // XXX: put in interest and send
    return 0
}
