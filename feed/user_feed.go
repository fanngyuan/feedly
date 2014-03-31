package feed

import (
	"github.com/fanngyuan/feedly/activity"
)

type UserFeedable interface{
	Follow(followActivity activity.FollowActivity)
	Unfollow(unfollowActivity activity.FollowActivity)
	AddFollower(followActivity activity.FollowActivity)
	RemoveFollower(unfollowActivity activity.FollowActivity)
	AddActivity(activity activity.Activable)
	RemoveActivity(activity activity.Activable)
	GetActivities(feedType string,sinceId,maxId uint64,page,count int)[]activity.Activable
	GetActivitiesCount(feedType string)int
	GetFollowing(sinceId,maxId uint64,page,count int)[]activity.FollowActivity
	GetFollower(sinceId,maxId uint64,page,count int)[]activity.FollowActivity
}

type UserFeed struct{
	UserId uint64
	FollowingFeed Feedable
	FollowerFeed Feedable
	ActivityFeedMap map[string] Feedable
}

func (this UserFeed) Follow(followActivity activity.FollowActivity){
	this.FollowingFeed.AddActivity(followActivity)
}

func (this UserFeed) Unfollow(unfollowActivity activity.FollowActivity){
	this.FollowingFeed.RemoveActivity(unfollowActivity)
}

func (this UserFeed) AddFollower(followActivity activity.FollowActivity){
	this.FollowerFeed.AddActivity(followActivity)
}

func (this UserFeed) RemoveFollower(unfollowActivity activity.FollowActivity){
	this.FollowerFeed.RemoveActivity(unfollowActivity)
}

func (this UserFeed) AddActivity(activity activity.Activable){
	activityFeed,ok:=this.ActivityFeedMap[activity.GetType()]
	if ok{
		activityFeed.AddActivity(activity)
	}
}

func (this UserFeed) GetActivities(feedType string,sinceId,maxId uint64,page,count int)[]activity.Activable{
	activityFeed,ok:=this.ActivityFeedMap[feedType]
	if ok{
		return activityFeed.GetActivities(sinceId,maxId,page,count)
	}
	return nil
}

func (this UserFeed) RemoveActivity(activity activity.Activable){
	activityFeed,ok:=this.ActivityFeedMap[activity.GetType()]
	if ok{
		activityFeed.RemoveActivity(activity)
	}
}

func (this UserFeed) GetActivitiesCount(feedType string)int{
	activityFeed,ok:=this.ActivityFeedMap[feedType]
	if ok{
		activityFeed.GetCount()
	}
	return 0
}

func (this UserFeed) GetFollowing(sinceId,maxId uint64,page,count int)[]activity.FollowActivity{
	activities:=this.FollowingFeed.GetActivities(sinceId,maxId,page,count)
	if activities==nil{
		return nil
	}
	followActivities:=make([]activity.FollowActivity,len(activities))
	for i,activityItem:=range activities{
		followActivities[i]=activityItem.(activity.FollowActivity)
	}
	return followActivities
}

func (this UserFeed) GetFollower(sinceId,maxId uint64,page,count int)[]activity.FollowActivity{
	activities:=this.FollowerFeed.GetActivities(sinceId,maxId,page,count)
	if activities==nil{
		return nil
	}
	followActivities:=make([]activity.FollowActivity,len(activities))
	for i,activityItem:=range activities{
		followActivities[i]=activityItem.(activity.FollowActivity)
	}
	return followActivities
}
