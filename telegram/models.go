package telegram

import (
	"time"

	"github.com/Laisky/laisky-blog-graphql/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

// AlertType type of alert
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
	ChatID     int64         `bson:"chat_id" json:"chat_id"`
}

func (db *MonitorDB) GetAlertTypesCol() *mgo.Collection {
	return db.dbcli.GetCol(alertTypeColName)
}
func (db *MonitorDB) GetUsersCol() *mgo.Collection {
	return db.dbcli.GetCol(usersColName)
}
