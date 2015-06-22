package ContentService

import (
  "github.com/gorilla/mux"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"

  "net/http"
  "mime/multipart"
  "encoding/json"
  "errors"
  "fmt"
  "strings"
  "os"
  "path"
)

var (
  PORT int
  DB string
  DB_NAME string
  ORIG_PATH string
)

type ContentService struct {
  r *mux.Router
  mongoSession *mgo.Session
  mongo *mgo.Database
}

// MongoDB Session that holds the information about logged in users
type Session struct {
  Id bson.ObjectId "_id,omitempty"
  Session string "session,omitempty"
}

func Create(originalsPath string, mongoURL string, mongoDBName string) (service *ContentService, err error){
  ORIG_PATH = originalsPath
  DB_NAME = mongoDBName
  DB = mongoURL

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
  service.r.PathPrefix("/").Handler(http.FileServer(http.Dir(ORIG_PATH))).Methods("GET")

  return service, err
}

func (service *ContentService) Start(port int) error{
  PORT = port
  
  var err error
  service.mongoSession, err = mgo.Dial(DB)
  if err != nil {
    return err
  }
  service.mongo = service.mongoSession.DB(DB_NAME)

  http.Handle("/", service.r)
  fmt.Printf("Pictures Server running on port %d\n", PORT)
  http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)

  return nil
}

func (service *ContentService) HandleUploadPicture(rw http.ResponseWriter, req *http.Request) {
  var err error

  vars := mux.Vars(req)
  itemId := vars["itemId"]
  pictureType := vars["pictureType"]

  err = service.Authorize(req)
  if err != nil {
    HandleError(403, err, rw)
    return
  }

  err = os.MkdirAll(path.Join(ORIG_PATH, itemId), 0777)
  if err != nil {
    HandleError(403, err, rw)
    return
  }

  var (
    file multipart.File
    fileHeader *multipart.FileHeader
  )

  file, fileHeader, err = req.FormFile("picture")
  defer file.Close()
  if err != nil {
    HandleError(500, err, rw)
    return
  }

  var filename string
  filename, err = service.UploadPicture(file, fileHeader, itemId, pictureType)
  if err != nil {
    HandleError(500, err, rw)
    return
  }

  rw.WriteHeader(200)
  rw.Write([]byte(filename))
}

func HandleError(code int, err error, rw http.ResponseWriter) {
  rw.WriteHeader(code)
  rw.Write([]byte(err.Error()))
}

func (service *ContentService) Authorize(req *http.Request) error {
  var err error
  var sessionCookie *http.Cookie
  var sessionId string

  sessionCookie, err = req.Cookie("connect.sid")
  if err != nil {
    return err
  }
  sessionId = strings.Split(sessionCookie.Value, ".")[0]
  sessionId = sessionId[4:len(sessionId)]

  c := service.mongo.C("sessions")
  sessionStruct := &Session{}

  err = c.Find(bson.M{"_id" : sessionId}).One(sessionStruct)
  if err != nil {
    return err
  }

  var session map[string]map[string]interface{}

  err = json.Unmarshal([]byte(sessionStruct.Session), &session)
  if err != nil {
    return err
  }

  userId := bson.ObjectIdHex(session["passport"]["user"].(string))

  c = service.mongo.C("users")

  var count int
  count, err = c.FindId(userId).Count()
  if err != nil {
    return err
  }
  if count == 0 {
    return errors.New("Not authorized")
  }

  return nil
}
