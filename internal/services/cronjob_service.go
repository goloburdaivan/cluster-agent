package services

import (
	"cluster-agent/internal/models"
	"context"
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type CronJobService interface {
	GetCronJobs(ctx context.Context, namespace string) ([]models.CronJobInfo, error)
	GetCronJob(ctx context.Context, namespace, name string) (*batchv1.CronJob, error)
	CreateCronJob(ctx context.Context, cronjob *batchv1.CronJob) error
	UpdateCronJob(ctx context.Context, cronjob *batchv1.CronJob) error
	DeleteCronJob(ctx context.Context, namespace, name string) error
}

type cronJobService struct {
	clientset kubernetes.Interface
}

func NewCronJobService(clientset kubernetes.Interface) CronJobService {
	return &cronJobService{
		clientset: clientset,
	}
}

func (c *cronJobService) GetCronJobs(ctx context.Context, namespace string) ([]models.CronJobInfo, error) {
	list, err := c.clientset.BatchV1().CronJobs(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list cronjobs in namespace %s: %w", namespace, err)
	}

	result := make([]models.CronJobInfo, 0, len(list.Items))
	for _, item := range list.Items {
		suspended := false
		if item.Spec.Suspend != nil {
			suspended = *item.Spec.Suspend
		}

		var lastSchedule *metav1.Time
		if item.Status.LastScheduleTime != nil {
			lastSchedule = item.Status.LastScheduleTime
		}

		result = append(result, models.CronJobInfo{
			Name:         item.Name,
			Namespace:    item.Namespace,
			Schedule:     item.Spec.Schedule,
			Suspended:    suspended,
			Active:       len(item.Status.Active),
			LastSchedule: lastSchedule,
			Age:          item.CreationTimestamp.Time,
		})
	}

	return result, nil
}

func (c *cronJobService) GetCronJob(ctx context.Context, namespace, name string) (*batchv1.CronJob, error) {
	cronjob, err := c.clientset.BatchV1().CronJobs(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get cronjob: %w", err)
	}

	cronjob.Kind = "CronJob"
	cronjob.APIVersion = "batch/v1"
	cronjob.ManagedFields = nil

	return cronjob, nil
}

func (c *cronJobService) CreateCronJob(ctx context.Context, cronjob *batchv1.CronJob) error {
	_, err := c.clientset.BatchV1().CronJobs(cronjob.Namespace).Create(ctx, cronjob, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create cronjob: %w", err)
	}
	return nil
}

func (c *cronJobService) UpdateCronJob(ctx context.Context, cronjob *batchv1.CronJob) error {
	_, err := c.clientset.BatchV1().CronJobs(cronjob.Namespace).Update(ctx, cronjob, metav1.UpdateOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to update cronjob: %w", err)
	}
	return nil
}

func (c *cronJobService) DeleteCronJob(ctx context.Context, namespace, name string) error {
	err := c.clientset.BatchV1().CronJobs(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to delete cronjob: %w", err)
	}
	return nil
}
