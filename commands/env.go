package commands

import (
	log "github.com/Sirupsen/logrus"
	"github.com/blablacar/green-garden/builder"
	"github.com/blablacar/green-garden/config"
	"github.com/blablacar/green-garden/work"
	"github.com/spf13/cobra"
)

func loadEnvCommands(rootCmd *cobra.Command) {
	log.WithField("path", config.GetConfig().WorkPath).Debug("Loading envs")
	work := work.NewWork(config.GetConfig().WorkPath)

	for _, f := range work.ListEnvs() {
		env := f
		envCmd := &cobra.Command{
			Use:   env,
			Short: "Run command for " + env,
		}

		checkCmd := &cobra.Command{
			Use:   "check",
			Short: "Check local units with what is running on fleet on " + env,
			Run: func(cmd *cobra.Command, args []string) {
				checkEnv(cmd, args, work, env)
			},
		}

		statusCmd := &cobra.Command{
			Use:   "status",
			Short: "Status of " + env,
			Run: func(cmd *cobra.Command, args []string) {
				statusEnv(cmd, args, work, env)
			},
		}

		fleetctlCmd := &cobra.Command{
			Use:   "fleetctl",
			Short: "Run fleetctl command on " + env,
			Run: func(cmd *cobra.Command, args []string) {
				fleetctl(cmd, args, work, env)
			},
		}

		generateCmd := &cobra.Command{
			Use:   "generate",
			Short: "Generate units for " + env,
			Run: func(cmd *cobra.Command, args []string) {
				generateEnv(cmd, args, work, env)
			},
		}
		envCmd.AddCommand(generateCmd, fleetctlCmd, checkCmd, statusCmd)

		rootCmd.AddCommand(envCmd)

		for _, g := range work.LoadEnv(env).ListServices() {
			var service = g
			var serviceCmd = &cobra.Command{
				Use:   service,
				Short: "run command for " + service + " on env " + env,
			}

			var checkCmd = &cobra.Command{
				Use:   "check",
				Short: "Check local units matches what is running on " + env + " for " + service,
				Run: func(cmd *cobra.Command, args []string) {
					checkService(cmd, args, work, env, service)
				},
			}

			var generateCmd = &cobra.Command{
				Use:   "generate [manifest...]",
				Short: "generate units for " + service + " on env " + env,
				Long:  `generate units using remote resolved or local pod/aci manifests`,
				Run: func(cmd *cobra.Command, args []string) {
					generateService(cmd, args, work, env, service)
				},
			}

			var statusCmd = &cobra.Command{
				Use:   "status [manifest...]",
				Short: "status units for " + service + " on env " + env,
				Run: func(cmd *cobra.Command, args []string) {
					statusService(cmd, args, work, env, service)
				},
			}

			var ttl string
			var lockCmd = &cobra.Command{
				Use:   "lock [message...]",
				Short: "lock " + service + " on env " + env,
				Run: func(cmd *cobra.Command, args []string) {
					lock(cmd, args, work, env, service, ttl)
				},
			}
			lockCmd.Flags().StringVarP(&ttl, "duration", "t", "1h", "lock duration")

			var unlockCmd = &cobra.Command{
				Use:   "unlock",
				Short: "unlock " + service + " on env " + env,
				Run: func(cmd *cobra.Command, args []string) {
					unLock(cmd, args, work, env, service)
				},
			}

			var updateCmd = &cobra.Command{
				Use:   "update",
				Short: "update " + service + " on env " + env,
				Run: func(cmd *cobra.Command, args []string) {
					update(cmd, args, work, env, service)
				},
			}
			updateCmd.Flags().BoolVarP(&builder.BuildFlags.All, "all", "a", false, "process all units, even up to date")
			updateCmd.Flags().BoolVarP(&builder.BuildFlags.Yes, "yes", "y", false, "process units without asking")

			serviceCmd.AddCommand(generateCmd, checkCmd, lockCmd, unlockCmd, updateCmd, statusCmd)

			envCmd.AddCommand(serviceCmd)
		}
	}
}
