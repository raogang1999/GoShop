package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"shop_srvs/user_srv/global"
	"shop_srvs/user_srv/model"
	"shop_srvs/user_srv/proto"
	"strings"
	"time"
)

type UserServer struct {
	proto.UnimplementedUserServer
}

//type UserServer interface {
//	GetUserList(context.Context, *PageInfo) (*UserListResponse, error)
//	GetUserByMobile(context.Context, *UserMobileRequest) (*UserInfoResponse, error)
//	GetUserById(context.Context, *IdRequest) (*UserInfoResponse, error)
//	CreateUser(context.Context, *CreateUserInfo) (*UserInfoResponse, error)
//	UpdateUser(context.Context, *UpdateUserInfo) (*empty.Empty, error)
//	CheckPassword(context.Context, *CheckPasswordInfo) (*CheckResponse, error)
//	mustEmbedUnimplementedUserServer()
//}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func ModelToResponse(user model.User) proto.UserInfoResponse {
	//grpc中字段有默认值，不能赋值为nil
	//需要知道哪些字段是空的
	userInfoRsp := proto.UserInfoResponse{
		Id:       user.ID,
		Mobile:   user.Mobile,
		NickName: user.NickName,
		Gender:   user.Gender,
		Role:     user.Role,
		Password: user.Password,
	}
	if user.Birthday != nil {
		userInfoRsp.Birthday = uint64(user.Birthday.Unix())
	}
	return userInfoRsp
}

func (s *UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	fmt.Println("获取用户列表....")
	//获取用户列表
	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	//构造返回
	rsp := &proto.UserListResponse{}
	rsp.Total = int32(result.RowsAffected)
	//分页
	global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users)
	for _, user := range users {
		userInfoRsp := ModelToResponse(user)
		rsp.Data = append(rsp.Data, &userInfoRsp)
	}
	return rsp, nil
}

// GetUserByMobile 根据手机号获取用户信息
func (s *UserServer) GetUserByMobile(ctx context.Context, req *proto.UserMobileRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp, nil
}

// GetUserById 根据id获取用户信息
func (s *UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp, nil
}

// CreateUser 创建用户
func (s *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	//新建用户
	//首先查询存在否
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if result.RowsAffected > 0 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}
	//不存在则创建
	user.Mobile = req.Mobile
	user.NickName = req.NickName

	//密码
	option := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode(req.Password, option)
	newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)

	//解析
	//passwordInfo := strings.Split(my_password, "$")
	//
	//verify := password.Verify("123456", passwordInfo[2], passwordInfo[3], option)
	//fmt.Println(verify)
	user.Password = newPassword

	//保存
	result = global.DB.Create(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	//返回
	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp, nil
}

// UpdateUser
func (s *UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserInfo) (*empty.Empty, error) {
	//查询用户
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	//更新
	//生日
	birthDay := time.Unix(int64(req.Birthday), 0)

	if req.NickName != "" {
		user.NickName = req.NickName
	}
	user.Birthday = &birthDay
	user.Gender = req.Gender

	//保存
	result = global.DB.Save(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	//返回

	return &empty.Empty{}, nil
}

// CheckPassword
func (s *UserServer) CheckPassword(ctx context.Context, req *proto.CheckPasswordInfo) (*proto.CheckResponse, error) {
	//校验密码
	//解析
	option := &password.Options{16, 100, 32, sha512.New}
	passwordInfo := strings.Split(req.EncryptedPassword, "$")
	verify := password.Verify(req.Password, passwordInfo[2], passwordInfo[3], option)
	return &proto.CheckResponse{Success: verify}, nil
}
