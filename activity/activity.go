package activity

type Activable interface{
	GetId()uint64
	GetType()string
}

type Activity struct{
	Id uint64
	Type string
}

func (this Activity) GetId()uint64{
	return this.Id
}

func (this Activity) GetType()string{
	return this.Type
}

type FollowActivity struct{
	FollowId uint64
	UserId uint64
	TargetId uint64
	Type string
}

func (this FollowActivity) GetId()uint64{
	return this.FollowId
}

func (this FollowActivity) GetType()string{
	return this.Type
}
