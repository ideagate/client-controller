package workermanagement

import (
	"context"

	"github.com/bayu-aditya/ideagate/backend/client/controller/adapter/kubernetes"
)

type usecase struct {
	adapterKubernetes kubernetes.IKubernetesAdapter
	workerName        string // worker name worker label, deployment and container
}

func New(adapterKubernetes kubernetes.IKubernetesAdapter) IWorkerManagementUsecase {
	return &usecase{
		adapterKubernetes: adapterKubernetes,
		workerName:        "worker-api",
	}
}

func (u *usecase) CreateWorker(ctx context.Context, version string) error {
	workerDeploymentSpec := constructWorkerDeploymentSpec(u.workerName, version)

	if err := u.adapterKubernetes.CreateDeployment(ctx, workerDeploymentSpec); err != nil {
		return err
	}

	return nil
}

func (u *usecase) ListWorker(ctx context.Context) ([]ListWorkerOutputItem, error) {
	pods, err := u.adapterKubernetes.ListPods(ctx, u.workerName)
	if err != nil {
		return nil, err
	}

	var output []ListWorkerOutputItem
	for _, pod := range pods {
		// determine worker image tag
		workerImageTag := ""
		for _, container := range pod.Spec.Containers {
			if container.Name == u.workerName {
				workerImageTag = container.Image
				break
			}
		}

		output = append(output, ListWorkerOutputItem{
			Name:      pod.Name,
			Version:   workerImageTag,
			CreatedAt: pod.CreationTimestamp.Time,
			Status:    string(pod.Status.Phase),
		})
	}

	return output, nil
}

func (u *usecase) SetWorkerNum(ctx context.Context, num int32) error {
	if err := u.adapterKubernetes.SetReplicas(ctx, u.workerName, num); err != nil {
		return err
	}

	return nil
}

func (u *usecase) RestartWorker(ctx context.Context) error {
	if err := u.adapterKubernetes.RestartDeployment(ctx, u.workerName); err != nil {
		return err
	}

	return nil
}
