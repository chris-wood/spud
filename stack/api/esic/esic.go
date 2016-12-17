package esic

type ESIC struct {
    writeEncKey []byte
    writeMacKey []byte
    readEncKey []byte
    readMacKey []byte
}

func (n *ESIC) Read() []byte {
    // XXX: do the handshake establishment here...
    return nil
}

func (n *ESIC) Write(data []byte) int {
    // XXX: put in interest and send
    return 0
}
