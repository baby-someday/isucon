package util

import "github.com/baby-someday/isucon/pkg/interaction"

func HandleError(err error) error {
	interaction.Mark(2)
	return err
}
