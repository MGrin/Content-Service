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

func (service *ContentService) SupportsType(pictureType string) bool{
  return pictureType == HEADER_TYPE || pictureType == AVATAR_TYPE
}

func (service *ContentService) SupportsExtension(ext string) bool {
  return ext == ".jpg" || ext == ".jpeg" || ext == ".png"
}

func (service *ContentService) RemovePicture(itemId, pictureType, pictureName string) error {
  var err error
  filename := path.Join(ORIG_PATH, itemId, pictureType + pictureName)
  err = os.Remove(filename)
  return err
}

func (service *ContentService) ConfirmPicture(itemId, pictureType, pictureName string) error {
  var err error

  err = os.MkdirAll(path.Join(ORIG_PATH, itemId), 0777)
  if err != nil {
    return err
  }

  tmpPath := path.Join(ORIG_PATH, "temp", pictureType + pictureName)
  dstPath := path.Join(ORIG_PATH, itemId, pictureType + pictureName)
  
  err = os.Rename(tmpPath, dstPath)
  return err
}

func (service *ContentService) UploadPicture(file multipart.File, fileHeader *multipart.FileHeader, itemId string, pictureType string) (name string, err error){
  var outputFile *os.File

  ext := path.Ext(fileHeader.Filename)

  if !service.SupportsExtension(ext) {
      return "", errors.New("Unsupported extension")
  }

  name = randomstring(15) + path.Ext(fileHeader.Filename)

  outputFilePath := path.Join(ORIG_PATH, itemId, pictureType + name)

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

func (service *ContentService) GenerateThumbnails(itemId string, itemType string, name string) {

  return
}
