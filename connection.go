package gopwn

import(
  "bufio"
  "net"
  "os"
)

type Connection struct {
  url     string
  conn    net.Conn
  reader  *bufio.Reader
  writer  *bufio.Writer
}

func Remote(url string)(Connection, error){
  // connect
  printfLog("Connecting to %s...\n", url)
  conn, err := net.Dial("tcp", url)
  if err != nil {
    return Connection{}, err
  }
  printfLog("%s connected.\n", conn.RemoteAddr())

  // build reader, writer
  c := Connection{
    url:    url,
    conn:   conn,
    reader: bufio.NewReader(conn),
    writer: bufio.NewWriter(conn),
  }

  return c, nil
}

func (conn Connection) Recvuntil(delim byte) (string, error) {
  recv, err := conn.reader.ReadBytes(delim)
  return string(recv), err
}

func (conn Connection) Sendline(s string) error {
  _, err := conn.writer.WriteString(s + "\n")
  conn.writer.Flush()
  return err
}

func (conn Connection) Send(s string) error {
  _, err := conn.writer.WriteString(s)
  conn.writer.Flush()
  return err
}

func (conn Connection) RemoteAddr() string {
  return conn.url
}

func (conn Connection) Close() error {
  return conn.conn.Close()
}

func (conn Connection) Interactive() {
  printfLog("Switch to interactive mode.\n")
  go conn.writer.ReadFrom(os.Stdin)
  conn.reader.WriteTo(os.Stdout)
}
