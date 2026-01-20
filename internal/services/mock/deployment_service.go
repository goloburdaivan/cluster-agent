package mock

import (
	"cluster-agent/internal/models"
	"context"
	"github.com/stretchr/testify/mock"
	v1 "k8s.io/api/apps/v1"
)

type DeploymentServiceMock struct {
	mock.Mock
}

func (d *DeploymentServiceMock) GetDeployments(ctx context.Context, namespace string) ([]models.DeploymentInfo, error) {
	args := d.Called(ctx, namespace)
	return args.Get(0).([]models.DeploymentInfo), args.Error(1)
}

func (d *DeploymentServiceMock) GetDeployment(ctx context.Context, namespace string, deploymentName string) (*v1.Deployment, error) {
	args := d.Called(ctx, namespace, deploymentName)
	return args.Get(0).(*v1.Deployment), args.Error(1)
}

func (d *DeploymentServiceMock) CreateDeployment(ctx context.Context, deployment *v1.Deployment) error {
	args := d.Called(ctx, deployment)
	return args.Error(0)
}

func (d *DeploymentServiceMock) DeleteDeployment(ctx context.Context, namespace string, deploymentName string) error {
	args := d.Called(ctx, namespace, deploymentName)
	return args.Error(0)
}

func (d *DeploymentServiceMock) ScaleDeployment(ctx context.Context, params models.ScaleDeploymentParams) error {
	args := d.Called(ctx, params)
	return args.Error(0)
}
