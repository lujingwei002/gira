package proj

import (
	"context"
	"path"

	"github.com/lujingwei002/gira"
)

func LoadResource(ctx context.Context, loader gira.ResourceLoader, config *gira.Config) error {
	// 连接resourcedb
	if resourceDbClient, err := NewResourceDbClient(ctx, config); err != nil {
		return err
	} else if err := loader.LoadResource(ctx, resourceDbClient, path.Join(Config.ResourceDir, "conf"), true); err != nil {
		return err
	} else {
		return nil
	}
}
