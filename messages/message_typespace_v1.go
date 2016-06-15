packae messages

import "fmt"

// hop-by-hop headers
const T_INT_LIFE int = 0x0001
const T_CACHE_TIME int = 0x0002
const T_MSG_HASH int = 0x0003

// top-level TLVs
const T_INTEREST int = 0x0001
const T_OBJECT int = 0x0002
const T_VALALG int = 0x0003
const T_VALSIG int = 0x0004
const T_MANIFEST int = 0x0006

// Message body
const T_NAME int = 0x0000
const T_PAYLOAD int = 0x0001
const T_KEYID_REST int = 0x0002
const T_HASH_REST int = 0x0003
const T_PAYLDTYPE int = 0x0005
const T_EXPIRY int = 0x0006
const T_HASHGROUP int = 0x0007
const T_BLOCKHASHGROUP int = 0x0008

  // name sements
const T_NAMESEG_NAME int = 0x0001
const T_NAMESEG_IPID int = 0x0002
const T_NAMESEG_CHUNK int = 0x0010
const T_NAMESEG_VERSION int = 0x0013

const T_NAMESEG_APP0 int = 0x1000
const T_NAMESEG_APP1 int = 0x1001
const T_NAMESEG_APP2 int = 0x1002
const T_NAMESEG_APP3 int = 0x1003
const T_NAMESEG_APP4 int = 0x1004

// Payload type
constT_PAYLOADTYPE_DATA int = 0x00
const T_PAYLOADTYPE_KEY int = 0x01
const T_PAYLOADTYPE_LINK int = 0x02
const T_PAYLOADTYPE_MANIFEST int = 0x3

// Validation fields
const T_CRC32C int = 0x0002
const T_HMAC_SHA256 int = 0x0003
const T_RSA_SHA256 int = 0x0006

const T_KEYID int = 0x0009
const T_PUBLICKEY int = 0x000B
const T_SIGTIME int = 0x000F
