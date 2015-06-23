# contentService
Go server to deal with user media content like images, videos, etc...

Works with MongoDB backend and users logged in you server using mngo-sessions (you should have a collection caled session in order to check if user is logged in or not)

### Example of usage:
```
package main

import (
  "github.com/MGrin/ContentService"
)

const (
  PORT = 7896
  DB = "mongodb://localhost:27017"
  DB_NAME = "Eventorio-dev"
  ORIG_PATH = "./pictures"
)
func main() {
  var service, err = ContentService.Create(ORIG_PATH, DB, DB_NAME)
  if err != nil {
    panic(err)
  }

  service.Start(PORT)
}
```
