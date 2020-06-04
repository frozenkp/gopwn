package gopwn

import(
  "log"
)

var LOG bool = true

func printfLog(format string, v ...interface{}){
  if LOG {
    if v == nil {
      log.Printf(format)
    } else {
      log.Printf(format, v)
    }
  }
}
