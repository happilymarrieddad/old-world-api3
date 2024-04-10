// Code generated by MockGen. DO NOT EDIT.
// Source: ./items.go
//
// Generated by this command:
//
//	mockgen -source=./items.go -destination=./mocks/ItemsRepo.go -package=mock_repos ItemsRepo
//

// Package mock_repos is a generated GoMock package.
package mock_repos

import (
	context "context"
	reflect "reflect"

	repos "github.com/happilymarrieddad/old-world/api3/internal/repos"
	types "github.com/happilymarrieddad/old-world/api3/types"
	neo4j "github.com/neo4j/neo4j-go-driver/v5/neo4j"
	gomock "go.uber.org/mock/gomock"
)

// MockItemsRepo is a mock of ItemsRepo interface.
type MockItemsRepo struct {
	ctrl     *gomock.Controller
	recorder *MockItemsRepoMockRecorder
}

// MockItemsRepoMockRecorder is the mock recorder for MockItemsRepo.
type MockItemsRepoMockRecorder struct {
	mock *MockItemsRepo
}

// NewMockItemsRepo creates a new mock instance.
func NewMockItemsRepo(ctrl *gomock.Controller) *MockItemsRepo {
	mock := &MockItemsRepo{ctrl: ctrl}
	mock.recorder = &MockItemsRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockItemsRepo) EXPECT() *MockItemsRepoMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockItemsRepo) Create(ctx context.Context, itm types.CreateItem) (*types.Item, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, itm)
	ret0, _ := ret[0].(*types.Item)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockItemsRepoMockRecorder) Create(ctx, itm any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockItemsRepo)(nil).Create), ctx, itm)
}

// CreateTx mocks base method.
func (m *MockItemsRepo) CreateTx(ctx context.Context, tx neo4j.ManagedTransaction, itm types.CreateItem) (*types.Item, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTx", ctx, tx, itm)
	ret0, _ := ret[0].(*types.Item)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTx indicates an expected call of CreateTx.
func (mr *MockItemsRepoMockRecorder) CreateTx(ctx, tx, itm any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTx", reflect.TypeOf((*MockItemsRepo)(nil).CreateTx), ctx, tx, itm)
}

// Find mocks base method.
func (m *MockItemsRepo) Find(ctx context.Context, opts *repos.FindItemsOpts) ([]*types.Item, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", ctx, opts)
	ret0, _ := ret[0].([]*types.Item)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find.
func (mr *MockItemsRepoMockRecorder) Find(ctx, opts any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockItemsRepo)(nil).Find), ctx, opts)
}

// FindOrCreate mocks base method.
func (m *MockItemsRepo) FindOrCreate(ctx context.Context, itm types.CreateItem) (*types.Item, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOrCreate", ctx, itm)
	ret0, _ := ret[0].(*types.Item)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOrCreate indicates an expected call of FindOrCreate.
func (mr *MockItemsRepoMockRecorder) FindOrCreate(ctx, itm any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOrCreate", reflect.TypeOf((*MockItemsRepo)(nil).FindOrCreate), ctx, itm)
}

// FindOrCreateTx mocks base method.
func (m *MockItemsRepo) FindOrCreateTx(ctx context.Context, tx neo4j.ManagedTransaction, itm types.CreateItem) (*types.Item, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOrCreateTx", ctx, tx, itm)
	ret0, _ := ret[0].(*types.Item)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOrCreateTx indicates an expected call of FindOrCreateTx.
func (mr *MockItemsRepoMockRecorder) FindOrCreateTx(ctx, tx, itm any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOrCreateTx", reflect.TypeOf((*MockItemsRepo)(nil).FindOrCreateTx), ctx, tx, itm)
}

// FindTx mocks base method.
func (m *MockItemsRepo) FindTx(ctx context.Context, tx neo4j.ManagedTransaction, opts *repos.FindItemsOpts) ([]*types.Item, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindTx", ctx, tx, opts)
	ret0, _ := ret[0].([]*types.Item)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindTx indicates an expected call of FindTx.
func (mr *MockItemsRepoMockRecorder) FindTx(ctx, tx, opts any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindTx", reflect.TypeOf((*MockItemsRepo)(nil).FindTx), ctx, tx, opts)
}

// Get mocks base method.
func (m *MockItemsRepo) Get(ctx context.Context, id, gameID string) (*types.Item, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id, gameID)
	ret0, _ := ret[0].(*types.Item)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockItemsRepoMockRecorder) Get(ctx, id, gameID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockItemsRepo)(nil).Get), ctx, id, gameID)
}

// GetByName mocks base method.
func (m *MockItemsRepo) GetByName(ctx context.Context, gameID string, armyTypeID *string, name string) (*types.Item, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByName", ctx, gameID, armyTypeID, name)
	ret0, _ := ret[0].(*types.Item)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByName indicates an expected call of GetByName.
func (mr *MockItemsRepoMockRecorder) GetByName(ctx, gameID, armyTypeID, name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByName", reflect.TypeOf((*MockItemsRepo)(nil).GetByName), ctx, gameID, armyTypeID, name)
}

// GetByNameTx mocks base method.
func (m *MockItemsRepo) GetByNameTx(ctx context.Context, tx neo4j.ManagedTransaction, gameID string, armyTypeID *string, name string) (*types.Item, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByNameTx", ctx, tx, gameID, armyTypeID, name)
	ret0, _ := ret[0].(*types.Item)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByNameTx indicates an expected call of GetByNameTx.
func (mr *MockItemsRepoMockRecorder) GetByNameTx(ctx, tx, gameID, armyTypeID, name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByNameTx", reflect.TypeOf((*MockItemsRepo)(nil).GetByNameTx), ctx, tx, gameID, armyTypeID, name)
}

// GetTx mocks base method.
func (m *MockItemsRepo) GetTx(ctx context.Context, tx neo4j.ManagedTransaction, id, gameID string) (*types.Item, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTx", ctx, tx, id, gameID)
	ret0, _ := ret[0].(*types.Item)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTx indicates an expected call of GetTx.
func (mr *MockItemsRepoMockRecorder) GetTx(ctx, tx, id, gameID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTx", reflect.TypeOf((*MockItemsRepo)(nil).GetTx), ctx, tx, id, gameID)
}
