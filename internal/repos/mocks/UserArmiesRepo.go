// Code generated by MockGen. DO NOT EDIT.
// Source: ./userArmies.go
//
// Generated by this command:
//
//	mockgen -source=./userArmies.go -destination=./mocks/UserArmiesRepo.go -package=mock_repos UserArmiesRepo
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

// MockUserArmiesRepo is a mock of UserArmiesRepo interface.
type MockUserArmiesRepo struct {
	ctrl     *gomock.Controller
	recorder *MockUserArmiesRepoMockRecorder
}

// MockUserArmiesRepoMockRecorder is the mock recorder for MockUserArmiesRepo.
type MockUserArmiesRepoMockRecorder struct {
	mock *MockUserArmiesRepo
}

// NewMockUserArmiesRepo creates a new mock instance.
func NewMockUserArmiesRepo(ctrl *gomock.Controller) *MockUserArmiesRepo {
	mock := &MockUserArmiesRepo{ctrl: ctrl}
	mock.recorder = &MockUserArmiesRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserArmiesRepo) EXPECT() *MockUserArmiesRepoMockRecorder {
	return m.recorder
}

// AddUnits mocks base method.
func (m *MockUserArmiesRepo) AddUnits(ctx context.Context, userArmyID string, uaus ...*types.CreateUserArmyUnit) ([]*types.UserArmyUnit, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, userArmyID}
	for _, a := range uaus {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddUnits", varargs...)
	ret0, _ := ret[0].([]*types.UserArmyUnit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddUnits indicates an expected call of AddUnits.
func (mr *MockUserArmiesRepoMockRecorder) AddUnits(ctx, userArmyID any, uaus ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, userArmyID}, uaus...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUnits", reflect.TypeOf((*MockUserArmiesRepo)(nil).AddUnits), varargs...)
}

// AddUnitsTx mocks base method.
func (m *MockUserArmiesRepo) AddUnitsTx(ctx context.Context, tx neo4j.ManagedTransaction, userArmyID string, uaus ...*types.CreateUserArmyUnit) ([]*types.UserArmyUnit, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, tx, userArmyID}
	for _, a := range uaus {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddUnitsTx", varargs...)
	ret0, _ := ret[0].([]*types.UserArmyUnit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddUnitsTx indicates an expected call of AddUnitsTx.
func (mr *MockUserArmiesRepoMockRecorder) AddUnitsTx(ctx, tx, userArmyID any, uaus ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, tx, userArmyID}, uaus...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUnitsTx", reflect.TypeOf((*MockUserArmiesRepo)(nil).AddUnitsTx), varargs...)
}

// Create mocks base method.
func (m *MockUserArmiesRepo) Create(ctx context.Context, nua types.CreateUserArmy) (*types.UserArmy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, nua)
	ret0, _ := ret[0].(*types.UserArmy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockUserArmiesRepoMockRecorder) Create(ctx, nua any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUserArmiesRepo)(nil).Create), ctx, nua)
}

// CreateTx mocks base method.
func (m *MockUserArmiesRepo) CreateTx(ctx context.Context, tx neo4j.ManagedTransaction, nua types.CreateUserArmy) (*types.UserArmy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTx", ctx, tx, nua)
	ret0, _ := ret[0].(*types.UserArmy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTx indicates an expected call of CreateTx.
func (mr *MockUserArmiesRepoMockRecorder) CreateTx(ctx, tx, nua any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTx", reflect.TypeOf((*MockUserArmiesRepo)(nil).CreateTx), ctx, tx, nua)
}

// Find mocks base method.
func (m *MockUserArmiesRepo) Find(ctx context.Context, userID string, opts *repos.FindUserArmyOpts) ([]*types.UserArmy, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", ctx, userID, opts)
	ret0, _ := ret[0].([]*types.UserArmy)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Find indicates an expected call of Find.
func (mr *MockUserArmiesRepoMockRecorder) Find(ctx, userID, opts any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockUserArmiesRepo)(nil).Find), ctx, userID, opts)
}

// FindTx mocks base method.
func (m *MockUserArmiesRepo) FindTx(ctx context.Context, tx neo4j.ManagedTransaction, userID string, opts *repos.FindUserArmyOpts) ([]*types.UserArmy, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindTx", ctx, tx, userID, opts)
	ret0, _ := ret[0].([]*types.UserArmy)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// FindTx indicates an expected call of FindTx.
func (mr *MockUserArmiesRepoMockRecorder) FindTx(ctx, tx, userID, opts any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindTx", reflect.TypeOf((*MockUserArmiesRepo)(nil).FindTx), ctx, tx, userID, opts)
}

// Get mocks base method.
func (m *MockUserArmiesRepo) Get(ctx context.Context, userID, userArmyID string) (*types.UserArmy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, userID, userArmyID)
	ret0, _ := ret[0].(*types.UserArmy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockUserArmiesRepoMockRecorder) Get(ctx, userID, userArmyID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockUserArmiesRepo)(nil).Get), ctx, userID, userArmyID)
}

// GetTx mocks base method.
func (m *MockUserArmiesRepo) GetTx(ctx context.Context, tx neo4j.ManagedTransaction, userID, userArmyID string) (*types.UserArmy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTx", ctx, tx, userID, userArmyID)
	ret0, _ := ret[0].(*types.UserArmy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTx indicates an expected call of GetTx.
func (mr *MockUserArmiesRepoMockRecorder) GetTx(ctx, tx, userID, userArmyID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTx", reflect.TypeOf((*MockUserArmiesRepo)(nil).GetTx), ctx, tx, userID, userArmyID)
}

// GetUnit mocks base method.
func (m *MockUserArmiesRepo) GetUnit(ctx context.Context, userArmyUnitID string) (*types.UserArmyUnit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnit", ctx, userArmyUnitID)
	ret0, _ := ret[0].(*types.UserArmyUnit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUnit indicates an expected call of GetUnit.
func (mr *MockUserArmiesRepoMockRecorder) GetUnit(ctx, userArmyUnitID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnit", reflect.TypeOf((*MockUserArmiesRepo)(nil).GetUnit), ctx, userArmyUnitID)
}

// GetUnitTx mocks base method.
func (m *MockUserArmiesRepo) GetUnitTx(ctx context.Context, tx neo4j.ManagedTransaction, userArmyUnitID string) (*types.UserArmyUnit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnitTx", ctx, tx, userArmyUnitID)
	ret0, _ := ret[0].(*types.UserArmyUnit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUnitTx indicates an expected call of GetUnitTx.
func (mr *MockUserArmiesRepoMockRecorder) GetUnitTx(ctx, tx, userArmyUnitID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnitTx", reflect.TypeOf((*MockUserArmiesRepo)(nil).GetUnitTx), ctx, tx, userArmyUnitID)
}

// RemoveUnits mocks base method.
func (m *MockUserArmiesRepo) RemoveUnits(ctx context.Context, userArmyID string, userArmyUnitIDs ...string) error {
	m.ctrl.T.Helper()
	varargs := []any{ctx, userArmyID}
	for _, a := range userArmyUnitIDs {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RemoveUnits", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveUnits indicates an expected call of RemoveUnits.
func (mr *MockUserArmiesRepoMockRecorder) RemoveUnits(ctx, userArmyID any, userArmyUnitIDs ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, userArmyID}, userArmyUnitIDs...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveUnits", reflect.TypeOf((*MockUserArmiesRepo)(nil).RemoveUnits), varargs...)
}

// RemoveUnitsTx mocks base method.
func (m *MockUserArmiesRepo) RemoveUnitsTx(ctx context.Context, tx neo4j.ManagedTransaction, userArmyID string, userArmyUnitIDs ...string) error {
	m.ctrl.T.Helper()
	varargs := []any{ctx, tx, userArmyID}
	for _, a := range userArmyUnitIDs {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RemoveUnitsTx", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveUnitsTx indicates an expected call of RemoveUnitsTx.
func (mr *MockUserArmiesRepoMockRecorder) RemoveUnitsTx(ctx, tx, userArmyID any, userArmyUnitIDs ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, tx, userArmyID}, userArmyUnitIDs...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveUnitsTx", reflect.TypeOf((*MockUserArmiesRepo)(nil).RemoveUnitsTx), varargs...)
}

// Update mocks base method.
func (m *MockUserArmiesRepo) Update(ctx context.Context, obj types.UpdateUserArmy) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, obj)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockUserArmiesRepoMockRecorder) Update(ctx, obj any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUserArmiesRepo)(nil).Update), ctx, obj)
}

// UpdateTx mocks base method.
func (m *MockUserArmiesRepo) UpdateTx(ctx context.Context, tx neo4j.ManagedTransaction, obj types.UpdateUserArmy) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTx", ctx, tx, obj)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTx indicates an expected call of UpdateTx.
func (mr *MockUserArmiesRepoMockRecorder) UpdateTx(ctx, tx, obj any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTx", reflect.TypeOf((*MockUserArmiesRepo)(nil).UpdateTx), ctx, tx, obj)
}

// UpdateUnit mocks base method.
func (m *MockUserArmiesRepo) UpdateUnit(ctx context.Context, opts types.UpdateUserArmyUnit) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUnit", ctx, opts)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUnit indicates an expected call of UpdateUnit.
func (mr *MockUserArmiesRepoMockRecorder) UpdateUnit(ctx, opts any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUnit", reflect.TypeOf((*MockUserArmiesRepo)(nil).UpdateUnit), ctx, opts)
}

// UpdateUnitTx mocks base method.
func (m *MockUserArmiesRepo) UpdateUnitTx(ctx context.Context, tx neo4j.ManagedTransaction, opts types.UpdateUserArmyUnit) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUnitTx", ctx, tx, opts)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUnitTx indicates an expected call of UpdateUnitTx.
func (mr *MockUserArmiesRepoMockRecorder) UpdateUnitTx(ctx, tx, opts any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUnitTx", reflect.TypeOf((*MockUserArmiesRepo)(nil).UpdateUnitTx), ctx, tx, opts)
}
