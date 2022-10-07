package testing

import (
	"errors"
	"strconv"
	"testing"

	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"goravel/app/models"
	"goravel/bootstrap"
	testingmodels "goravel/testing/resources/models"
)

type OrmTestSuite struct {
	suite.Suite
	start bool
}

func TestOrmTestSuite(t *testing.T) {
	bootstrap.Boot()
	suite.Run(t, new(OrmTestSuite))
}

func (s *OrmTestSuite) SetupTest() {
	if !s.start {
		migrate(s.T())
		s.start = true
	}
}

func (s *OrmTestSuite) TestMakeMigration() {
	t := s.T()
	Equal(t, "make:migration create_users_table", "Created Migration: create_users_table")
	assert.True(t, file.Exist("./database/migrations"))
	assert.True(t, file.Remove("./database"))
}

func (s *OrmTestSuite) TestSelect() {
	t := s.T()
	user := models.User{Name: "user"}
	assert.Nil(t, facades.Orm.Query().Create(&user))
	assert.Equal(t, uint(1), user.ID)

	var user1 models.User
	assert.Nil(t, facades.Orm.Query().Where("name = ?", "user").First(&user1))
	assert.Equal(t, uint(1), user1.ID)

	var user2 models.User
	assert.Nil(t, facades.Orm.Query().Find(&user2, user.ID))
	assert.Equal(t, uint(1), user2.ID)

	var user3 []models.User
	assert.Nil(t, facades.Orm.Query().Find(&user3, []uint{user.ID}))
	assert.Equal(t, 1, len(user3))

	var user4 []models.User
	assert.Nil(t, facades.Orm.Query().Where("id in ?", []uint{user.ID}).Find(&user4))
	assert.Equal(t, 1, len(user4))

	var user5 []models.User
	assert.Nil(t, facades.Orm.Query().Where("id in ?", []uint{user.ID}).Get(&user5))
	assert.Equal(t, 1, len(user5))

	clearData(t)
}

func (s *OrmTestSuite) TestFirstOrCreate() {
	t := s.T()
	var user models.User
	assert.Nil(t, facades.Orm.Query().Where("avatar = ?", "avatar").FirstOrCreate(&user, models.User{Name: "user"}))
	assert.True(t, user.ID == 1)

	var user1 models.User
	assert.Nil(t, facades.Orm.Query().Where("avatar = ?", "avatar").FirstOrCreate(&user1, models.User{Name: "user"}, models.User{Avatar: "avatar1"}))
	assert.True(t, user1.ID == 2)
	assert.True(t, user1.Avatar == "avatar1")

	clearData(t)
}

func (s *OrmTestSuite) TestWhere() {
	t := s.T()
	user := models.User{Name: "user", Avatar: "avatar"}
	assert.Nil(t, facades.Orm.Query().Create(&user))
	assert.True(t, user.ID == 1)

	user1 := models.User{Name: "user1", Avatar: "avatar1"}
	assert.Nil(t, facades.Orm.Query().Create(&user1))
	assert.True(t, user1.ID == 2)

	var user2 []models.User
	assert.Nil(t, facades.Orm.Query().Where("name = ?", "user").OrWhere("avatar = ?", "avatar1").Find(&user2))
	assert.True(t, len(user2) == 2)

	var user3 models.User
	assert.Nil(t, facades.Orm.Query().Where("name = 'user'").Find(&user3))
	assert.True(t, user3.ID == 1)

	var user4 models.User
	assert.Nil(t, facades.Orm.Query().Where("name", "user").Find(&user4))
	assert.True(t, user4.ID == 1)

	clearData(t)
}

func (s *OrmTestSuite) TestLimit() {
	t := s.T()
	user := models.User{Name: "user", Avatar: "avatar"}
	assert.Nil(t, facades.Orm.Query().Create(&user))
	assert.True(t, user.ID == 1)

	user1 := models.User{Name: "user", Avatar: "avatar1"}
	assert.Nil(t, facades.Orm.Query().Create(&user1))
	assert.True(t, user1.ID == 2)

	var user2 []models.User
	assert.Nil(t, facades.Orm.Query().Where("name = ?", "user").Limit(1).Get(&user2))
	assert.True(t, len(user2) == 1)
	assert.True(t, user2[0].ID == 1)

	clearData(t)
}

func (s *OrmTestSuite) TestOffset() {
	t := s.T()
	user := models.User{Name: "user", Avatar: "avatar"}
	assert.Nil(t, facades.Orm.Query().Create(&user))
	assert.True(t, user.ID == 1)

	user1 := models.User{Name: "user", Avatar: "avatar1"}
	assert.Nil(t, facades.Orm.Query().Create(&user1))
	assert.True(t, user1.ID == 2)

	var user2 []models.User
	assert.Nil(t, facades.Orm.Query().Where("name = ?", "user").Offset(1).Limit(1).Get(&user2))
	assert.True(t, len(user2) == 1)
	assert.True(t, user2[0].ID == 2)

	clearData(t)
}

func (s *OrmTestSuite) TestOrder() {
	t := s.T()
	user := models.User{Name: "user", Avatar: "avatar"}
	assert.Nil(t, facades.Orm.Query().Create(&user))
	assert.True(t, user.ID == 1)

	user1 := models.User{Name: "user", Avatar: "avatar1"}
	assert.Nil(t, facades.Orm.Query().Create(&user1))
	assert.True(t, user1.ID == 2)

	var user2 []models.User
	assert.Nil(t, facades.Orm.Query().Where("name = ?", "user").Order("id desc").Order("name asc").Get(&user2))
	assert.True(t, len(user2) == 2)
	assert.True(t, user2[0].ID == 2)

	clearData(t)
}

func (s *OrmTestSuite) TestPluck() {
	t := s.T()
	user := models.User{Name: "user", Avatar: "avatar"}
	assert.Nil(t, facades.Orm.Query().Create(&user))
	assert.True(t, user.ID == 1)

	user1 := models.User{Name: "user", Avatar: "avatar1"}
	assert.Nil(t, facades.Orm.Query().Create(&user1))
	assert.True(t, user1.ID == 2)

	var avatars []string
	assert.Nil(t, facades.Orm.Query().Model(&models.User{}).Where("name = ?", "user").Pluck("avatar", &avatars))
	assert.True(t, len(avatars) == 2)
	assert.True(t, avatars[0] == "avatar")

	clearData(t)
}

func (s *OrmTestSuite) TestCount() {
	t := s.T()
	user := models.User{Name: "user", Avatar: "avatar"}
	assert.Nil(t, facades.Orm.Query().Create(&user))
	assert.True(t, user.ID == 1)

	user1 := models.User{Name: "user", Avatar: "avatar1"}
	assert.Nil(t, facades.Orm.Query().Create(&user1))
	assert.True(t, user1.ID == 2)

	var count int64
	assert.Nil(t, facades.Orm.Query().Model(&models.User{}).Where("name = ?", "user").Count(&count))
	assert.True(t, count == 2)

	var count1 int64
	assert.Nil(t, facades.Orm.Query().Table("users").Where("name = ?", "user").Count(&count1))
	assert.True(t, count1 == 2)

	clearData(t)
}

func (s *OrmTestSuite) TestSelectColumn() {
	t := s.T()
	user := models.User{Name: "user", Avatar: "avatar"}
	assert.Nil(t, facades.Orm.Query().Create(&user))
	assert.True(t, user.ID == 1)

	user1 := models.User{Name: "user", Avatar: "avatar1"}
	assert.Nil(t, facades.Orm.Query().Create(&user1))
	assert.True(t, user1.ID == 2)

	user2 := models.User{Name: "user1", Avatar: "avatar1"}
	assert.Nil(t, facades.Orm.Query().Create(&user2))
	assert.True(t, user2.ID == 3)

	type Result struct {
		Name  string
		Count string
	}
	var result []Result
	assert.Nil(t, facades.Orm.Query().Model(&models.User{}).Select("name, count(avatar) as count").Group("name").Get(&result))
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "user", result[0].Name)
	assert.Equal(t, "2", result[0].Count)
	assert.Equal(t, "user1", result[1].Name)
	assert.Equal(t, "1", result[1].Count)

	var result1 []Result
	assert.Nil(t, facades.Orm.Query().Model(&models.User{}).Select("name, count(avatar) as count").Group("name").Having("name = ?", "user").Get(&result1))

	assert.Equal(t, 1, len(result1))
	assert.Equal(t, "user", result1[0].Name)
	assert.Equal(t, "2", result1[0].Count)

	clearData(t)
}

func (s *OrmTestSuite) TestJoin() {
	t := s.T()
	user := models.User{Name: "user", Avatar: "avatar"}
	assert.Nil(t, facades.Orm.Query().Create(&user))
	assert.True(t, user.ID == 1)

	userAddress := testingmodels.UserAddress{UserId: user.ID, Name: "address", Province: "province"}
	assert.Nil(t, facades.Orm.Query().Create(&userAddress))
	assert.True(t, userAddress.ID == 1)

	type Result struct {
		UserName        string
		UserAddressName string
	}
	var result []Result
	assert.Nil(t, facades.Orm.Query().Model(&models.User{}).Join("left join user_addresses ua on users.id = ua.user_id").
		Select("users.name user_name, ua.name user_address_name").Get(&result))
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "user", result[0].UserName)
	assert.Equal(t, "address", result[0].UserAddressName)

	clearData(t)
}

func (s *OrmTestSuite) TestUpdate() {
	t := s.T()
	user := models.User{Name: "user", Avatar: "avatar"}
	assert.Nil(t, facades.Orm.Query().Create(&user))
	assert.True(t, user.ID == 1)

	user.Name = "user1"
	assert.Nil(t, facades.Orm.Query().Save(&user))
	assert.Nil(t, facades.Orm.Query().Model(&models.User{}).Where("id = ?", user.ID).Update("avatar", "avatar1"))

	var user1 models.User
	assert.Nil(t, facades.Orm.Query().Find(&user1, user.ID))
	assert.Equal(t, "user1", user1.Name)
	assert.Equal(t, "avatar1", user1.Avatar)

	clearData(t)
}

func (s *OrmTestSuite) TestDelete() {
	t := s.T()
	user := models.User{Name: "user", Avatar: "avatar"}
	assert.Nil(t, facades.Orm.Query().Create(&user))
	assert.Equal(t, uint(1), user.ID)

	assert.Nil(t, facades.Orm.Query().Delete(&user))

	var user1 models.User
	assert.Nil(t, facades.Orm.Query().Find(&user1, user.ID))
	assert.Equal(t, uint(0), user1.ID)

	user2 := models.User{Name: "user", Avatar: "avatar"}
	assert.Nil(t, facades.Orm.Query().Create(&user2))
	assert.Equal(t, uint(2), user2.ID)

	assert.Nil(t, facades.Orm.Query().Delete(&models.User{}, user2.ID))

	var user3 models.User
	assert.Nil(t, facades.Orm.Query().Find(&user3, user2.ID))
	assert.Equal(t, uint(0), user3.ID)

	users := []models.User{{Name: "user", Avatar: "avatar"}, {Name: "user1", Avatar: "avatar1"}}
	assert.Nil(t, facades.Orm.Query().Create(&users))
	assert.Equal(t, uint(3), users[0].ID)
	assert.Equal(t, uint(4), users[1].ID)

	assert.Nil(t, facades.Orm.Query().Delete(&models.User{}, []uint{users[0].ID, users[1].ID}))

	var count int64
	assert.Nil(t, facades.Orm.Query().Model(&models.User{}).Where("name", "user").OrWhere("name", "user1").Count(&count))
	assert.True(t, count == 0)

	clearData(t)
}

func (s *OrmTestSuite) TestSoftDelete() {
	t := s.T()
	user := models.User{Name: "user", Avatar: "avatar"}
	assert.Nil(t, facades.Orm.Query().Create(&user))
	assert.Equal(t, uint(1), user.ID)

	assert.Nil(t, facades.Orm.Query().Where("name = ?", "user").Delete(&models.User{}))

	var user1 models.User
	assert.Nil(t, facades.Orm.Query().Find(&user1, user.ID))
	assert.Equal(t, uint(0), user1.ID)

	var user2 models.User
	assert.Nil(t, facades.Orm.Query().WithTrashed().Find(&user2, user.ID))
	assert.Equal(t, uint(1), user2.ID)

	assert.Nil(t, facades.Orm.Query().Where("name = ?", "user").ForceDelete(&models.User{}))

	var user3 models.User
	assert.Nil(t, facades.Orm.Query().WithTrashed().Find(&user3, user.ID))
	assert.Equal(t, uint(0), user3.ID)

	clearData(t)
}

func (s *OrmTestSuite) TestRaw() {
	t := s.T()
	user := models.User{Name: "user", Avatar: "avatar"}
	assert.Nil(t, facades.Orm.Query().Create(&user))
	assert.Equal(t, uint(1), user.ID)

	var user1 models.User
	assert.Nil(t, facades.Orm.Query().Raw("SELECT id, name FROM users WHERE name = ?", "user").Scan(&user1))
	assert.Equal(t, uint(1), user1.ID)
	assert.Equal(t, "user", user1.Name)
	assert.Equal(t, "", user1.Avatar)

	clearData(t)
}

func (s *OrmTestSuite) TestTransactionSuccess() {
	t := s.T()
	assert.Nil(t, facades.Orm.Transaction(func(tx orm.Transaction) error {
		user := models.User{Name: "user", Avatar: "avatar"}
		assert.Nil(t, tx.Create(&user))

		user1 := models.User{Name: "user1", Avatar: "avatar1"}
		assert.Nil(t, tx.Create(&user1))

		return nil
	}))

	var user2, user3 models.User
	assert.Nil(t, facades.Orm.Query().Find(&user2, 1))
	assert.Nil(t, facades.Orm.Query().Find(&user3, 2))

	clearData(t)
}

func (s *OrmTestSuite) TestTransactionError() {
	t := s.T()
	assert.NotNil(t, facades.Orm.Transaction(func(tx orm.Transaction) error {
		user := models.User{Name: "user", Avatar: "avatar"}
		assert.Nil(t, tx.Create(&user))

		user1 := models.User{Name: "user1", Avatar: "avatar1"}
		assert.Nil(t, tx.Create(&user1))

		return errors.New("error")
	}))

	var users []models.User
	assert.Nil(t, facades.Orm.Query().Find(&users))
	assert.Equal(t, 0, len(users))

	clearData(t)
}

func (s *OrmTestSuite) TestScope() {
	t := s.T()
	users := []models.User{{Name: "user", Avatar: "avatar"}, {Name: "user1", Avatar: "avatar1"}}
	assert.Nil(t, facades.Orm.Query().Create(&users))
	assert.Equal(t, uint(1), users[0].ID)
	assert.Equal(t, uint(2), users[1].ID)

	var users1 []models.User
	assert.Nil(t, facades.Orm.Query().Scopes(paginator("1", "1")).Find(&users1))

	assert.Equal(t, 1, len(users1))
	assert.Equal(t, uint(1), users1[0].ID)

	clearData(t)
}

func migrate(t *testing.T) {
	clearTables(t)

	outStr, errStr, err := RunCommand("cp -R stubs/database database")
	assert.Empty(t, outStr)
	assert.Empty(t, errStr)
	assert.NoError(t, err)

	Equal(t, "migrate", "Migration success")
	user := models.User{Name: "user"}
	assert.Nil(t, facades.Orm.Query().Create(&user))
	assert.Nil(t, facades.Orm.Query().Create(&testingmodels.UserAddress{Name: "address", UserId: user.ID}))

	clearData(t)
	file.Remove("./database")
}

func paginator(page string, limit string) func(methods orm.Query) orm.Query {
	return func(query orm.Query) orm.Query {
		page, _ := strconv.Atoi(page)
		limit, _ := strconv.Atoi(limit)
		offset := (page - 1) * limit

		return query.Offset(offset).Limit(limit)
	}
}

func clearTables(t *testing.T) {
	facades.Orm.Query().Exec("DROP TABLE users")
	facades.Orm.Query().Exec("DROP TABLE user_addresses;")
	facades.Orm.Query().Exec("DROP TABLE migrations;")
}

func clearData(t *testing.T) {
	assert.Nil(t, facades.Orm.Query().Exec("TRUNCATE table users"))
	assert.Nil(t, facades.Orm.Query().Exec("TRUNCATE table user_addresses"))
}
