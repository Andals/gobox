package mongo

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/andals/golog"

	"fmt"
)

type CmdLogFmtFunc func(cmd string, args ...interface{}) []byte

type Client struct {
	config *Config
	logger golog.ILogger
	clff   CmdLogFmtFunc

	conn      *mgo.Session
	db        *mgo.Database
	coll      *mgo.Collection
	connected bool

	pipeCnt int
}

func NewClient(config *Config, logger golog.ILogger) *Client {
	if config.LogLevel == 0 {
		config.LogLevel = golog.LEVEL_INFO
	}

	if logger == nil {
		logger = new(golog.NoopLogger)
	}

	this := &Client{
		config: config,
		logger: logger,
	}
	this.clff = this.cmdLogFmt
	this.SetLogger(logger)

	return this
}

func (this *Client) SetLogger(logger golog.ILogger) *Client {
	this.logger = logger
	mgo.SetLogger(NewMongoLogger(logger))

	return this
}

func (this *Client) SetDebug(debug bool) {
	mgo.SetDebug(debug)
}

func (this *Client) SetCmdLogFmtFunc(clff CmdLogFmtFunc) *Client {
	this.clff = clff

	return this
}

func (this *Client) Connected() bool {
	return this.connected
}

func (this *Client) Free() {
	if this.conn != nil {
		this.conn.Close()
	}

	this.connected = false
}

func (this *Client) Connect() error {
	url := "mongodb://" + this.config.User + ":" + this.config.Pass + "@" + this.config.Host + ":" + this.config.Port
	session, err := mgo.Dial(url)
	if err != nil {
		return err
	}

	//session.SetMode(mgo.Monotonic, true)
	session.SetMode(mgo.Eventual, true)

	this.conn = session
	this.db = session.DB(this.config.DBName)
	this.connected = true

	return nil
}

func (this *Client) connectCheck() {
	if !this.Connected() {
		if err := this.Connect(); err != nil {
			panic(err)
		}
	}
}

func (this *Client) DB(name string) *mgo.Database {
	database := this.conn.DB(name)
	this.db = database
	return database
}

func (this *Client) Collection(col string) *mgo.Collection {
	this.connectCheck()
	collection := this.db.C(col)
	this.coll = collection
	return collection
}

func (this *Client) Count() (n int, err error) {
	n, err = this.coll.Count()
	if err != nil {
		return 0, err
	}
	return n, err
}

func (this *Client) BuildQuery(query *Query) *mgo.Query {
	q := this.coll.Find(query.query)
	if query.selector != nil {
		q = q.Select(query.selector)
	}
	if query.limit != 0 {
		q = q.Limit(query.limit)
	}
	if query.skip != 0 {
		q = q.Skip(query.skip)
	}
	if query.setMaxTime != 0 {
		q = q.SetMaxTime(query.setMaxTime)
	}

	return q
}

func (this *Client) Query(query *Query) (result []bson.M, err error) {
	this.connectCheck()
	err = this.BuildQuery(query).All(&result)
	if err != nil {
		this.log("Query Fail, Query:", query,
			", Error:", err)
	}
	return result, err
}

func (this *Client) QueryOne(query *Query) (result bson.M, err error) {
	this.connectCheck()
	err = this.BuildQuery(query).One(&result)
	if err != nil {
		this.log("QueryOne Fail, Query:", query,
			", Error:", err)
	}
	return result, err
}

func (this *Client) QueryId(id interface{}) (result bson.M, err error) {
	this.connectCheck()
	err = this.coll.FindId(id).One(&result)
	if err != nil {
		this.log("QueryId Fail, Id:", id,
			", Error:", err)
	}
	return result, err
}

func (this *Client) QueryCount(query *Query) (n int, err error) {
	this.connectCheck()
	n, err = this.BuildQuery(query).Count()
	if err != nil {
		this.log("QueryCount Fail, Query:", query,
			", Error:", err)
	}
	return n, err
}

func (this *Client) Find(query interface{}) *mgo.Query {
	this.connectCheck()
	return this.coll.Find(query)
}

func (this *Client) FindId(id interface{}) *mgo.Query {
	this.connectCheck()
	return this.coll.FindId(id)
}

func (this *Client) Indexes() (indexes []mgo.Index, err error) {
	this.connectCheck()
	indexes, err = this.coll.Indexes()
	if err != nil {
		this.log("Indexes Fail, Indexes:", indexes,
			", Error:", err)
	}
	return indexes, err
}

func (this *Client) Insert(docs ...interface{}) error {
	this.connectCheck()
	err := this.coll.Insert(docs...)
	if err != nil {
		this.log("Insert Fail, docs:", docs,
			", Error:", err)
	}
	return err
}

func (this *Client) Update(selector, updater interface{}) error {
	this.connectCheck()
	err := this.coll.Update(selector, updater)
	if err != nil {
		this.log("Update Fail, Selector:", selector,
			", Updater:", updater,
			", Error:", err)
	}
	return err
}

func (this *Client) UpdateAll(selector, updater interface{}) error {
	this.connectCheck()
	_, err := this.coll.UpdateAll(selector, updater)
	if err != nil {
		this.log("UpdateAll Fail, Selector:", selector,
			", Updater:", updater,
			", Error:", err)
	}
	return err
}

func (this *Client) UpdateId(id interface{}, updater interface{}) error {
	this.connectCheck()
	err := this.coll.UpdateId(id, updater)
	if err != nil {
		this.log("UpdateId Fail, Id:", id,
			", Updater:", updater,
			", Error:", err)
	}
	return err
}

func (this *Client) Upsert(selector, updater interface{}) error {
	this.connectCheck()
	_, err := this.coll.Upsert(selector, updater)
	if err != nil {
		this.log("Upsert Fail, Selector:", selector,
			", Updater:", updater,
			", Error:", err)
	}
	return err
}

func (this *Client) Remove(selector interface{}) error {
	this.connectCheck()
	err := this.coll.Remove(selector)
	if err != nil {
		this.log("Remove Fail, Selector:", selector,
			", Error:", err)
	}
	return err
}

func (this *Client) RemoveAll(selector interface{}) error {
	this.connectCheck()
	_, err := this.coll.RemoveAll(selector)
	if err != nil {
		this.log("RemoveAll Fail, Selector:", selector,
			", Error:", err)
	}
	return err
}

func (this *Client) RemoveId(id interface{}) error {
	this.connectCheck()
	err := this.coll.RemoveId(id)
	if err != nil {
		this.log("RemoveId Fail, Id:", id,
			", Error:", err)
	}
	return err
}

func (this *Client) log(cmd string, args ...interface{}) {
	if len(cmd) == 0 {
		return
	}

	msg := this.clff(cmd, args...)
	if msg != nil {
		this.logger.Log(this.config.LogLevel, msg)
	}
}

func (this *Client) cmdLogFmt(cmd string, args ...interface{}) []byte {
	for _, arg := range args {
		cmd += " " + fmt.Sprint(arg)
	}

	return []byte(cmd)
}
