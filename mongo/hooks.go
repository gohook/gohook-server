package mongo

import (
	"errors"
	"github.com/gohook/gohook-server/gohookd"
	"github.com/gohook/gohook-server/user"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const HookDoc = "hook"

type MongoHookStore struct {
	db        string
	session   *mgo.Session
	accountId user.AccountId
	scoped    bool
}

func NewMongoHookStore(db string, session *mgo.Session) (gohookd.HookStore, error) {
	d := &MongoHookStore{
		db:      db,
		session: session,
		scoped:  false,
	}

	index := mgo.Index{
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	sess := d.session.Copy()
	defer sess.Close()

	c := sess.DB(d.db).C(HookDoc)

	if err := c.EnsureIndex(index); err != nil {
		return nil, err
	}

	return d, nil
}

func (d MongoHookStore) Scope(accountId user.AccountId) gohookd.HookStore {
	d.accountId = accountId
	d.scoped = true
	return &d
}

func (d *MongoHookStore) Add(m *gohookd.Hook) error {
	sess := d.session.Copy()
	defer sess.Close()

	if d.scoped {
		m.AccountId = d.accountId
	}

	c := sess.DB(d.db).C(HookDoc)

	id := bson.NewObjectId()
	_, err := c.UpsertId(id, bson.M{"$set": m})
	return err
}

func (d *MongoHookStore) Find(id gohookd.HookID) (*gohookd.Hook, error) {
	sess := d.session.Copy()
	defer sess.Close()

	c := sess.DB(d.db).C(HookDoc)

	q := bson.M{"id": id}

	if d.scoped {
		q = bson.M{"id": id, "accountid": d.accountId}
	}

	var result gohookd.Hook
	err := c.Find(q).One(&result)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, errors.New("Not Found")
		}
		return nil, err
	}

	return &result, nil
}

func (d *MongoHookStore) FindAll() (gohookd.HookList, error) {
	sess := d.session.Copy()
	defer sess.Close()

	c := sess.DB(d.db).C(HookDoc)

	q := bson.M{}

	if d.scoped {
		q = bson.M{"accountid": d.accountId}
	}

	var result gohookd.HookList
	err := c.Find(q).All(&result)
	if err != nil {
		return gohookd.HookList{}, errors.New("Not Found")
	}

	return result, nil
}

func (d *MongoHookStore) Remove(id gohookd.HookID) (*gohookd.Hook, error) {
	sess := d.session.Copy()
	defer sess.Close()

	c := sess.DB(d.db).C(HookDoc)

	q := bson.M{"id": id}

	if d.scoped {
		q = bson.M{"id": id, "accountid": d.accountId}
	}

	err := c.Remove(q)
	if err == mgo.ErrNotFound {
		return nil, errors.New("Not Found")
	}

	return nil, err
}
