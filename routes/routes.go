package routes

import (
	"net/http"

	"github.com/txix-open/isp-kit/cluster"
	"github.com/txix-open/isp-kit/grpc"
	"github.com/txix-open/isp-kit/grpc/endpoint"
	httpEndpoint "github.com/txix-open/isp-kit/http/endpoint"
	"github.com/txix-open/isp-kit/http/router"

	"social-network/controller"
)

type Controllers struct {
	User controller.User
}

func EndpointDescriptors() []cluster.EndpointDescriptor {
	return endpointDescriptors(Controllers{})
}

func Handler(wrapper endpoint.Wrapper, c Controllers) *grpc.Mux {
	muxer := grpc.NewMux()
	for _, descriptor := range endpointDescriptors(c) {
		muxer.Handle(descriptor.Path, wrapper.Endpoint(descriptor.Handler))
	}
	return muxer
}

func HttpHandler(wrapper httpEndpoint.Wrapper, c Controllers) *router.Router {
	r := router.New()

	r.POST("/users/register", wrapper.Endpoint(c.User.RegisterUserCtrl))
	r.POST("/users/update", wrapper.Endpoint(c.User.UpdateUserCtrl))
	r.POST("/users/delete", wrapper.Endpoint(c.User.DeleteUserCtrl))

	return r
}

func endpointDescriptors(c Controllers) []cluster.EndpointDescriptor {
	return []cluster.EndpointDescriptor{{
		Path:             "social-network/users/register",
		Inner:            false,
		UserAuthRequired: false,
		Handler:          c.User.RegisterUserCtrl,
		HttpMethod:       http.MethodPost,
	}, {
		Path:             "social-network/users/update",
		Inner:            false,
		UserAuthRequired: false,
		Handler:          c.User.UpdateUserCtrl,
		HttpMethod:       http.MethodPost,
	}, {
		Path:             "social-network/users/delete",
		Inner:            false,
		UserAuthRequired: false,
		Handler:          c.User.DeleteUserCtrl,
		HttpMethod:       http.MethodPost,
	}}
}
