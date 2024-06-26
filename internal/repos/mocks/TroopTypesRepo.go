// Code generated by MockGen. DO NOT EDIT.
// Source: ./troopTypes.go
//
// Generated by this command:
//
//	mockgen -source=./troopTypes.go -destination=./mocks/TroopTypesRepo.go -package=mock_repos TroopTypesRepo
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

// MockTroopTypesRepo is a mock of TroopTypesRepo interface.
type MockTroopTypesRepo struct {
	ctrl     *gomock.Controller
	recorder *MockTroopTypesRepoMockRecorder
}

// MockTroopTypesRepoMockRecorder is the mock recorder for MockTroopTypesRepo.
type MockTroopTypesRepoMockRecorder struct {
	mock *MockTroopTypesRepo
}

// NewMockTroopTypesRepo creates a new mock instance.
func NewMockTroopTypesRepo(ctrl *gomock.Controller) *MockTroopTypesRepo {
	mock := &MockTroopTypesRepo{ctrl: ctrl}
	mock.recorder = &MockTroopTypesRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTroopTypesRepo) EXPECT() *MockTroopTypesRepoMockRecorder {
	return m.recorder
}

// Find mocks base method.
func (m *MockTroopTypesRepo) Find(ctx context.Context, opts *repos.FindTroopTypeOpts) ([]*types.TroopType, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", ctx, opts)
	ret0, _ := ret[0].([]*types.TroopType)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Find indicates an expected call of Find.
func (mr *MockTroopTypesRepoMockRecorder) Find(ctx, opts any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockTroopTypesRepo)(nil).Find), ctx, opts)
}

// FindOrCreate mocks base method.
func (m *MockTroopTypesRepo) FindOrCreate(ctx context.Context, at types.CreateTroopType) (*types.TroopType, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOrCreate", ctx, at)
	ret0, _ := ret[0].(*types.TroopType)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOrCreate indicates an expected call of FindOrCreate.
func (mr *MockTroopTypesRepoMockRecorder) FindOrCreate(ctx, at any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOrCreate", reflect.TypeOf((*MockTroopTypesRepo)(nil).FindOrCreate), ctx, at)
}

// Update mocks base method.
func (m *MockTroopTypesRepo) Update(ctx context.Context, tt types.UpdateTroopType) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, tt)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockTroopTypesRepoMockRecorder) Update(ctx, tt any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockTroopTypesRepo)(nil).Update), ctx, tt)
}

// UpdateTx mocks base method.
func (m *MockTroopTypesRepo) UpdateTx(ctx context.Context, tx neo4j.ManagedTransaction, tt types.UpdateTroopType) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTx", ctx, tx, tt)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTx indicates an expected call of UpdateTx.
func (mr *MockTroopTypesRepoMockRecorder) UpdateTx(ctx, tx, tt any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTx", reflect.TypeOf((*MockTroopTypesRepo)(nil).UpdateTx), ctx, tx, tt)
}
