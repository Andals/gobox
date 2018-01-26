package mongo

import (
	"encoding/json"
	"time"

	"github.com/andals/golog"
	"gopkg.in/mgo.v2/bson"

	"testing"
)

func getTestClient() *Client {
	w, _ := golog.NewFileWriter("/tmp/test_mongo.log")
	logger, _ := golog.NewSimpleLogger(w, golog.LEVEL_INFO, golog.NewSimpleFormater())

	config := NewConfig("localhost", "myport", "myuser", "mypass", "mydb")
	config.LogLevel = golog.LEVEL_DEBUG

	return NewClient(config, logger)
}

func getTestCollection(client *Client) *Client {
	client.Collection("mycoll")
	return client
}

var client *Client

func init() {
	client = getTestClient()
	client = getTestCollection(client)

	//client.Free()
}

func TestInsert(t *testing.T) {
	var err error
	err = client.Insert(bson.M{"_id": 1, "a": 1, "b": 2})
	if err != nil {
		t.Error(err)
	}
	err = client.Insert(bson.M{"_id": 2, "a": 3, "b": 4})
	if err != nil {
		t.Error(err)
	}
	err = client.Insert(bson.M{"_id": 3, "a": 4, "b": 5})
	if err != nil {
		t.Error(err)
	}
}

func TestUpdate(t *testing.T) {
	selector := bson.M{"_id": 1}
	updater := bson.M{
		"$inc":         bson.M{"view_count": 1},
		"$currentDate": bson.M{"edit_time": true},
	}
	err := client.Update(selector, updater)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateAll(t *testing.T) {
	selector := bson.M{"_id": bson.M{"$gt": 0}}
	updater := bson.M{
		"$inc":         bson.M{"view_count": 1},
		"$currentDate": bson.M{"edit_time": true},
	}
	err := client.UpdateAll(selector, updater)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateId(t *testing.T) {
	id := 1
	updater := bson.M{
		"$inc":         bson.M{"view_count": 1},
		"$currentDate": bson.M{"edit_time": true},
	}
	err := client.UpdateId(id, updater)
	if err != nil {
		t.Error(err)
	}
}

func TestUpsert(t *testing.T) {
	selector := bson.M{"_id": 4}
	updater := bson.M{
		"$inc":         bson.M{"view_count": 1},
		"$currentDate": bson.M{"edit_time": true},
		"$setOnInsert": bson.M{"add_time": "2018-06-23 09:00:00"},
	}
	err := client.Upsert(selector, updater)
	if err != nil {
		t.Error(err)
	}
}

func TestQuery(t *testing.T) {
	query := NewQuery().SetMaxTime(1 * time.Second)
	result, err := client.Query(query)
	if err != nil {
		t.Error(err)
	}
	jsonData, _ := json.Marshal(result)
	t.Logf("%s", jsonData)
}

func TestQueryOne(t *testing.T) {
	query := NewQuery().Query(bson.M{"_id": bson.M{"$gt": 0}}).Select(bson.M{"_id": 0}).Skip(1)
	result, err := client.QueryOne(query)
	if err != nil {
		t.Error(err)
	}
	jsonData, _ := json.Marshal(result)
	t.Logf("%s", jsonData)
}

func TestQurtyId(t *testing.T) {
	result, err := client.QueryId(4)
	if err != nil {
		t.Error(err)
	}
	jsonData, _ := json.Marshal(result)
	t.Logf("%s", jsonData)
}

func TestQurtyCount(t *testing.T) {
	query := NewQuery()
	result, err := client.QueryCount(query)
	if err != nil {
		t.Error(err)
	}
	jsonData, _ := json.Marshal(result)
	t.Logf("%s", jsonData)
}

func TestFind(t *testing.T) {
	result := []bson.M{}
	err := client.Find(bson.M{"_id": bson.M{"$gt": 0}}).All(&result)
	if err != nil {
		t.Error(err)
	}
	jsonData, _ := json.Marshal(result)
	t.Logf("%s", jsonData)
}

func TestFindId(t *testing.T) {
	result := bson.M{}
	err := client.FindId(4).One(&result)
	if err != nil {
		t.Error(err)
	}
	jsonData, _ := json.Marshal(result)
	t.Logf("%s", jsonData)
}

func TestRemove(t *testing.T) {
	selector := bson.M{"_id": 4}
	err := client.Remove(selector)
	if err != nil {
		t.Error(err)
	}
}

func TestRemoveAll(t *testing.T) {
	selector := bson.M{"_id": bson.M{"$gt": 1}}
	err := client.RemoveAll(selector)
	if err != nil {
		t.Error(err)
	}
}

func TestRemoveId(t *testing.T) {
	id := 1
	err := client.RemoveId(id)
	if err != nil {
		t.Error(err)
	}
}
