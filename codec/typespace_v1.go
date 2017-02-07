package codec

// hop-by-hop headers
const T_INT_LIFE uint16 = 0x0001
const T_CACHE_TIME uint16 = 0x0002
const T_MSG_HASH uint16 = 0x0003

// top-level TLVs
const T_INTEREST uint16 = 0x0001
const T_OBJECT uint16 = 0x0002
const T_VALALG uint16 = 0x0003
const T_VALPAYLOAD uint16 = 0x0004
const T_MANIFEST uint16 = 0x0006

// Message body containers
const T_NAME uint16 = 0x0000
const T_PAYLOAD uint16 = 0x0001
const T_KEYID_REST uint16 = 0x0002
const T_HASH_REST uint16 = 0x0003
const T_PAYLDTYPE uint16 = 0x0005
const T_EXPIRY uint16 = 0x0006
const T_HASHGROUP uint16 = 0x0007
const T_KEX uint16 = 0x0009

// HashGroup containers
const T_POINTER_SIZE uint16 = 0x0001
const T_DATA_POINTER uint16 = 0x0002
const T_MANIFEST_POINTER uint16 = 0x0003

// name sements
const T_NAMESEG_NAME uint16 = 0x0001
const T_NAMESEG_IPID uint16 = 0x0002
const T_NAMESEG_CHUNK uint16 = 0x0010
const T_NAMESEG_VERSION uint16 = 0x0013
const T_NAMESEG_APP0 uint16 = 0x1000
const T_NAMESEG_APP1 uint16 = 0x1001
const T_NAMESEG_APP2 uint16 = 0x1002
const T_NAMESEG_APP3 uint16 = 0x1003
const T_NAMESEG_APP4 uint16 = 0x1004

// KEX message types
const T_KEX_BAREHELLO uint16 = 0x0001
const T_KEX_REJECT uint16 = 0x0002
const T_KEX_HELLO uint16 = 0x0003
const T_KEX_ACCEPT uint16 = 0x0004

// KEX extension types
const T_KEX_MESSAGE_TYPE uint16 = 0x0000
const T_KEX_SOURCE_CHALLENGE uint16 = 0x0001
const T_KEX_SOURCE_TOKEN uint16 = 0x0002
const T_KEX_SOURCE_PROOF uint16 = 0x0003
const T_KEX_MOVE_CHALLENGE uint16 = 0x0004
const T_KEX_MOVE_TOKEN uint16 = 0x0005
const T_KEX_MOVE_PROOF uint16 = 0x0006
const T_KEX_SESSION_ID uint16 = 0x0007
const T_KEX_TIMESTAMP uint16 = 0x0008
const T_KEX_KEYSHARE uint16 = 0x0009

// Payload type
const T_PAYLOADTYPE_DATA uint16 = 0x0000
const T_PAYLOADTYPE_KEY uint16 = 0x0001
const T_PAYLOADTYPE_LINK uint16 = 0x0002
const T_PAYLOADTYPE_MANIFEST uint16 = 0x0003
const T_PAYLOADTYPE_ENCAP uint16 = 0x0004

// Validation fields - crypto suites
const T_CRC32C uint16 = 0x0002
const T_HMAC_SHA256 uint16 = 0x0003
const T_RSA_SHA256 uint16 = 0x0006

// Validation fields - validation dependent data
const T_KEYID uint16 = 0x0009
const T_PUBLICKEY uint16 = 0x000B
const T_SIGTIME uint16 = 0x000F

// TODO
const T_HASH uint16 = 0x0001
const T_HASH_SHA256 uint16 = 0x0001
