package gopwn

import(
  "bufio"
  "net"
  "os"
  "os/exec"
  "io"
)

type RemoteDial struct {
  conn    net.Conn
}

type LocalDial struct {
  cmd     *exec.Cmd
  stdin   io.WriteCloser
  stdout  io.ReadCloser
}

type Dial interface {
  Close() error
  Interactive(*bufio.Writer)
}

type Connection struct {
  dial    Dial
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
    dial: RemoteDial{conn},
    reader: bufio.NewReader(conn),
    writer: bufio.NewWriter(conn),
  }

  return c, nil
}

func Process(path string)(Connection, error){
  // exec
  printfLog("Starting process for %s...\n", path)
  cmd := exec.Command(path)

  // build reader, writer
  stdin, err := cmd.StdinPipe()
  if err != nil {
    return Connection{}, err
  }

  stdout, err := cmd.StdoutPipe()
  if err != nil {
    return Connection{}, err
  }

  // start
  err = cmd.Start()
  if err != nil {
    return Connection{}, err
  }

  printfLog("%s executed at pid %d.\n", path, cmd.Process.Pid)

  c := Connection{
    dial:   LocalDial{cmd, stdin, stdout},
    reader: bufio.NewReader(stdout),
    writer: bufio.NewWriter(stdin),
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

func (conn Connection) Close() error {
  return conn.dial.Close()
}

func (dial RemoteDial) Close() error {
  return dial.conn.Close()
}

func (dial LocalDial) Close() error {
  errIn := dial.stdin.Close()
  errOut := dial.stdout.Close()

  if errIn != nil {
    return errIn
  } else if errOut != nil {
    return errOut
  }

  return nil
}

func (conn Connection) Interactive() {
  printfLog("Switch to interactive mode.\n")
  go conn.dial.Interactive(conn.writer)
  conn.reader.WriteTo(os.Stdout)
}

func (dial RemoteDial) Interactive(w *bufio.Writer) {
  w.ReadFrom(os.Stdin)
}

func (dial LocalDial) Interactive(w *bufio.Writer) {
  io.Copy(dial.stdin, os.Stdin)
}

