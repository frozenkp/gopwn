package gopwn

import(
  "log"
)

var LOG bool = true

func printfLog(format string, v ...interface{}){
  if LOG {
    log.Printf(format, v...)
  }
}
