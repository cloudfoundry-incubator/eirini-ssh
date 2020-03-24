package main

import (
	"fmt"
	"os"

	extension "github.com/SUSE/eirini-ssh/extension"
	"github.com/SUSE/eirini-ssh/version"
	eirinix "github.com/SUSE/eirinix"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc" // from https://github.com/kubernetes/client-go/issues/345
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the eirini extension",
	PreRun: func(cmd *cobra.Command, args []string) {

		viper.BindPFlag("kubeconfig", cmd.Flags().Lookup("kubeconfig"))
		viper.BindPFlag("namespace", cmd.Flags().Lookup("namespace"))
		viper.BindPFlag("operator-webhook-host", cmd.Flags().Lookup("operator-webhook-host"))
		viper.BindPFlag("operator-webhook-port", cmd.Flags().Lookup("operator-webhook-port"))
		viper.BindPFlag("operator-service-name", cmd.Flags().Lookup("operator-service-name"))
		viper.BindPFlag("operator-webhook-namespace", cmd.Flags().Lookup("operator-webhook-namespace"))
		viper.BindPFlag("register", cmd.Flags().Lookup("register"))

		viper.BindEnv("kubeconfig")
		viper.BindEnv("namespace", "NAMESPACE")
		viper.BindEnv("operator-webhook-host", "OPERATOR_WEBHOOK_HOST")
		viper.BindEnv("operator-webhook-port", "OPERATOR_WEBHOOK_PORT")
		viper.BindEnv("operator-service-name", "OPERATOR_SERVICE_NAME")
		viper.BindEnv("operator-webhook-namespace", "OPERATOR_WEBHOOK_NAMESPACE")
		viper.BindEnv("register", "EIRINI_EXTENSION_REGISTER")
	},
	Run: func(cmd *cobra.Command, args []string) {
		defer log.Sync()
		kubeConfig := viper.GetString("kubeconfig")

		namespace := viper.GetString("namespace")

		log.Infof("Starting %s with namespace %s", version.Version, namespace)

		webhookHost := viper.GetString("operator-webhook-host")
		webhookPort := viper.GetInt32("operator-webhook-port")
		serviceName := viper.GetString("operator-service-name")
		webhookNamespace := viper.GetString("operator-webhook-namespace")
		register := viper.GetBool("register")

		if webhookHost == "" {
			log.Fatal("required flag 'operator-webhook-host' not set (env variable: OPERATOR_WEBHOOK_HOST)")
		}

		RegisterWebhooks := true
		if !register {
			log.Info("The extension will start without registering")
			RegisterWebhooks = false
		}

		filterEiriniApps := true
		x := eirinix.NewManager(
			eirinix.ManagerOptions{
				FilterEiriniApps:    &filterEiriniApps,
				OperatorFingerprint: "eirini-ssh",
				Namespace:           namespace,
				Host:                webhookHost,
				Port:                webhookPort,
				KubeConfig:          kubeConfig,
				ServiceName:         serviceName,
				WebhookNamespace:    webhookNamespace,
				RegisterWebHook:     &RegisterWebhooks,
			})

		x.AddExtension(&extension.SSH{Namespace: namespace})
		x.AddWatcher(&extension.CleanupWatcher{})

		var v chan error
		go func() {
			fmt.Println("Starting watcher")
			err := x.Watch()
			if err != nil {
				v <- err
				fmt.Println(err.Error())
				return
			}
		}()
		go func() {
			fmt.Println("Starting extension")
			err := x.Start()
			if err != nil {
				v <- err
				fmt.Println(err.Error())
				return
			}
		}()

		for err := range v {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	startCmd.Flags().BoolP("register", "r", true, "Register the extension")

	rootCmd.AddCommand(startCmd)
}
