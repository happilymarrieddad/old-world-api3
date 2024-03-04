// Code generated by MockGen. DO NOT EDIT.
// Source: ./compositionTypes.go
//
// Generated by this command:
//
//	mockgen -source=./compositionTypes.go -destination=./mocks/CompositionTypesRepo.go -package=mock_repos CompositionTypesRepo
//

// Package mock_repos is a generated GoMock package.
package mock_repos

import (
	context "context"
	reflect "reflect"

	types "github.com/happilymarrieddad/old-world/api3/types"
	gomock "go.uber.org/mock/gomock"
)

// MockCompositionTypesRepo is a mock of CompositionTypesRepo interface.
type MockCompositionTypesRepo struct {
	ctrl     *gomock.Controller
	recorder *MockCompositionTypesRepoMockRecorder
}

// MockCompositionTypesRepoMockRecorder is the mock recorder for MockCompositionTypesRepo.
type MockCompositionTypesRepoMockRecorder struct {
	mock *MockCompositionTypesRepo
}

// NewMockCompositionTypesRepo creates a new mock instance.
func NewMockCompositionTypesRepo(ctrl *gomock.Controller) *MockCompositionTypesRepo {
	mock := &MockCompositionTypesRepo{ctrl: ctrl}
	mock.recorder = &MockCompositionTypesRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCompositionTypesRepo) EXPECT() *MockCompositionTypesRepoMockRecorder {
	return m.recorder
}

// Find mocks base method.
func (m *MockCompositionTypesRepo) Find(ctx context.Context, gameID string, limit, offset int) ([]*types.CompositionType, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", ctx, gameID, limit, offset)
	ret0, _ := ret[0].([]*types.CompositionType)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find.
func (mr *MockCompositionTypesRepoMockRecorder) Find(ctx, gameID, limit, offset any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockCompositionTypesRepo)(nil).Find), ctx, gameID, limit, offset)
}

// FindOrCreate mocks base method.
func (m *MockCompositionTypesRepo) FindOrCreate(ctx context.Context, at types.CreateCompositionType) (*types.CompositionType, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOrCreate", ctx, at)
	ret0, _ := ret[0].(*types.CompositionType)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOrCreate indicates an expected call of FindOrCreate.
func (mr *MockCompositionTypesRepoMockRecorder) FindOrCreate(ctx, at any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOrCreate", reflect.TypeOf((*MockCompositionTypesRepo)(nil).FindOrCreate), ctx, at)
}