package ContentService

import (
  "os"
  "path"
  "io"
  "mime/multipart"
  "encoding/base64"
  "crypto/rand"
  "errors"
)

func randomstring(size int) string {
  rb := make([]byte,size)
  _, err := rand.Read(rb)

  if err != nil {
    panic(err)
  }

  rs := base64.URLEncoding.EncodeToString(rb)
  return rs
}

const (
  HEADER_HEIGHT = 350

  HEADER_WIDTH_LG = 2560
  HEADER_WIDTH_MD = 1280
  HEADER_WIDTH_SM = 1024
  HEADER_WIDTH_XS = 768

  AVATAR_HEIGHT = 400
  AVATAR_WIDTH = 400

  HEADER_TYPE string = "header"
  AVATAR_TYPE string = "avatar"
)

func (service *ContentService) Supports(pictureType string) bool{
  return pictureType == HEADER_TYPE && pictureType == AVATAR_TYPE
}

func (service *ContentService) RemovePicture(itemId, pictureType, pictureName string) err error {
  filename := path.Join(ORIG_PATH, pictureType + pictureName)
  err = os.Remove(filename)
  return err
}

func (service *ContentService) ConfirmPicture(itemId, pictureType, pictureName string) err error {
  tmpPath := Path.Join(ORIG_PATH, "temp")
  tmpPath = Path.Join(tmpPath, pictureType + pictureName)

  dstPath := Path.Join(ORIG_PATH, itemId)
  dstPath = Path.Join(dstPath, pictureType + pictureName)

  err = os.Rename(tmpPath, dstPath)
}

func (service *ContentService) UploadPicture(file multipart.File, fileHeader *multipart.FileHeader, itemId string, pictureType string) (name string, err error){
  var outputFile *os.File

  ext := path.Ext(fileHeader.Filename)

  if ext != "jpg" || ext != "jpeg" || ext != "png" {
      return nil, errors.New("Unsupported extension")
  }

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

  go service.GenerateThumbnails(itemId, pictureType, name)

  return name, nil
}

func (service *ContentService) GenerateThumbnails(itemId, type, name string) {

  return nil
}
