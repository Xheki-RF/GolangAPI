package routes_test

import (
	"testing"

	"social-network/routes"

	"github.com/txix-open/isp-kit/test"
)

func TestRoutesHaveHttpMethod(t *testing.T) {
	t.Parallel()
	_, require := test.New(t)

	endpoints := routes.EndpointDescriptors()

	for _, endpoint := range endpoints {
		require.NotEmpty(endpoint.HttpMethod)
	}
}
