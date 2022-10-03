package testing

import (
	"os"
	"testing"

	"goravel/bootstrap"
)

func TestMain(m *testing.M) {
	bootstrap.Boot()

	os.Exit(m.Run())
}

//func TestCache(t *testing.T) {
//	mockCache := mocks.Cache()
//	mockCache.On("Put", "name", "goravel", mock.Anything).Return(nil).Once()
//	mockCache.On("Get", "name", "test").Return("hwb").Once()
//
//	res := Cache()
//	assert.Equal(t, res, "hwb")
//}
//
//func TestConfig(t *testing.T) {
//	mockConfig := mocks.Config()
//	mockConfig.On("GetString", "app.name", "test").Return("hwb").Once()
//
//	res := Config()
//	assert.Equal(t, res, "hwb")
//}
//
//func TestArtisan(t *testing.T) {
//	mockConsole := mocks.Console()
//	mockConsole.On("Call", "list").Once()
//
//	assert.NotPanics(t, func() {
//		Artisan()
//	})
//}
//
//func TestOrm(t *testing.T) {
//	mockOrm, mockOrmDB, _ := mocks.Orm()
//	mockOrm.On("Query").Return(mockOrmDB)
//
//	mockOrmDB.On("Create", mock.Anything).Return(nil).Once()
//	mockOrmDB.On("Where", "id = ?", 1).Return(mockOrmDB).Once()
//	mockOrmDB.On("Find", mock.Anything).Return(nil).Once()
//
//	assert.NoError(t, Orm())
//}
//
//func TestTransaction(t *testing.T) {
//	mockOrm, _, mockOrmTransaction := mocks.Orm()
//	mockOrm.On("Transaction", mock.Anything).Return(func(txFunc func(tx orm.Transaction) error) error {
//		return txFunc(mockOrmTransaction)
//	})
//
//	var test Test
//	mockOrmTransaction.On("Create", &test).Return(func(test2 interface{}) error {
//		test2.(*Test).ID = 1
//
//		return nil
//	}).Once()
//	mockOrmTransaction.On("Where", "id = ?", uint(1)).Return(mockOrmTransaction).Once()
//	mockOrmTransaction.On("Find", mock.Anything).Return(nil).Once()
//
//	assert.NoError(t, Transaction())
//}
//
//func TestEvent(t *testing.T) {
//	mockEvent, mockTask := mocks.Event()
//	mockEvent.On("Job", mock.Anything, mock.Anything).Return(mockTask).Once()
//	mockTask.On("Dispatch").Return(nil).Once()
//
//	assert.NoError(t, Event())
//}
//
//func TestLog(t *testing.T) {
//	mocks.Log()
//
//	Log()
//}
//
//func TestQueue(t *testing.T) {
//	mockQueue, mockTask := mocks.Queue()
//	mockQueue.On("Job", mock.Anything, mock.Anything).Return(mockTask).Once()
//	mockTask.On("Dispatch").Return(nil).Once()
//
//	assert.NoError(t, Queue())
//}
