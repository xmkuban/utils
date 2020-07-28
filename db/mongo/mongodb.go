package mongo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/x/mongo/driver/auth"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/xmkuban/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbMongoMap map[string]*MongoDBWrap

func init() {
	dbMongoMap = make(map[string]*MongoDBWrap)
}

type MongoDBWrap struct {
	Host           string
	Username       string
	Password       string
	AuthMechanisma string
	Database       string

	ConnectTimeOut int64

	MinPoolSize     uint64
	MaxPoolSize     uint64
	MaxConnLifetime int64

	client *mongo.Client
	d      *mongo.Database
}

func (db *MongoDBWrap) getConnString() string {
	if db.Host == "" {
		db.Host = "localhost:27017"
	}
	hostArr := strings.Split(db.Host, ":")
	if len(hostArr) == 1 {
		db.Host = fmt.Sprintf("%s:27017", db.Host)
	}
	if db.Username == "" {
		return fmt.Sprintf("mongodb://%s", db.Host)
	}
	if db.AuthMechanisma == "" {
		db.AuthMechanisma = auth.SCRAMSHA1
	}
	return fmt.Sprintf("mongodb://%s:%s@%s/%s?authMechanism=%s", db.Username, db.Password, db.Host, db.Database, db.AuthMechanisma)
}
func (db *MongoDBWrap) init() error {
	link := db.getConnString()
	opts := options.Client().ApplyURI(link)
	if db.MinPoolSize > 0 {
		opts.SetMinPoolSize(db.MinPoolSize)
	}
	if db.MaxPoolSize > 0 {
		opts.SetMaxPoolSize(db.MaxPoolSize)
	}
	if db.MaxConnLifetime > 0 {
		opts.SetMaxConnIdleTime(time.Duration(db.MaxConnLifetime) * time.Second)
	}
	if db.ConnectTimeOut > 0 {
		opts.SetConnectTimeout(time.Duration(db.ConnectTimeOut) * time.Second)
	}
	client, err := mongo.NewClient(opts)
	if err != nil {
		logger.Error(err)
		return err
	}
	err = client.Connect(context.Background())
	if err != nil {
		logger.Error(err)
		return err
	}
	db.d = client.Database(db.Database)
	return nil
}

func (db *MongoDBWrap) DB() *mongo.Database {
	if db.d != nil {
		return db.d
	}
	return nil
}

func (db *MongoDBWrap) Table(name string) *Collection {
	collection := db.DB().Collection(name)
	return &Collection{
		Collection: collection,
		filter:     bson.M{},
	}
}

func (db *MongoDBWrap) Close() {
	if db.client == nil {
		return
	}
	err := db.client.Disconnect(context.Background())
	if err != nil {
		logger.Error(err)
	}
}

func Mo(key string) *MongoDBWrap {
	db, ok := dbMongoMap[key]
	if !ok {
		return nil
	}

	return db
}

func InitMongoDB(key string, wrap *MongoDBWrap) error {
	err := wrap.init()
	if err != nil {
		return err
	}
	dbMongoMap[key] = wrap
	return nil
}

type Collection struct {
	*mongo.Collection
	filter bson.M
	sort   bson.M
	limit  int64
	skip   int64
}

func (c *Collection) Where(key string, value interface{}) *Collection {
	c.filter[key] = value
	return c
}

//$ne相当与"！"
func (c *Collection) NoEqual(key string, value interface{}) *Collection {
	if _, ok := c.filter[key]; ok {
		c.filter[key].(bson.M)["$ne"] = value
		return c
	}
	c.filter[key] = bson.M{
		"$ne": value,
	}
	return c
}

//$gt相当于">"
func (c *Collection) MoreThan(key string, value interface{}) *Collection {
	if _, ok := c.filter[key]; ok {
		c.filter[key].(bson.M)["$gt"] = value
		return c
	}
	c.filter[key] = bson.M{
		"$gt": value,
	}
	return c
}

//$lt相当与"<"
func (c *Collection) LessThan(key string, value interface{}) *Collection {
	if _, ok := c.filter[key]; ok {
		c.filter[key].(bson.M)["$lt"] = value
		return c
	}
	c.filter[key] = bson.M{
		"$lt": value,
	}
	return c
}

//$gte相当于">="
func (c *Collection) MoreThanAndEqual(key string, value interface{}) *Collection {
	if _, ok := c.filter[key]; ok {
		c.filter[key].(bson.M)["$gte"] = value
		return c
	}
	c.filter[key] = bson.M{
		"$gte": value,
	}
	return c
}

//$lte相当于"<="
func (c *Collection) LessThanAndEqual(key string, value interface{}) *Collection {
	if _, ok := c.filter[key]; ok {
		c.filter[key].(bson.M)["$lte"] = value
		return c
	}
	c.filter[key] = bson.M{
		"$lte": value,
	}
	return c
}

//$in在集合内
func (c *Collection) In(key string, value interface{}) *Collection {
	if _, ok := c.filter[key]; ok {
		c.filter[key].(bson.M)["$in"] = value
		return c
	}
	c.filter[key] = bson.M{
		"$in": value,
	}
	return c
}

//$nin不在集合内
func (c *Collection) NotIn(key string, value interface{}) *Collection {
	if _, ok := c.filter[key]; ok {
		c.filter[key].(bson.M)["$nin"] = value
		return c
	}
	c.filter[key] = bson.M{
		"$nin": value,
	}
	return c
}

//$exists是否包含键(isExist=true,key键存在的所有记录)
func (c *Collection) Exists(key string, isExist bool) *Collection {
	if _, ok := c.filter[key]; ok {
		c.filter[key].(bson.M)["$exists"] = isExist
		return c
	}
	c.filter[key] = bson.M{
		"$exists": isExist,
	}
	return c
}

//键值为null（键存在，键值为null）
func (c *Collection) ValueNullByKey(key string) *Collection {
	c.filter[key] = bson.M{
		"$in":     []interface{}{nil},
		"$exists": true,
	}
	return c
}

func (c *Collection) SortAsc(key string) *Collection {
	if c.sort == nil {
		c.sort = bson.M{}
	}
	c.sort[key] = 1
	return c
}

func (c *Collection) SortDesc(key string) *Collection {
	if c.sort == nil {
		c.sort = bson.M{}
	}
	c.sort[key] = -1
	return c
}

func (c *Collection) Limit(num int64) *Collection {
	c.limit = num
	return c
}

func (c *Collection) Skip(num int64) *Collection {
	c.skip = num
	return c
}

func (c *Collection) Insert(document interface{}) (*mongo.InsertOneResult, error) {
	return c.Collection.InsertOne(context.Background(), document)
}

func (c *Collection) InsertMany(documents []interface{}) (*mongo.InsertManyResult, error) {
	return c.Collection.InsertMany(context.Background(), documents)
}

func (c *Collection) Update(document interface{}) (*mongo.UpdateResult, error) {
	return c.Collection.UpdateMany(context.Background(), c.filter, document)
}

func (c *Collection) Get(result interface{}) (bool, error) {
	opts := new(options.FindOneOptions)
	if len(c.sort) != 0 {
		opts.Sort = c.sort
	}
	if c.skip > 0 {
		opts.Skip = &c.skip
	}
	err := c.Collection.FindOne(context.Background(), c.filter).Decode(result)
	if err == nil {
		return true, nil
	}
	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	return false, nil
}

func (c *Collection) Find(results interface{}) error {
	opts := new(options.FindOptions)
	if len(c.sort) != 0 {
		opts.Sort = c.sort
	}
	if c.limit > 0 {
		opts.Limit = &c.limit
	}
	if c.skip > 0 {
		opts.Skip = &c.skip
	}
	cursor, err := c.Collection.Find(context.Background(), c.filter, opts)
	if err != nil {
		return err
	}
	return cursor.All(context.Background(), results)
}

func (c *Collection) Delete() (*mongo.DeleteResult, error) {
	return c.Collection.DeleteMany(context.Background(), c.filter)
}
