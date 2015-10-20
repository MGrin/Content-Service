# contentService
Go server to deal with user images.

Works with MongoDB backend and users logged in you server using mongo-sessions (you should have a collection called sessions in order to check if user is logged in or not)

### Example of usage:
```go
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
### Supported Paths
* `POST /{itemId}/{pictureType}` - upload a picture of type {pictureType} (now avatar or header are supported) to a ORIG_PATH/itemId/headerRandomString.extension. PNG of JPG formats are supported. If no itemId exists (for examle the user is creating this item and it does not exist yet in MongoDB) the itemId should be `temp`.
* `DELETE /{itemId}/{pictureType}/{pictureName}` - remove a picture of type pictureType (header or avatar) with name pictureName (a random string) from a folder ORIG_PATH/itemId/
* `PUT /{itemId}/{pictureType}/{pictureName}` - Move a picture of type pictureType and name pictureName from the ORIG_PATH/temp folder to the ORIG_PATH/itemId folder. Is used when you are creating pictures for an unexisting items, and then you need to move them to the item-related folder. THe temp folder can be cleaned up every n minutes (the script will be added to this repo afterwards)
* `GET /{itemId}/{pictureName}` - Return the requested picture. In this case, the pictureName should have a form [header|avatar]pictureName.ext where pictureName is a random string returned by the POST request
