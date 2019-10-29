package telegram

import (
	"time"

	"github.com/Laisky/go-utils"

	"github.com/Laisky/laisky-blog-graphql/models"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	monitorDBName    = "monitor"
	alertTypeColName = "alert_types"
	usersColName     = "users"
)

// MonitorDB db
type MonitorDB struct {
	dbcli *models.DB
}

// NewMonitorDB create new MonitorDB
func NewMonitorDB(dbcli *models.DB) *MonitorDB {
	return &MonitorDB{
		dbcli: dbcli,
	}
}

// AlertTypes type of alert
type AlertTypes struct {
	ID         bson.ObjectId `bson:"_id,omitempty" json:"mongo_id"`
	Name       string        `bson:"name" json:"name"`
	PushToken  string        `bson:"push_token" json:"push_token"`
	JoinKey    string        `bson:"join_key" json:"join_key"`
	CreatedAt  time.Time     `bson:"created_at" json:"created_at"`
	ModifiedAt time.Time     `bson:"modified_at" json:"modified_at"`
}

type Users struct {
	ID         bson.ObjectId `bson:"_id,omitempty" json:"mongo_id"`
	CreatedAt  time.Time     `bson:"created_at" json:"created_at"`
	ModifiedAt time.Time     `bson:"modified_at" json:"modified_at"`
	Name       string        `bson:"name" json:"name"`
	UID        int           `bson:"uid" json:"uid"`
}

func (db *MonitorDB) GetAlertTypesCol() *mgo.Collection {
	return db.dbcli.GetCol(alertTypeColName)
}
func (db *MonitorDB) GetUsersCol() *mgo.Collection {
	return db.dbcli.GetCol(usersColName)
}

func (db *MonitorDB) CreateOrGetUser(user *tb.User) (u *Users, err error) {
	u = new(Users)
	if err := db.GetUsersCol().Find(bson.M{"uid": user.ID}).One(u); err == mgo.ErrNotFound {
		u.CreatedAt = utils.Clock.GetUTCNow()
		u.ModifiedAt = utils.Clock.GetUTCNow()
		u.Name = user.Username
		u.UID = user.ID
		if err = db.GetUsersCol().Insert(u); err != nil {
			return nil, errors.Wrap(err, "insert new user")
		}
	} else if err != nil {
		return nil, errors.Wrap(err, "load user from db")
	}

	return u, nil
}

func (db *MonitorDB) CreateAlertType(name string) (at *AlertTypes, err error) {
	// check if exists
	if _, err = db.GetAlertTypesCol().Upsert(
		bson.M{"name": name},
		bson.M{"$setOnInsert": bson.M{
			"name":        name,
			"push_token":  utils.RandomStringWithLength(20),
			"join_key":    utils.RandomStringWithLength(6),
			"created_at":  utils.Clock.GetUTCNow(),
			"modified_at": utils.Clock.GetUTCNow(),
		}},
	); err != nil {
		return nil, errors.Wrap(err, "upsert docu")
	}

	at = new(AlertTypes)
	if err = db.GetAlertTypesCol().Find(bson.M{"name": name}).One(at); err != nil {
		return nil, errors.Wrap(err, "load alert_types")
	}

	return at, nil
}
