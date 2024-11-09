package sub

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/kirilngusi/go-reverse-proxy/pkg/config"
	"github.com/kirilngusi/go-reverse-proxy/pkg/config/v1/validation"
)

func init() {
	rootCmd.AddCommand(verifyCmd)
}

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify that the configures is valid",
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfgFile == "" {
			fmt.Println("frpc: the configuration file is not specified")
			return nil
		}

		cliCfg, proxyCfgs, visitorCfgs, _, err := config.LoadClientConfig(cfgFile, strictConfigMode)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		warning, err := validation.ValidateAllClientConfig(cliCfg, proxyCfgs, visitorCfgs)
		if warning != nil {
			fmt.Printf("WARNING: %v\n", warning)
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("frpc: the configuration file %s syntax is ok\n", cfgFile)
		return nil
	},
}
