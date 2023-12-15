package service

import (
	"errors"
	"gin-use-demo/common"
	"gin-use-demo/dto"
	"gin-use-demo/model"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"log"
)

type IUserService interface {
	Login(name, number string) (*model.User, error)
	AddUser(user *model.User) (bool, error)
	GetUserInfo(id string) (userInfo *dto.UserInfoDto, err error)
	GetUserInfos(id string) (userInfo *model.User, err error)
	GetUserPageList(pageNum, pageSize int, name, number string) (u *dto.UserPageList, err error)
	DeleteUser(id string) error
	UpdateUser(id string, user *model.User) error
	IsExistUser(user *model.User) (bool, error)
}

type UserService struct {
	DB *gorm.DB
}

func (u UserService) Login(name, number string) (*model.User, error) {
	var user model.User
	u.DB.Where("number = ? AND name = ?", number, name).First(&user)
	if user.Id != 0 {
		return &user, nil
	} else {
		return nil, errors.New("用户不存在")
	}

}

func (u UserService) AddUser(user *model.User) (bool, error) {
	// 是否存在的user
	var isExistUser model.User
	// 判断
	u.DB.Where("number = ? OR name= ?", user.Number, user.Name).First(&isExistUser)
	if isExistUser.Id != 0 {
		return false, errors.New("手机号或者用户名已存在")
	}
	// 不存在用户,新增
	err := u.DB.Debug().Create(user).Error
	if err != nil {
		return false, err
	}
	return true, nil

}

func (u UserService) IsExistUser(user *model.User) (bool, error) {
	// 是否存在的user
	var isExistUser model.User
	// 判断
	u.DB.Where("number = ? OR name= ?", user.Number, user.Name).First(&isExistUser)
	if isExistUser.Id != 0 {
		return false, errors.New("手机号或者用户名已存在")
	}
	return true, nil
}

func (u UserService) GetUserInfo(id string) (userInfo *dto.UserInfoDto, err error) {
	var user *model.User
	userInfo = &dto.UserInfoDto{}
	err = u.DB.Debug().Where("id = ?", id).First(&user).Error
	if err != nil {
		return &dto.UserInfoDto{}, err
	}
	err = copier.Copy(&userInfo, &user)

	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

func (u UserService) GetUserInfos(id string) (userInfo *model.User, err error) {
	var user *model.User
	err = u.DB.Debug().Where("id = ?", id).First(&user).Error
	if err != nil {
		return &model.User{}, err
	}
	return user, nil
}

func (u UserService) GetUserPageList(pageNum, pageSize int, name, number string) (userPageList *dto.UserPageList, err error) {
	var userList *[]model.User
	var total int64

	query := u.DB.Offset((pageNum - 1) * pageSize).Limit(pageSize)

	// 构建查询条件
	if name != "" && number != "" {
		query = query.Where("name LIKE ? AND number LIKE ?", "%"+name+"%", "%"+number+"%")
	} else if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	} else if number != "" {
		query = query.Where("number LIKE ?", "%"+number+"%")
	}

	if err = query.Find(&userList).Error; err != nil {
		return &dto.UserPageList{}, err
	}

	userInfoList := make([]*dto.UserInfoDto, len(*userList))

	// 复制

	err = copier.Copy(&userInfoList, &userList)
	if err != nil {
		log.Default().Println("copier失效")
		return &dto.UserPageList{}, err
	}
	// 查询查询结果的总数
	query.Find(&userList).Count(&total)

	return &dto.UserPageList{
		UserInfoDtoList: userInfoList,
		Total:           total}, nil
}

func (u UserService) DeleteUser(id string) error {
	err := u.DB.Where("id =?", id).Delete(&model.User{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (u UserService) UpdateUser(id string, user *model.User) error {
	err := u.DB.Model(&user).Where("id = ?", id).Updates(model.User{
		Name:     user.Name,
		Password: user.Password,
		Number:   user.Number,
		Address:  user.Address,
		Tag:      user.Tag,
		Role:     user.Role,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func NewUserService() IUserService {
	return UserService{DB: common.GetDB()}
}
