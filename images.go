package contentService

import (
  "os"
  "path"
  "io"
  "mime/multipart"
  "math/rand"
)

// Change to the real random!!!
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

const (
  HEADER_HEIGHT = 350

  HEADER_WIDTH_LG = 2560
  HEADER_WIDTH_MD = 1280
  HEADER_WIDTH_SM = 1024
  HEADER_WIDTH_XS = 768

  AVATAR_HEIGHT = 400
  AVATAR_WIDTH = 400
)

func (service *ContentService) UploadPicture(file multipart.File, fileHeader *multipart.FileHeader, itemId string, pictureType string) (name string, err error){
  var outputFile *os.File

  name = randSeq(15) + path.Ext(fileHeader.Filename)

  outputFilePath := path.Join(ORIG_PATH, itemId)
  outputFilePath = path.Join(outputFilePath, pictureType + name)

  outputFile, err = os.Create(outputFilePath)
  defer outputFile.Close()

  if err != nil {
    return "", err
  }

  _, err = io.Copy(outputFile, file)
  if err != nil {
    return "", err
  }

  // TODO thumbnails of diferent sizes

  return name, nil
}
