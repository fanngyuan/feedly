package feed

import (
	"testing"
	mcstorage "github.com/fanngyuan/mcstorage"
	"reflect"
	"github.com/fanngyuan/feedly/activity"
)

func newRedisListStorage()mcstorage.RedisListStorage{
	redisListStorage,_ := mcstorage.NewRedisListStorage(":6379", "test_list", 0, mcstorage.DecodeIntReversedSlice)
	return redisListStorage
}

func newRedisStorage()mcstorage.RedisStorage{
	activity:=activity.Activity{uint64(1),"new note"}
	jsonEncoding:=mcstorage.JsonEncoding{reflect.TypeOf(activity)}
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
	storages:=make([]mcstorage.Storage,3)
	storages[0]=redisStorage
	storages[1]=redisListStorage
	storages[2]=redisCounterStorage

	feed:=BaseFeed{redisStorage,redisListStorage,redisCounterStorage,"fanngyuan"}
	activity1:=activity.Activity{uint64(1),"new note"}
	feed.AddActivity(activity1)

	result:=feed.GetActivities(uint64(0),uint64(0),1,10)
	if len(result)!=1{
		t.Error("len should be 1")
	}
	if result[0].GetId()!=1{
		t.Error("id should be 1")
	}
	count:=feed.GetCount()
	if count!=1{
		t.Error("count should be 1")
	}

	activity2:=activity.Activity{uint64(2),"new note"}
	feed.AddActivity(activity2)

	result=feed.GetActivities(uint64(0),uint64(0),1,10)
	if len(result)!=2{
		t.Error("len should be 2")
	}
	if result[0].GetId()!=2{
		t.Error("id should be 2")
	}
	if result[1].GetId()!=1{
		t.Error("id should be 1")
	}
	count=feed.GetCount()
	if count!=2{
		t.Error("count should be 2")
	}

	activity3:=activity.Activity{uint64(3),"new note"}
	feed.AddActivity(activity3)

	result=feed.GetActivities(uint64(0),uint64(0),1,10)

	if len(result)!=3{
		t.Error("len should be 3")
	}
	if result[0].GetId()!=3{
		t.Error("id should be 3")
	}
	if result[1].GetId()!=2{
		t.Error("id should be 2")
	}
	if result[2].GetId()!=1{
		t.Error("id should be 1")
	}
	count=feed.GetCount()
	if count!=3{
		t.Error("count should be 3")
	}

	feed.RemoveActivity(activity2)
	result=feed.GetActivities(uint64(0),uint64(0),1,10)

	if len(result)!=2{
		t.Error("len should be 2")
	}
	if result[0].GetId()!=3{
		t.Error("id should be 3")
	}
	if result[1].GetId()!=1{
		t.Error("id should be 1")
	}
	count=feed.GetCount()
	if count!=2{
		t.Error("count should be 2")
	}

	flush(storages)
}
