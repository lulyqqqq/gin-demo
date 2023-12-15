package dto

import "gin-use-demo/model"

type UserInfoDto struct {
	Id      int    `gorm:"type:int" json:"id"`
	Name    string `gorm:"varchar(11);not null;unique" json:"name"`
	Number  string `gorm:"varchar(11);not null;unique" json:"number"`
	Address string `gorm:"varchar(256)" json:"address"`
	Tag     string `gorm:"varchar(5)" json:"tag"`
	Role    string `gorm:"varchar(2)" json:"role"` // 0-管理员 1-正常用户 2-禁止使用的用户
}

func ToUserInfoDto(user model.User) UserInfoDto {
	return UserInfoDto{
		Id:      user.Id,
		Name:    user.Name,
		Number:  user.Number,
		Address: user.Address,
		Tag:     user.Tag,
		Role:    user.Role,
	}
}
