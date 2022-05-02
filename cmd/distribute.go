package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/baby-someday/isucon/internal/distribute"
	"github.com/baby-someday/isucon/pkg/remote"
	"github.com/spf13/cobra"
)

const (
	FLAG_HOSTS    = "hosts"
	FLAG_SRC      = "src"
	FLAG_DST      = "dst"
	FLAG_AUTH     = "auth"
	FLAG_USER     = "user"
	FLAG_PASSWORD = "password"
	FLAG_IGNORE   = "ignore"
)

var distributeCmd = &cobra.Command{
	Use:   "distribute",
	Short: "distribute",
	Long:  `distribute`,
	Args:  validateDistributeArgs,
	Run:   runDistributeCommand,
}

func init() {
	distributeCmd.Flags().StringSlice(
		FLAG_HOSTS,
		make([]string, 0),
		"remote host",
	)
	distributeCmd.Flags().String(
		FLAG_SRC,
		"",
		"source file path",
	)
	distributeCmd.Flags().String(
		FLAG_DST,
		"",
		"dest path",
	)
	distributeCmd.Flags().String(
		FLAG_AUTH,
		remote.AUTHENTICATION_METHOD_PASSWORD,
		"ssh authentication method",
	)
	distributeCmd.Flags().String(
		FLAG_USER,
		"",
		"ssh user",
	)
	distributeCmd.Flags().String(
		FLAG_PASSWORD,
		"",
		"ssh password",
	)
	distributeCmd.Flags().StringSlice(
		FLAG_IGNORE,
		make([]string, 0),
		"files should be ignored",
	)
	distributeCmd.MarkFlagRequired(FLAG_HOSTS)
	distributeCmd.MarkFlagRequired(FLAG_SRC)
	distributeCmd.MarkFlagRequired(FLAG_DST)
	distributeCmd.MarkFlagRequired(FLAG_AUTH)
	rootCmd.AddCommand(distributeCmd)
}

func validateDistributeArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("command is required")
	}
	return nil
}

func runDistributeCommand(cmd *cobra.Command, args []string) {
	command := args[0]

	hosts, err := cmd.Flags().GetStringSlice(FLAG_HOSTS)
	if err != nil {
		log.Fatal(err.Error())
	}

	src, err := cmd.Flags().GetString(FLAG_SRC)
	if err != nil {
		log.Fatal(err.Error())
	}

	dst, err := cmd.Flags().GetString(FLAG_DST)
	if err != nil {
		log.Fatal(err.Error())
	}

	auth, err := cmd.Flags().GetString(FLAG_AUTH)
	if err != nil {
		log.Fatal(err.Error())
	}

	ignore, err := cmd.Flags().GetStringSlice(FLAG_IGNORE)
	if err != nil {
		log.Fatal(err.Error())
	}

	switch auth {
	case remote.AUTHENTICATION_METHOD_PASSWORD:
		err = distributeUsingPasswordAuthentication(
			cmd,
			hosts,
			src,
			dst,
			command,
			ignore,
		)

	case remote.AUTHENTICATION_METHOD_KEY:

	default:
		log.Fatal(fmt.Sprintf(
			"%s flag should be followings: %s, %s",
			FLAG_AUTH,
			remote.AUTHENTICATION_METHOD_PASSWORD,
			remote.AUTHENTICATION_METHOD_KEY,
		))
	}

	if err != nil {
		log.Fatal(err.Error())
	}
}

func distributeUsingPasswordAuthentication(cmd *cobra.Command, hosts []string, src, dst, command string, ignore []string) error {
	user, err := cmd.Flags().GetString(FLAG_USER)
	if err != nil {
		return err
	}
	password, err := cmd.Flags().GetString(FLAG_PASSWORD)
	if err != nil {
		return err
	}

	return distribute.DistributeUsingPasswordAuthentication(
		context.Background(),
		hosts,
		src,
		dst,
		user,
		password,
		command,
		ignore,
	)
}
