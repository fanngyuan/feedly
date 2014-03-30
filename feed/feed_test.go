package feed

import (
	"testing"
	mcstorage "github.com/fanngyuan/mcstorage"
	"reflect"
	"github.com/fanngyuan/feedly/activity"
)

type T struct {
	A int
}

func (this T)GetId()uint64{
	return uint64(this.A)
}

func (this T)GetType()string{
	return "T"
}

func newRedisListStorage()mcstorage.RedisListStorage{
	redisListStorage,_ := mcstorage.NewRedisListStorage(":6379", "test_list", 0, mcstorage.DecodeIntReversedSlice)
	return redisListStorage
}

func newRedisStorage()mcstorage.RedisStorage{
	tt := T{1}
	jsonEncoding:=mcstorage.JsonEncoding{reflect.TypeOf(&tt)}
	redisStorage,_ := mcstorage.NewRedisStorage(":6379", "test", 0, jsonEncoding)
	return redisStorage
}

func newRedisCounterStorage()mcstorage.RedisStorage{
	jsonEncoding:=mcstorage.JsonEncoding{reflect.TypeOf(1)}
	redisStorage,_ := mcstorage.NewRedisStorage(":6379", "test_count", 0, jsonEncoding)
	return redisStorage
}

func flush(storages []mcstorage.Storage){
	for _,storage:=range storages{
		storage.FlushAll()
	}

}

func TestAddRemoveActivities(t *testing.T) {
	redisStorage:=newRedisStorage()
	redisListStorage:=newRedisListStorage()
	redisCounterStorage:=newRedisCounterStorage()

	feed:=BaseFeed{redisStorage,redisListStorage,redisCounterStorage,"fanngyuan"}
	activity:=activity.Activity{uint64(1),"new note"}
	feed.AddActivity(activity)

	result:=feed.GetActivities(int64(0),int64(0),1,10)
	if len(result)!=1{
		t.Error("len should be 1")
	}
	if result[0].GetId()!=0{
		t.Error("id should be 0")
	}
	if result[0].GetType()!="T"{
		t.Error("Type should be T")
	}

}
