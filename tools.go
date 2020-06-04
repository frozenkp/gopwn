package gopwn

import(
  "encoding/binary"
)

func P64(i int) string {
  bytes := make([]byte, 8)
  binary.LittleEndian.PutUint64(bytes, uint64(i))

  return string(bytes)
}

func P32(i int) string {
  bytes := make([]byte, 4)
  binary.LittleEndian.PutUint32(bytes, uint32(i))

  return string(bytes)
}

func U64(s string) int {
  bytes := make([]byte, 8)
  copy(bytes, []byte(s))
  return int(binary.LittleEndian.Uint64(bytes))
}

func U32(s string) int {
  bytes := make([]byte, 4)
  copy(bytes, []byte(s))
  return int(binary.LittleEndian.Uint32(bytes))
}
