package websocketmanagement

import (
	"context"

	"github.com/ideagate/client-controller/model"
	"github.com/ideagate/client-controller/usecase/workermanagement"
	v1 "github.com/ideagate/model/gen-go/client/controller/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IRouter interface {
	Switch(ctx context.Context, eventType string, request any) (any, error)
}

func NewRouter(usecaseWorkerManagement workermanagement.IWorkerManagementUsecase) IRouter {
	return &router{
		usecaseWorkerManagement: usecaseWorkerManagement,
	}
}

type router struct {
	usecaseWorkerManagement workermanagement.IWorkerManagementUsecase
}

func (r *router) Switch(ctx context.Context, eventType string, request any) (result any, err error) {
	switch eventType {
	case "worker:list":
		//// Parse request
		//_, ok := request.(v1.GetListPodRequest)
		//if !ok {
		//	return nil, errors.New("invalid request")
		//}

		// Process request
		resultListWorker, err := r.usecaseWorkerManagement.ListWorker(ctx)
		if err != nil {
			return nil, err
		}

		// Construct response
		listPods := &v1.GetListPodResponse{}
		for _, worker := range resultListWorker {
			listPods.Pods = append(listPods.Pods, &v1.Pod{
				Name:      worker.Name,
				CreatedAt: timestamppb.New(worker.CreatedAt),
				Status:    model.ConvertPodStatus(worker.Status),
			})
		}
		result = listPods
	}

	return result, nil
}

type IRoute interface {
	ParseRequest(ctx context.Context, request any) error
	Process(ctx context.Context) (any, error)
}
