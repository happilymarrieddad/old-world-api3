// Code generated by MockGen. DO NOT EDIT.
// Source: ./armyTypes.go
//
// Generated by this command:
//
//	mockgen -source=./armyTypes.go -destination=./mocks/ArmyTypesRepo.go -package=mock_repos ArmyTypesRepo
//

// Package mock_repos is a generated GoMock package.
package mock_repos

import (
	context "context"
	reflect "reflect"

	types "github.com/happilymarrieddad/old-world/api3/types"
	gomock "go.uber.org/mock/gomock"
)

// MockArmyTypesRepo is a mock of ArmyTypesRepo interface.
type MockArmyTypesRepo struct {
	ctrl     *gomock.Controller
	recorder *MockArmyTypesRepoMockRecorder
}

// MockArmyTypesRepoMockRecorder is the mock recorder for MockArmyTypesRepo.
type MockArmyTypesRepoMockRecorder struct {
	mock *MockArmyTypesRepo
}

// NewMockArmyTypesRepo creates a new mock instance.
func NewMockArmyTypesRepo(ctrl *gomock.Controller) *MockArmyTypesRepo {
	mock := &MockArmyTypesRepo{ctrl: ctrl}
	mock.recorder = &MockArmyTypesRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockArmyTypesRepo) EXPECT() *MockArmyTypesRepoMockRecorder {
	return m.recorder
}

// Find mocks base method.
func (m *MockArmyTypesRepo) Find(ctx context.Context, gameID string, limit, offset int) ([]*types.ArmyType, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", ctx, gameID, limit, offset)
	ret0, _ := ret[0].([]*types.ArmyType)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Find indicates an expected call of Find.
func (mr *MockArmyTypesRepoMockRecorder) Find(ctx, gameID, limit, offset any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockArmyTypesRepo)(nil).Find), ctx, gameID, limit, offset)
}

// FindOrCreate mocks base method.
func (m *MockArmyTypesRepo) FindOrCreate(ctx context.Context, at types.CreateArmyType) (*types.ArmyType, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOrCreate", ctx, at)
	ret0, _ := ret[0].(*types.ArmyType)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOrCreate indicates an expected call of FindOrCreate.
func (mr *MockArmyTypesRepoMockRecorder) FindOrCreate(ctx, at any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOrCreate", reflect.TypeOf((*MockArmyTypesRepo)(nil).FindOrCreate), ctx, at)
}
