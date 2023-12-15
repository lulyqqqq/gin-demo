package dto

type UserPageList struct {
	UserInfoDtoList []*UserInfoDto
	Total           int64
}
