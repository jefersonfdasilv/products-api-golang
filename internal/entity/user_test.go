package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewUser(t *testing.T) {
	type args struct {
		name     string
		email    string
		password string
	}
	tests := []struct {
		name string
		args args
	}{
		{"test1", args{"test1", "test1@a.com", "test1"}},
		{"test2", args{"test2", "test2@b.com", "test2"}},
		{"test3", args{"test3", "test3@a.com", "test3"}},
		{"test4", args{"test4", "test4@a.com", "test4"}},
	}
	for _, tt := range tests {
		user, err := NewUser(tt.args.name, tt.args.email, tt.args.password)
		assert.Nil(t, err, "error should be nil")
		assert.NotNilf(t, user, "user should not be nil")
		assert.NotEmpty(t, user.ID, "ID should not be empty")
		assert.NotEmpty(t, user.Name, "name should not be empty")
		assert.NotEmpty(t, user.Email, "email should not be empty")
		assert.Equal(t, tt.args.name, user.Name, "name should be the same")
		assert.Equal(t, tt.args.email, user.Email, "email should be the same")
		assert.True(t, user.CheckPassword(tt.args.password), "password should be the same")
	}
}

func TestUser_ValidatePassword(t *testing.T) {
	type args struct {
		name     string
		email    string
		password string
	}
	tests := []struct {
		name string
		args args
	}{
		{"test1", args{"test1", "test1@a.com", "test1"}},
		{"test2", args{"test2", "test2@b.com", "test2"}},
	}

	for _, tt := range tests {
		user, err := NewUser(tt.args.name, tt.args.email, tt.args.password)
		assert.Nil(t, err)
		assert.NotNil(t, user, "user should not be nil")
		assert.True(t, user.CheckPassword(tt.args.password), "password should be the same")
		assert.False(t, user.CheckPassword(""), "password should not be the same")
		assert.NotEqual(t, tt.args.password, user.Password, "password should not be the same")
	}
}
