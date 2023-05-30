package service

import (
	userv1alpha "github.com/MohammadBnei/go-realtime-chat/server/stubs/user/v1alpha"
	funk "github.com/thoas/go-funk"
)

type UserService interface {
	AddUser(roomid string, user *userv1alpha.User)
	DeleteUser(userId, roomid string)
	GetUser(userId, roomid string) *userv1alpha.User
	GetUsers(roomid string) []*userv1alpha.User
}

type userService struct {
	users map[string]map[string]*userv1alpha.User
}

func NewUserService() UserService {
	return &userService{
		users: make(map[string]map[string]*userv1alpha.User),
	}
}

func (s *userService) AddUser(roomid string, user *userv1alpha.User) {
	if _, ok := s.users[roomid]; !ok {
		s.users[roomid] = make(map[string]*userv1alpha.User)
	}

	s.users[roomid][user.Id] = user
}

func (s *userService) DeleteUser(userId, roomid string) {
	if _, ok := s.users[roomid]; !ok {
		return
	}

	delete(s.users[roomid], userId)
}

func (s *userService) GetUser(userId, roomid string) *userv1alpha.User {
	if _, ok := s.users[roomid]; !ok {
		return nil
	}

	return s.users[roomid][userId]
}

func (s *userService) GetUsers(roomid string) []*userv1alpha.User {
	if _, ok := s.users[roomid]; !ok {
		return nil
	}

	return funk.Map(s.users[roomid], func(user *userv1alpha.User) *userv1alpha.User {
		return user
	}).([]*userv1alpha.User)
}
