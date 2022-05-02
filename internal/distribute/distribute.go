package distribute

import (
	"context"
	"path"

	"github.com/baby-someday/isucon/pkg/build"
	"github.com/baby-someday/isucon/pkg/remote"
)

func DistributeUsingPasswordAuthentication(ctx context.Context, hosts []string, src, dst, user, password, command string, ignore []string) error {
	zipPath := path.Base(src) + ".zip"
	err := build.Compress(src, zipPath, ignore)
	if err != nil {
		return err
	}

	for _, host := range hosts {
		authenticationMethod := remote.PasswordAuthentication{
			User:     user,
			Password: password,
		}

		err = remote.Copy(
			ctx,
			host,
			zipPath,
			dst,
			authenticationMethod,
		)
		if err != nil {
			return err
		}

		_, err = remote.Exec(
			host,
			command,
			authenticationMethod,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
