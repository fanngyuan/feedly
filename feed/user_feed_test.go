package feed

import (
	"testing"
	"github.com/fanngyuan/feedly/activity"
	mcstorage "github.com/fanngyuan/mcstorage"
	"reflect"
)

func newRedisFollowStorage()mcstorage.RedisStorage{
	followActivity:=activity.FollowActivity{uint64(1),uint64(1),uint64(2),"follow"}
	jsonEncoding:=mcstorage.JsonEncoding{reflect.TypeOf(followActivity)}
	redisStorage,_ := mcstorage.NewRedisStorage(":6379", "test", 0, jsonEncoding)
	return redisStorage
}

func newFollowIdMapStorage()mcstorage.RedisStorage{
	jsonEncoding:=mcstorage.JsonEncoding{reflect.TypeOf(1)}
	redisStorage,_ := mcstorage.NewRedisStorage(":6379", "test_follow_map", 0, jsonEncoding)
	return redisStorage
}

func newFollowerIdMapStorage()mcstorage.RedisStorage{
	jsonEncoding:=mcstorage.JsonEncoding{reflect.TypeOf(1)}
	redisStorage,_ := mcstorage.NewRedisStorage(":6379", "test_follower_map", 0, jsonEncoding)
	return redisStorage
}

func TestFollow(t *testing.T) {
	redisStorage:=newRedisStorage()
	redisListStorage:=newRedisListStorage()
	redisCounterStorage:=newRedisCounterStorage()
	redisFollowStorage:=newRedisFollowStorage()

	redisFollowMapStorage:=newFollowIdMapStorage()
	redisFollowerMapStorage:=newFollowerIdMapStorage()
	storages:=make([]mcstorage.Storage,4)
	storages[0]=redisStorage
	storages[1]=redisListStorage
	storages[2]=redisCounterStorage
	storages[3]=redisFollowStorage

	followingFeed:=BaseFeed{redisFollowStorage,redisListStorage,redisCounterStorage,"fanngyuan_follow"}
	followerFeed:=BaseFeed{redisFollowStorage,redisListStorage,redisCounterStorage,"fanngyuan_follower"}
	feed:=BaseFeed{redisStorage,redisListStorage,redisCounterStorage,"status"}
	activityFeedMap:=make(map[string] Feedable)
	activityFeedMap["status"]=feed
	userFeed:=UserFeed{uint64(1),followingFeed,followerFeed,activityFeedMap,redisFollowMapStorage,redisFollowerMapStorage}

	followActivity1:=activity.FollowActivity{uint64(1),uint64(1),uint64(2),"follow"}
	userFeed.Follow(followActivity1)
	following:=userFeed.GetFollowing(0,0,1,10)

	if len(following)!=1{
		t.Error("len should be 1")
	}
	if following[0].GetId()!=1{
		t.Error("id should be 1")
	}

	followActivity2:=activity.FollowActivity{uint64(2),uint64(1),uint64(3),"follow"}
	userFeed.Follow(followActivity2)
	following=userFeed.GetFollowing(0,0,1,10)

	if len(following)!=2{
		t.Error("len should be 2")
	}
	if following[0].GetId()!=2{
		t.Error("id should be 2")
	}
	if following[1].GetId()!=1{
		t.Error("id should be 1")
	}

	followActivity3:=activity.FollowActivity{uint64(3),uint64(1),uint64(4),"follow"}
	userFeed.Follow(followActivity3)
	following=userFeed.GetFollowing(0,0,1,10)

	if len(following)!=3{
		t.Error("len should be 3")
	}
	if following[0].GetId()!=3{
		t.Error("id should be 3")
	}
	if following[1].GetId()!=2{
		t.Error("id should be 2")
	}
	if following[2].GetId()!=1{
		t.Error("id should be 1")
	}

	userFeed.Unfollow(followActivity2)
	following=userFeed.GetFollowing(0,0,1,10)

	if len(following)!=2{
		t.Error("len should be 2")
	}
	if following[0].GetId()!=3{
		t.Error("id should be 2")
	}
	if following[1].GetId()!=1{
		t.Error("id should be 1")
	}

	userFeed.Unfollow(activity.FollowActivity{uint64(0),uint64(1),uint64(2),"follow"})
	following=userFeed.GetFollowing(0,0,1,10)

	if len(following)!=1{
		t.Error("len should be 1")
	}
	if following[0].GetId()!=3{
		t.Error("id should be 2")
	}

	flush(storages)
}

func TestFollower(t *testing.T) {
	redisStorage:=newRedisStorage()
	redisListStorage:=newRedisListStorage()
	redisCounterStorage:=newRedisCounterStorage()
	redisFollowStorage:=newRedisFollowStorage()
	storages:=make([]mcstorage.Storage,4)
	storages[0]=redisStorage
	storages[1]=redisListStorage
	storages[2]=redisCounterStorage
	storages[3]=redisFollowStorage

	redisFollowMapStorage:=newFollowIdMapStorage()
	redisFollowerMapStorage:=newFollowerIdMapStorage()

	followingFeed:=BaseFeed{redisFollowStorage,redisListStorage,redisCounterStorage,"fanngyuan_follow"}
	followerFeed:=BaseFeed{redisFollowStorage,redisListStorage,redisCounterStorage,"fanngyuan_follower"}
	feed:=BaseFeed{redisStorage,redisListStorage,redisCounterStorage,"status"}
	activityFeedMap:=make(map[string] Feedable)
	activityFeedMap["status"]=feed
	userFeed:=UserFeed{uint64(1),followingFeed,followerFeed,activityFeedMap,redisFollowMapStorage,redisFollowerMapStorage}

	followActivity1:=activity.FollowActivity{uint64(1),uint64(2),uint64(1),"follow"}
	userFeed.AddFollower(followActivity1)
	follower:=userFeed.GetFollower(0,0,1,10)

	if len(follower)!=1{
		t.Error("len should be 1")
	}
	if follower[0].GetId()!=1{
		t.Error("id should be 1")
	}

	followActivity2:=activity.FollowActivity{uint64(2),uint64(3),uint64(1),"follow"}
	userFeed.AddFollower(followActivity2)
	follower=userFeed.GetFollower(0,0,1,10)

	if len(follower)!=2{
		t.Error("len should be 2")
	}
	if follower[0].GetId()!=2{
		t.Error("id should be 2")
	}
	if follower[1].GetId()!=1{
		t.Error("id should be 1")
	}

	followActivity3:=activity.FollowActivity{uint64(3),uint64(4),uint64(1),"follow"}
	userFeed.AddFollower(followActivity3)
	follower=userFeed.GetFollower(0,0,1,10)

	if len(follower)!=3{
		t.Error("len should be 3")
	}
	if follower[0].GetId()!=3{
		t.Error("id should be 3")
	}
	if follower[1].GetId()!=2{
		t.Error("id should be 2")
	}
	if follower[2].GetId()!=1{
		t.Error("id should be 1")
	}

	userFeed.RemoveFollower(followActivity2)
	follower=userFeed.GetFollower(0,0,1,10)

	if len(follower)!=2{
		t.Error("len should be 2")
	}
	if follower[0].GetId()!=3{
		t.Error("id should be 3")
	}
	if follower[1].GetId()!=1{
		t.Error("id should be 1")
	}

	userFeed.RemoveFollower(activity.FollowActivity{uint64(0),uint64(4),uint64(1),"follow"})
	follower=userFeed.GetFollower(0,0,1,10)

	if len(follower)!=1{
		t.Error("len should be 1")
	}
	if follower[0].GetId()!=1{
		t.Error("id should be 1")
	}

	flush(storages)
}

func TestUserActivity(t *testing.T) {
	redisStorage:=newRedisStorage()
	redisListStorage:=newRedisListStorage()
	redisCounterStorage:=newRedisCounterStorage()
	redisFollowStorage:=newRedisFollowStorage()
	storages:=make([]mcstorage.Storage,4)
	storages[0]=redisStorage
	storages[1]=redisListStorage
	storages[2]=redisCounterStorage
	storages[3]=redisFollowStorage

	redisFollowMapStorage:=newFollowIdMapStorage()
	redisFollowerMapStorage:=newFollowerIdMapStorage()

	followingFeed:=BaseFeed{redisFollowStorage,redisListStorage,redisCounterStorage,"fanngyuan_follow"}
	followerFeed:=BaseFeed{redisFollowStorage,redisListStorage,redisCounterStorage,"fanngyuan_follower"}
	feed:=BaseFeed{redisStorage,redisListStorage,redisCounterStorage,"status"}
	activityFeedMap:=make(map[string] Feedable)
	activityFeedMap["status"]=feed
	userFeed:=UserFeed{uint64(1),followingFeed,followerFeed,activityFeedMap,redisFollowMapStorage,redisFollowerMapStorage}

	activity1:=activity.Activity{uint64(1),"status"}
	userFeed.AddActivity(activity1)

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

	flush(storages)
}
