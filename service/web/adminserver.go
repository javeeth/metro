package web

import (
	"context"

	"github.com/razorpay/metro/internal/subscription"
	"github.com/razorpay/metro/internal/topic"

	"github.com/razorpay/metro/internal/merror"
	"github.com/razorpay/metro/internal/project"
	"github.com/razorpay/metro/pkg/logger"
	metrov1 "github.com/razorpay/metro/rpc/proto/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type adminServer struct {
	projectCore      project.ICore
	subscriptionCore subscription.ICore
	topicCore        topic.ICore
}

func newAdminServer(projectCore project.ICore, subscriptionCore subscription.ICore, topicCore topic.ICore) *adminServer {
	return &adminServer{projectCore, subscriptionCore, topicCore}
}

// CreateProject creates a new project
func (s adminServer) CreateProject(ctx context.Context, req *metrov1.Project) (*metrov1.Project, error) {
	logger.Ctx(ctx).Infow("request received to create project", "id", req.ProjectId)
	p, err := project.GetValidatedModelForCreate(ctx, req)
	if err != nil {
		return nil, merror.ToGRPCError(err)
	}
	err = s.projectCore.CreateProject(ctx, p)
	if err != nil {
		return nil, merror.ToGRPCError(err)
	}
	return req, nil
}

// DeleteProject creates a new project
func (s adminServer) DeleteProject(ctx context.Context, req *metrov1.Project) (*emptypb.Empty, error) {
	logger.Ctx(ctx).Infow("request received to delete project", "id", req.ProjectId)
	p, err := project.GetValidatedModelForDelete(ctx, req)
	if err != nil {
		return nil, merror.ToGRPCError(err)
	}

	// Delete all subscriptions of the project first
	err = s.subscriptionCore.DeleteProjectSubscriptions(ctx, p.ProjectID)
	if err != nil {
		return nil, merror.ToGRPCError(err)
	}

	// Delete all topics of the project
	err = s.topicCore.DeleteProjectTopics(ctx, p.ProjectID)
	if err != nil {
		return nil, merror.ToGRPCError(err)
	}

	// Delete the project
	err = s.projectCore.DeleteProject(ctx, p)
	if err != nil {
		return nil, merror.ToGRPCError(err)
	}

	return &emptypb.Empty{}, nil
}
