package feedly

import (
	"github.com/fanngyuan/feedly/feed"
	"github.com/fanngyuan/feedly/activity"
	"github.com/fanngyuan/feedly/id"
	mcstorage "github.com/fanngyuan/mcstorage"
	"strconv"
)

type UserFeedly interface{
	Follow(userId,targetId uint64)
	UnFollow(userId,targetId uint64)
	GetFollowing(userId uint64,sinceId,maxId uint64,page,count int)[]User
	GetFollower(userId uint64,sinceId,maxId uint64,page,count int)[]User
	GetUserTimeline(userId uint64,activityType string,sinceId,maxId uint64,page,count int)[]activity.Activable
	AddActivity(userId uint64,activity activity.Activable)
	RemoveActivity(userId uint64,activity activity.Activable)
	ActivityCount(userId uint64,activityType string)int
}

type ActivityInit interface{
	InitUserFeed(userId uint64) feed.UserFeed
}

var Follow string = "follow"
var UnFollow string = "unfollow"

type AggregatorUserFeedly interface{
	UserFeedly
	GetFriendsTimeline(userId uint64,activityType string,sinceId,maxId uint64,page,count int)[]activity.Activable
	GetHomeTimeline(userId uint64,activityType string,sinceId,maxId uint64,page,count int)[]activity.Activable
}

type AggregatorFeed struct{
	UserFeedMap map[uint64] feed.UserFeed
	IdGenerator id.IdGenerable
	UserStorage mcstorage.Storage
	ActivityInit ActivityInit
}

type PullAgrregatorFeed struct{
	AggregatorFeed
}

type PushAggregatorFeed struct{
	AggregatorFeed
	AsyncBoundry int
}

func (this AggregatorFeed) Follow(userId,targetId uint64){
	userFeed,ok:=this.UserFeedMap[userId]
	if !ok{
		userFeed=this.ActivityInit.InitUserFeed(userId)
		this.UserFeedMap[userId]=userFeed
	}
	followId:=this.IdGenerator.GetId()
	followActivity := activity.FollowActivity{followId,userId,targetId,Follow}
	userFeed.Follow(followActivity)

	targetUserFeed,ok:=this.UserFeedMap[targetId]
	if !ok{
		targetUserFeed=this.ActivityInit.InitUserFeed(targetId)
		this.UserFeedMap[targetId]=targetUserFeed
	}
	targetUserFeed.AddFollower(followActivity)
}

func (this AggregatorFeed) UnFollow(userId,targetId uint64){
	userFeed,ok:=this.UserFeedMap[userId]
	if !ok{
		userFeed=this.ActivityInit.InitUserFeed(userId)
		this.UserFeedMap[userId]=userFeed
	}
	unfollowActivity := activity.FollowActivity{0,userId,targetId,UnFollow}
	userFeed.Unfollow(unfollowActivity)

	targetUserFeed,ok:=this.UserFeedMap[targetId]
	if !ok{
		targetUserFeed=this.ActivityInit.InitUserFeed(targetId)
		this.UserFeedMap[targetId]=targetUserFeed
	}
	targetUserFeed.RemoveFollower(unfollowActivity)
}

func (this AggregatorFeed) GetFollowing(userId uint64,sinceId,maxId uint64,page,count int)[]User{
	uids:=this.GetFollowingIds(userId,sinceId,maxId,page,count)
	if uids==nil{
		return nil
	}
	return this.RenderUsers(uids)
}

func (this AggregatorFeed) GetFollowingIds(userId uint64,sinceId,maxId uint64,page,count int)[]interface{}{
	userFeed,ok:=this.UserFeedMap[userId]
	if !ok{
		userFeed=this.ActivityInit.InitUserFeed(userId)
	}
	following:=userFeed.GetFollowing(sinceId,maxId,page,count)
	uids:=make([]interface{},len(following))
	for index,uid := range following{
		uids[index]=strconv.Itoa(int(uid.TargetId))
	}
	return uids
}

func (this AggregatorFeed) RenderUsers(uids []interface{})[]User{
	userMap,err:=this.UserStorage.MultiGet(uids)
	if err!=nil{
		return nil
	}
	var users []User
	for _,uid := range uids{
		user,ok:=userMap[uid]
		if ok{
			users=append(users,user.(User))
		}
	}
	return users
}

func (this AggregatorFeed) GetFollower(userId uint64,sinceId,maxId uint64,page,count int)[]User{
	userFeed,ok:=this.UserFeedMap[userId]
	if !ok{
		userFeed=this.ActivityInit.InitUserFeed(userId)
	}
	followers:=userFeed.GetFollower(sinceId,maxId,page,count)
	uids:=make([]interface{},len(followers))
	for index,uid := range followers{
		uids[index]=strconv.Itoa(int(uid.UserId))
	}
	return this.RenderUsers(uids)
}

func (this AggregatorFeed) GetUserTimeline(userId uint64,activityType string,sinceId,maxId uint64,page,count int)[]activity.Activable{
	userFeed,ok:=this.UserFeedMap[userId]
	if !ok{
		userFeed=this.ActivityInit.InitUserFeed(userId)
	}
	return userFeed.GetActivities(activityType,sinceId,maxId,page,count)
}
/**
func (this PullAggregatorFeed) GetFriendsTimeline(userId uint64,activityType string,sinceId,maxId uint64,page,count int)[]activity.Activable{
	uids:=this.GetFollowingIds(userId)
	if uids==nil{
		return nil
	}
	followingIds:=make([]int64,len(uids))
	for i,uid :=range uids{
		followingId,err:=strconv.Atoi(uid)
		followingIds[i]=followingId
	}
}

func (this PullAggregatorFeed) GetHomeTimeline(userId uint64,activityType string,sinceId,maxId uint64,page,count int)[]activity.Activable{

}
*/
func (this AggregatorFeed) AddActivity(userId uint64,activity activity.Activable){
	userFeed,ok:=this.UserFeedMap[userId]
	if !ok{
		userFeed=this.ActivityInit.InitUserFeed(userId)
	}
	userFeed.AddActivity(activity)
}

func (this AggregatorFeed) RemoveActivity(userId uint64,activity activity.Activable){
	userFeed,ok:=this.UserFeedMap[userId]
	if !ok{
		userFeed=this.ActivityInit.InitUserFeed(userId)
	}
	userFeed.RemoveActivity(activity)
}

func (this AggregatorFeed) ActivityCount(userId uint64,feedType string)int{
	userFeed,ok:=this.UserFeedMap[userId]
	if !ok{
		userFeed=this.ActivityInit.InitUserFeed(userId)
	}
	return userFeed.GetActivitiesCount(feedType)
}
