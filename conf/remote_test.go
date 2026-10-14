package conf_test

import (
	"testing"

	"social-network/conf"

	"github.com/txix-open/isp-kit/test/rct"
)

func TestDefaultRemoteConfig(t *testing.T) {
	t.Parallel()
	rct.Test(t, "default_remote_config.json", conf.Remote{})
}
