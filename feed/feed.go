package feed

import (
	. "github.com/fanngyuan/mcstorage"
	"github.com/fanngyuan/feedly/activity"
	"strconv"
)

type Feedable interface{
	AddActivity(activity activity.Activable)
	AddActivities(activities []activity.Activable)
	RemoveActivity(activity activity.Activable)
	RemoveActivities(activities []activity.Activable)
	GetActivities(sinceId ,maxId uint64,page ,count int)[]activity.Activable
	GetActivityIds(sinceId ,maxId uint64,page ,count int)[]uint64
	MultiGet(ids []uint64)[]activity.Activable
	GetCount()int
	GetId()string
}

type BaseFeed struct{
	ActivityStorage Storage
	TimelimeStorage ListStorage
	CounterStorage  CounterStorage
	Id string
}

func (this BaseFeed) GetId()string{
	return this.Id
}

func (this BaseFeed) AddActivity(activity activity.Activable){
	activityKey:=strconv.Itoa(int(activity.GetId()))
	this.ActivityStorage.Set(activityKey,activity)
	this.TimelimeStorage.AddItem(this.GetId(),activity.GetId())
	this.CounterStorage.Incr(this.GetId(),1)
}

func (this BaseFeed) AddActivities(activities []activity.Activable){
	activityMap:=make(map[interface{}] interface{})
	ids:=make([]uint64,len(activities))
	for i,activity := range activities{
		key:=strconv.Itoa(int(activity.GetId()))
		activityMap[key]=activity
		ids[i]=activity.GetId()
	}
	this.ActivityStorage.MultiSet(activityMap)
	this.TimelimeStorage.Set(this.GetId(),ids)
	this.CounterStorage.Incr(this.GetId(),uint64(len(ids)))
}

func (this BaseFeed) RemoveActivity(activity activity.Activable){
	key:=strconv.Itoa(int(activity.GetId()))
	this.ActivityStorage.Delete(key)
	this.TimelimeStorage.DeleteItem(this.GetId(),activity.GetId())
	this.CounterStorage.Decr(this.GetId(),1)
}

func (this BaseFeed) RemoveActivities(activities []activity.Activable){
	for _,activity := range activities{
		this.RemoveActivity(activity)
	}
}

func (this BaseFeed) GetActivities(sinceId ,maxId uint64,page ,count int)[]activity.Activable{
	ids,err:=this.TimelimeStorage.Getlimit(this.GetId(),sinceId,maxId,page,count)
	if err!=nil{
		return nil
	}
	var keys []interface{}
	for _,id :=range(ids.(IntReversedSlice)){
		keys=append(keys,strconv.Itoa(id))
	}
	values,err:=this.ActivityStorage.MultiGet(keys)
	result:=make([]activity.Activable,len(values))
	i:=0
	for _,value := range(values){
		result[i]=value.(activity.Activable)
		i=i+1
	}
	return result
}

func (this BaseFeed) GetActivityIds(sinceId ,maxId uint64,page ,count int)[]uint64{
	ids,err:=this.TimelimeStorage.Getlimit(this.GetId(),sinceId,maxId,page,count)
	if err!=nil{
		return nil
	}
	result:=make([]uint64,ids.(IntReversedSlice).Len())
	for i,id :=range(ids.(IntReversedSlice)){
		result[i]=uint64(id)
	}
	return result
}

func (this BaseFeed) MultiGet(ids []uint64)[]activity.Activable{
	var keys []interface{}
	for _,id :=range ids {
		keys=append(keys,strconv.Itoa(int(id)))
	}
	values,err:=this.ActivityStorage.MultiGet(keys)
	if err!=nil{
		return make([]activity.Activable,0)
	}
	result:=make([]activity.Activable,len(values))
	i:=0
	for _,value := range(values){
		result[i]=value.(activity.Activable)
		i=i+1
	}
	return result
}

func (this BaseFeed) GetCount()int{
	count,err:=this.CounterStorage.Get(this.GetId())
	if err!=nil{
		return 0
	}
	return count.(int)
}
