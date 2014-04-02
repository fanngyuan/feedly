package feedly

import (
	"testing"
	"github.com/fanngyuan/feedly/id"
	"github.com/fanngyuan/feedly/feed"
	"github.com/fanngyuan/feedly/activity"
	mcstorage "github.com/fanngyuan/mcstorage"
	"strconv"
	"reflect"
)

func newUserRedisStorage()mcstorage.RedisStorage{
	user:=User{uint64(1)}
	jsonEncoding:=mcstorage.JsonEncoding{reflect.TypeOf(user)}
	redisStorage,_ := mcstorage.NewRedisStorage(":6379", "test_user", 0, jsonEncoding)
	return redisStorage
}

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

func newRedisFollowStorage()mcstorage.RedisStorage{
	followActivity:=activity.FollowActivity{uint64(1),uint64(1),uint64(2),"follow"}
	jsonEncoding:=mcstorage.JsonEncoding{reflect.TypeOf(followActivity)}
	redisStorage,_ := mcstorage.NewRedisStorage(":6379", "user_follow", 0, jsonEncoding)
	return redisStorage
}

func newFollowIdMapStorage(prefix string)mcstorage.RedisStorage{
	jsonEncoding:=mcstorage.JsonEncoding{reflect.TypeOf(1)}
	redisStorage,_ := mcstorage.NewRedisStorage(":6379", "test_follow_map_"+prefix, 0, jsonEncoding)
	return redisStorage
}

func newFollowerIdMapStorage(prefix string)mcstorage.RedisStorage{
	jsonEncoding:=mcstorage.JsonEncoding{reflect.TypeOf(1)}
	redisStorage,_ := mcstorage.NewRedisStorage(":6379", "test_follower_map_"+prefix, 0, jsonEncoding)
	return redisStorage
}

func addUsers(storage mcstorage.Storage,userIds []uint64)[]User{
	users:=make([]User,len(userIds))
	for i,userId := range userIds{
		user:=User{userId}
		storage.Set(strconv.Itoa(int(userId)),user)
		users[i]=user
	}
	return users
}

func flush(storage mcstorage.Storage){
	storage.FlushAll()
}

type ActivityFeedInit struct{

}

func (this ActivityFeedInit) InitUserFeed(userId uint64)feed.UserFeed{
	redisStorage:=newRedisStorage()
	redisListStorage:=newRedisListStorage()
	redisCounterStorage:=newRedisCounterStorage()
	redisFollowStorage:=newRedisFollowStorage()

	followingFeed:=feed.BaseFeed{redisFollowStorage,redisListStorage,redisCounterStorage,strconv.Itoa(int(userId))+"_following"}
	followerFeed:=feed.BaseFeed{redisFollowStorage,redisListStorage,redisCounterStorage,strconv.Itoa(int(userId))+"_follower"}
	statusFeed:=feed.BaseFeed{redisStorage,redisListStorage,redisCounterStorage,"status"}
	activityFeedMap:=make(map[string] feed.Feedable)
	activityFeedMap["status"]=statusFeed
	followMapStorage:=newFollowIdMapStorage(strconv.Itoa(int(userId)))
	followerMapStorage:=newFollowerIdMapStorage(strconv.Itoa(int(userId)))
	userFeed:=feed.UserFeed{userId,followingFeed,followerFeed,activityFeedMap,followMapStorage,followerMapStorage}
	return userFeed
}

func TestFollowUnfollow(t *testing.T) {
	userRedisStorage:=newUserRedisStorage()
	zero:=int64(0)
	idgen:=id.IdGenerator{&zero}
	userFeedMap:=make(map[uint64] feed.UserFeed)

	aggregateFeed:=AggregatorFeed{userFeedMap,idgen,userRedisStorage,ActivityFeedInit{}}

	userIds:=make([]uint64,20)
	for i:=1;i<=20;i++{
		userIds[i-1]=uint64(i)
	}
	addUsers(userRedisStorage,userIds)

	aggregateFeed.Follow(uint64(1),uint64(2))
	aggregateFeed.Follow(uint64(1),uint64(3))
	aggregateFeed.Follow(uint64(1),uint64(4))
	aggregateFeed.Follow(uint64(1),uint64(5))

	users:=aggregateFeed.GetFollowing(uint64(1),uint64(0),uint64(0),1,20)
	if len(users)!=4{
		t.Errorf("length of users should be 4 result is ",len(users))
	}
	if users[0].UserId!=5{
		t.Errorf("last one should be 5 ,result is ",users[0].UserId)
	}
	if users[3].UserId!=2{
		t.Errorf("last one should be 2 ,result is ",users[3].UserId)
	}

	aggregateFeed.UnFollow(uint64(1),uint64(4))
	users=aggregateFeed.GetFollowing(uint64(1),uint64(0),uint64(0),1,20)
	if len(users)!=3{
		t.Errorf("length of users should be 3 result is ",len(users))
	}
	if users[0].UserId!=5{
		t.Errorf("last one should be 5 ,result is ",users[0].UserId)
	}
	if users[2].UserId!=2{
		t.Errorf("last one should be 2 ,result is ",users[3].UserId)
	}

	aggregateFeed.Follow(uint64(2),uint64(5))
	aggregateFeed.Follow(uint64(3),uint64(5))

	users=aggregateFeed.GetFollower(uint64(5),uint64(0),uint64(0),1,20)
	if len(users)!=3{
		t.Errorf("length of users should be 3 result is ",len(users))
	}
	if users[0].UserId!=3{
		t.Errorf("last one should be 3 ,result is ",users[0].UserId)
	}
	if users[2].UserId!=1{
		t.Errorf("last one should be 1 ,result is ",users[3].UserId)
	}

	aggregateFeed.UnFollow(uint64(2),uint64(5))
	users=aggregateFeed.GetFollower(uint64(5),uint64(0),uint64(0),1,20)
	if len(users)!=2{
		t.Errorf("length of users should be 2 result is ",len(users))
	}
	if users[0].UserId!=3{
		t.Errorf("last one should be 3 ,result is ",users[0].UserId)
	}
	if users[1].UserId!=1{
		t.Errorf("last one should be 1 ,result is ",users[3].UserId)
	}

	flush(userRedisStorage)
}

func TestAddRemoveActiviy(t *testing.T) {
	userRedisStorage:=newUserRedisStorage()
	zero:=int64(0)
	idgen:=id.IdGenerator{&zero}
	userFeedMap:=make(map[uint64] feed.UserFeed)

	aggregateFeed:=AggregatorFeed{userFeedMap,idgen,userRedisStorage,ActivityFeedInit{}}

	userIds:=make([]uint64,20)
	for i:=1;i<=20;i++{
		userIds[i-1]=uint64(i)
	}
	addUsers(userRedisStorage,userIds)

	activity1:=activity.Activity{uint64(1),"status"}
	activity2:=activity.Activity{uint64(2),"status"}
	activity3:=activity.Activity{uint64(3),"status"}
	activity4:=activity.Activity{uint64(4),"status"}

	aggregateFeed.AddActivity(uint64(1),activity1)
	aggregateFeed.AddActivity(uint64(1),activity2)
	aggregateFeed.AddActivity(uint64(1),activity3)
	aggregateFeed.AddActivity(uint64(1),activity4)

	activities:=aggregateFeed.GetUserTimeline(uint64(1),"status",uint64(0),uint64(0),1,20)
	if len(activities)!=4{
		t.Errorf("length of users should be 4 result is ",len(activities))
	}
	if activities[0].GetId()!=4{
		t.Errorf("last one should be 4 ,result is ",activities)
	}
	if activities[3].GetId()!=1{
		t.Errorf("last one should be 1 ,result is ",activities)
	}

	aggregateFeed.RemoveActivity(uint64(1),activity3)

	activities=aggregateFeed.GetUserTimeline(uint64(1),"status",uint64(0),uint64(0),1,20)
	if len(activities)!=3{
		t.Errorf("length of users should be 3 result is ",len(activities))
	}
	if activities[0].GetId()!=4{
		t.Errorf("last one should be 4 ,result is ",activities)
	}
	if activities[2].GetId()!=1{
		t.Errorf("last one should be 1 ,result is ",activities)
	}

	flush(userRedisStorage)
}
