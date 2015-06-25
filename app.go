package ContentService

import (
  "github.com/gorilla/mux"

  "net/http"
  "mime/multipart"
  "encoding/json"
  "errors"
  "fmt"
  "strings"
  "os"
  "path"
)

type ContentService struct {
  r *mux.Router
  mongo *MongoDBBackend
}

var ORIG_PATH string

func Create(originalsPath string, mongoURL string, mongoDBName string) (service *ContentService, err error){
  ORIG_PATH = originalsPath

  // Create folder for pictures
  err = os.MkdirAll(path.Join(ORIG_PATH), 0777)
  if err != nil {
    return service, err
  }

  // Creating the service to store user content
  service = &ContentService{}

  service.r = mux.NewRouter()
  if err != nil {
    return service, err
  }

  service.r.HandleFunc("/{itemId}/{pictureType}", service.HandleUploadPicture).Methods("POST")
  service.r.HandleFunc("/{itemId}/{pictureType}/{pictureName}", service.HandlDeletePicture).Methods("DELETE", "OPTIONS")
  service.r.HandleFunc("/{itemId}/{pictureType}/{pictureName}", service.HandlePictureConfirmation).Methods("PUT")
  service.r.HandleFunc("/{itemId}/{pictureName}", service.HandlePictureRequest).Methods("GET")

  service.mongo, err = CreateMongoBackend(mongoURL, mongoDBName)

  return service, err
}

func (service *ContentService) Start(port int) error{
  http.Handle("/", service.r)
  fmt.Println("Without access-control-origin")
  fmt.Printf("Pictures Server running on port %d\n", port)
  http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

  return nil
}

func (service *ContentService) HandleUploadPicture(rw http.ResponseWriter, req *http.Request) {
  var err error

  err = SetHeaders(rw, req)
  if err != nil {
    HandleError(500, err, rw)
    return
  }

  vars := mux.Vars(req)
  itemId := vars["itemId"]
  pictureType := vars["pictureType"]

  if !service.SupportsType(pictureType) {
    HandleError(404, errors.New("Not found"), rw)
    return
  }

  err = service.Authorize(req)
  if err != nil {
    HandleError(403, err, rw)
    return
  }

  var itemType string

  itemType, err = service.mongo.GetItemType(itemId)
  if err != nil {
    HandleError(500, err, rw)
  }

  if itemType == TEMP_TYPE{
    itemId = "temp"
  }

  err = os.MkdirAll(path.Join(ORIG_PATH, itemId), 0777)
  if err != nil {
    HandleError(500, err, rw)
    return
  }

  var (
    file multipart.File
    fileHeader *multipart.FileHeader
  )

  file, fileHeader, err = req.FormFile("picture")
  if err != nil {
    HandleError(500, err, rw)
    return
  }
  if file == nil {
    HandleError(500, errors.New("No picture provided"), rw)
    return
  }
  defer file.Close()

  var filename string
  filename, err = service.UploadPicture(file, fileHeader, itemId, pictureType)
  if err != nil {
    HandleError(500, err, rw)
    return
  }

  rw.WriteHeader(200)
  rw.Write([]byte(filename))
}

func (service *ContentService) HandlDeletePicture(rw http.ResponseWriter, req *http.Request) {
  var err error

  err = SetHeaders(rw, req)
  if err != nil {
    HandleError(500, err, rw)
    return
  }

  if req.Method == "OPTIONS" {
    rw.WriteHeader(200)
    return
  }
  vars := mux.Vars(req)
  itemId := vars["itemId"]
  pictureType := vars["pictureType"]
  pictureName := vars["pictureName"]

  err = service.RemovePicture(itemId, pictureType, pictureName)
  if err != nil {
    HandleError(500, err, rw)
    return
  }

  rw.WriteHeader(200)
}

func (service *ContentService) HandlePictureConfirmation(rw http.ResponseWriter, req *http.Request) {
  var err error

  err = SetHeaders(rw, req)
  if err != nil {
    HandleError(500, err, rw)
    return
  }

  vars := mux.Vars(req)
  itemId := vars["itemId"]
  pictureType := vars["pictureType"]
  pictureName := vars["pictureName"]

  err = service.ConfirmPicture(itemId, pictureType, pictureName)
  if err != nil {
    HandleError(500, err, rw)
    return
  }

  rw.WriteHeader(200)
}

func (service *ContentService) HandlePictureRequest(rw http.ResponseWriter, req *http.Request) {
  err := SetHeaders(rw, req)
  if err != nil {
    HandleError(500, err, rw)
    return
  }

  vars := mux.Vars(req)
  itemId := vars["itemId"]
  pictureName := vars["pictureName"]

  http.ServeFile(rw, req, path.Join(ORIG_PATH, itemId, pictureName))
}

func HandleError(code int, err error, rw http.ResponseWriter) {
  fmt.Println(err.Error());
  rw.WriteHeader(code)
  rw.Write([]byte(err.Error()))
}

func SetHeaders(rw http.ResponseWriter, req *http.Request) error{
  rw.Header().Set("Access-Control-Allow-Origin", req.Header.Get("Origin"))
  rw.Header().Set("Access-Control-Allow-Credentials", "true")
  rw.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
  rw.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST, GET, PUT, DELETE")

  return nil
}

func (service *ContentService) Authorize(req *http.Request) error {
  var err error
  var sessionId string

  sessionId = req.FormValue("connect.sid")
  if sessionId == "" {
    return errors.New("Not authorized")
  }
  sessionId = strings.Split(sessionId, ".")[0]
  sessionId = sessionId[2:len(sessionId)]

  var sessionStruct *SessionModel
  sessionStruct, err = service.mongo.FindSessionById(sessionId)
  if err != nil {
    return err
  }

  var session map[string]map[string]interface{}

  err = json.Unmarshal([]byte(sessionStruct.Session), &session)
  if err != nil {
    return err
  }

  sessionUser := session["passport"]["user"]
  if sessionUser == nil {
    return errors.New("Not authorized")
  }
  userId := sessionUser.(string)
  var count int

  count, err = service.mongo.UsersCount(userId)
  if err != nil {
    return err
  }
  if count == 0 {
    return errors.New("Not authorized")
  }

  return nil
}
