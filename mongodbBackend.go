package ContentService

import (
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
 )
type MongoDBBackend struct {
	session *mgo.Session
	database *mgo.Database
}

const (
	SESSIONS_COLLECTION string = "sessions"
	USERS_COLLECTION string = "users"
	EVENTS_COLLECTION string = "events"

	USER_TYPE string = "user"
	EVENT_TYPE string = "event"
	TEMP_TYPE string = "temp"
)
// MongoDB Session that holds the information about logged in users
type SessionModel struct {
  Id bson.ObjectId "_id,omitempty"
  Session string "session,omitempty"
}

func CreateMongoBackend (url, name string) (MongoDBBackend, error){
	var err error
	var mongo := &MongoDBBackend{}

  mongo.session, err = mgo.Dial(url)
  if err != nil {
    return mongo, err
  }
  mongo.database = service.mongoSession.DB(name)
  return mongo, nil
}

func (mongo *MongoDBBackend) FindSessionById(id string) (*SessionModel, error) {
	var session := &SessionModel
	err := mongo.database.C(SESSIONS_COLLECTION).Find(bson.M{"_id" : id}).One(session)

	return *session, err
}

func (mongo *MongoDBBackend) UsersCount(id string) (int, error) {
	count, err := mongo.database.C(USERS_COLLECTION).FindId(bson.ObjectIdHex(id)).Count()
	return count, err
}

func (mongo *MongoDBBackend) GetItemType(id string) (string, error) {
	count, err := mongo.database.C(USERS_COLLECTION).FindId(bson.ObjectIdHex(id)).Count()
	if err != nil {
		return nil, err
	}
	if count != 0 {
		return USER_TYPE, nil
	}

	count, err := mongo.database.C(EVENTS_COLLECTION).FindId(bson.ObjectIdHex(id)).Count()
	if err != nil {
		return nil, err
	}
	if count != 0 {
		return EVENT_TYPE, nil
	}

	return TEMP_TYPE, nil
}