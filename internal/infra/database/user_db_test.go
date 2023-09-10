package database

import (
	"testing"

	"apis/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestUser_Create(t *testing.T) {
	db, cleanup := initDBTest(t, &entity.User{})
	defer cleanup()

	user, _ := entity.NewUser("test", "test1@gmail.com", "test123")
	userDB := NewUser(db)

	err := userDB.Create(user)
	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, user, "user should not be nil")
	assert.NotEmpty(t, user.ID, "ID should not be empty")

	var userCreated entity.User
	err = db.Where("id = ?", user.ID).First(&userCreated).Error
	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, userCreated, "user should not be nil")
	assert.Equal(t, user.ID, userCreated.ID, "ID should be the same")
	assert.Equal(t, user.Name, userCreated.Name, "name should be the same")
	assert.Equal(t, user.Email, userCreated.Email, "email should be the same")
	assert.Equal(t, user.Password, userCreated.Password, "password should be the same")
}

func TestUser_FindByEmail(t *testing.T) {
	db, cleanup := initDBTest(t, &entity.User{})
	defer cleanup()

	user, _ := entity.NewUser("test", "test@gmail.com", "test123")
	userDB := NewUser(db)

	err := userDB.Create(user)
	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, user, "user should not be nil")
	assert.NotEmpty(t, user.Email, "email should not be empty")

	userCreated, err := userDB.FindByEmail(user.Email)
	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, userCreated, "user should not be nil")
	assert.Equal(t, user.Email, userCreated.Email, "email should be the same")
}

func TestUser_FindById(t *testing.T) {
	db, cleanup := initDBTest(t, &entity.User{})
	defer cleanup()
	user, _ := entity.NewUser("test", "test@gmail.com", "test123")
	userDB := NewUser(db)

	err := userDB.Create(user)
	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, user, "user should not be nil")
	createdUser, _ := userDB.FindById(user.ID.String())
	assert.NotNil(t, createdUser, "user should not be nil")
	assert.Equal(t, user.ID, createdUser.ID, "ID should be the same")
}

func TestUser_Delete(t *testing.T) {
	db, cleanup := initDBTest(t, &entity.User{})
	defer cleanup()

	user, _ := entity.NewUser("test", "user@gmail.com", "test123")
	userDB := NewUser(db)

	err := userDB.Create(user)
	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, user, "user should not be nil")
	assert.NotEmpty(t, user.ID, "ID should not be empty")

	userCreated, err := userDB.FindById(user.ID.String())
	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, userCreated, "user should not be nil")

	err = userDB.Delete(user.ID.String())
	assert.Nil(t, err, "error should be nil")

	_, err = userDB.FindById(user.ID.String())
	assert.NotNil(t, err, "error should not be nil")
}

func TestUser_Update(t *testing.T) {
	db, cleanup := initDBTest(t, &entity.User{})
	defer cleanup()

	user, err := entity.NewUser("test", "user@gmail.com", "test123")
	if err != nil {
		t.Error(err)
	}
	userDB := NewUser(db)
	err = userDB.Create(user)
	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, user, "user should not be nil")

	user.Email = "user.new.email@gmail.com"
	err = user.SetPassword("newpassword")
	if err != nil {
		t.Error(err)
	}
	err = userDB.Update(user)
	assert.Nil(t, err, "error should be nil")
	updatededUser, err := userDB.FindById(user.ID.String())
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, user.Email, updatededUser.Email, "email should be the same")
	assert.Equal(t, user.Password, updatededUser.Password, "password should be the same")

}
