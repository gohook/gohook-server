package mongo

import (
	"errors"
	"github.com/gohook/gohook-server/user"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const AccountDoc = "account"

type MongoAccountStore struct {
	db      string
	session *mgo.Session
}

func NewMongoAccountStore(db string, session *mgo.Session) user.AccountStore {
	return &MongoAccountStore{
		db:      db,
		session: session,
	}
}

func (d *MongoAccountStore) Add(u *user.Account) error {
	sess := d.session.Copy()
	defer sess.Close()

	c := sess.DB(d.db).C(AccountDoc)

	id := bson.NewObjectId()
	u.Id = user.AccountId(id.String())
	_, err := c.UpsertId(id, bson.M{"$set": u})
	return err
}

func (d *MongoAccountStore) Remove(id user.AccountId) (*user.Account, error) {
	sess := d.session.Copy()
	defer sess.Close()

	c := sess.DB(d.db).C(AccountDoc)

	err := c.Remove(bson.M{"_id": id})
	if err == mgo.ErrNotFound {
		return nil, errors.New("Not Found")
	}

	return nil, err
}

func (d *MongoAccountStore) Find(id user.AccountId) (*user.Account, error) {
	sess := d.session.Copy()
	defer sess.Close()

	c := sess.DB(d.db).C(AccountDoc)

	var result user.Account
	err := c.Find(bson.M{"_id": id}).One(&result)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, errors.New("Not Found")
		}
		return nil, err
	}

	return &result, nil
}

func (d *MongoAccountStore) FindByToken(token user.AccountToken) (*user.Account, error) {
	sess := d.session.Copy()
	defer sess.Close()

	c := sess.DB(d.db).C(AccountDoc)

	var result user.Account
	err := c.Find(bson.M{"token": token}).One(&result)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, errors.New("Not Found")
		}
		return nil, err
	}

	return &result, nil
}
