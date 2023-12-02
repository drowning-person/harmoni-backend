package v1

import "harmoni/app/harmoni/internal/entity/user"

func (u *UserBasic) FromDomain(user *user.UserBasicInfo) *UserBasic {
	u.Id = user.UserID
	u.Name = user.Name
	u.Avatar = user.Avatar
	return u
}

func ListFromDomain(users []*user.UserBasicInfo) []*UserBasic {
	var list []*UserBasic
	for i := range users {
		list = append(list, (&UserBasic{}).FromDomain(users[i]))
	}
	return list
}

type UserBasicList []*UserBasic

func (l UserBasicList) ToMap() map[int64]*UserBasic {
	m := map[int64]*UserBasic{}
	for i := range l {
		m[l[i].GetId()] = l[i]
	}
	return m
}
