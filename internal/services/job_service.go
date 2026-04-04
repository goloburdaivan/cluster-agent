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

type JobService interface {
	GetJobs(ctx context.Context, namespace string) ([]models.JobInfo, error)
	GetJob(ctx context.Context, namespace, name string) (*batchv1.Job, error)
	CreateJob(ctx context.Context, job *batchv1.Job) error
	UpdateJob(ctx context.Context, job *batchv1.Job) error
	DeleteJob(ctx context.Context, namespace, name string) error
}

type jobService struct {
	clientset kubernetes.Interface
}

func NewJobService(clientset kubernetes.Interface) JobService {
	return &jobService{
		clientset: clientset,
	}
}

func (j *jobService) GetJobs(ctx context.Context, namespace string) ([]models.JobInfo, error) {
	list, err := j.clientset.BatchV1().Jobs(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list jobs in namespace %s: %w", namespace, err)
	}

	result := make([]models.JobInfo, 0, len(list.Items))
	for _, item := range list.Items {
		completions := int32(0)
		if item.Spec.Completions != nil {
			completions = *item.Spec.Completions
		}

		result = append(result, models.JobInfo{
			Name:        item.Name,
			Namespace:   item.Namespace,
			Completions: completions,
			Succeeded:   item.Status.Succeeded,
			Failed:      item.Status.Failed,
			Active:      item.Status.Active,
			Age:         item.CreationTimestamp.Time,
		})
	}

	return result, nil
}

func (j *jobService) GetJob(ctx context.Context, namespace, name string) (*batchv1.Job, error) {
	job, err := j.clientset.BatchV1().Jobs(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get job: %w", err)
	}

	job.Kind = "Job"
	job.APIVersion = "batch/v1"
	job.ManagedFields = nil

	return job, nil
}

func (j *jobService) CreateJob(ctx context.Context, job *batchv1.Job) error {
	_, err := j.clientset.BatchV1().Jobs(job.Namespace).Create(ctx, job, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create job: %w", err)
	}
	return nil
}

func (j *jobService) UpdateJob(ctx context.Context, job *batchv1.Job) error {
	_, err := j.clientset.BatchV1().Jobs(job.Namespace).Update(ctx, job, metav1.UpdateOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to update job: %w", err)
	}
	return nil
}

func (j *jobService) DeleteJob(ctx context.Context, namespace, name string) error {
	err := j.clientset.BatchV1().Jobs(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to delete job: %w", err)
	}
	return nil
}
