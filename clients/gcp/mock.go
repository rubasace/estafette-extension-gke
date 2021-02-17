// Code generated by MockGen. DO NOT EDIT.
// Source: client.go

// Package gcp is a generated GoMock package.
package gcp

import (
	context "context"
	api "github.com/estafette/estafette-extension-gke/api"
	gomock "github.com/golang/mock/gomock"
	container "google.golang.org/api/container/v1beta1"
	reflect "reflect"
)

// MockClient is a mock of Client interface
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// LoadGKEClusterKubeConfig mocks base method
func (m *MockClient) LoadGKEClusterKubeConfig(ctx context.Context, credential *api.GKECredentials) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadGKEClusterKubeConfig", ctx, credential)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadGKEClusterKubeConfig indicates an expected call of LoadGKEClusterKubeConfig
func (mr *MockClientMockRecorder) LoadGKEClusterKubeConfig(ctx, credential interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadGKEClusterKubeConfig", reflect.TypeOf((*MockClient)(nil).LoadGKEClusterKubeConfig), ctx, credential)
}

// GetGKECluster mocks base method
func (m *MockClient) GetGKECluster(ctx context.Context, projectID, location, clusterID string) (*container.Cluster, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGKECluster", ctx, projectID, location, clusterID)
	ret0, _ := ret[0].(*container.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGKECluster indicates an expected call of GetGKECluster
func (mr *MockClientMockRecorder) GetGKECluster(ctx, projectID, location, clusterID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGKECluster", reflect.TypeOf((*MockClient)(nil).GetGKECluster), ctx, projectID, location, clusterID)
}