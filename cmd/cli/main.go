package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/janwiemers/shelly_cmd"
	"github.com/spf13/cobra"
)

func prettyPrint(i interface{}) (string, error) {
	s, err := json.MarshalIndent(i, "", "  ")

	if err != nil {
		return "", err
	}

	return string(s), nil
}

func main() {
	var ip string
	var relay int
	var switchOffDelay = -1
	var shelly *shelly_cmd.RpcApi

	var cmdSwitch = &cobra.Command{
		Use:   "switch",
		Short: "returns information about the chosen switch",
		Long:  `returns information about the chosen switch`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			config, err := shelly.SwitchGetConfig(relay)
			if err != nil {
				log.Fatal(err)
			}
			c, err := prettyPrint(config)

			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("%v\n", c)
		},
	}

	var cmdSwitchOn = &cobra.Command{
		Use:   "on",
		Short: "power on switch",
		Long:  `power on switch`,
		Run: func(cmd *cobra.Command, args []string) {
			var was *shelly_cmd.SwitchWasResponse
			var err error
			
			if switchOffDelay == -1 {
				was, err = shelly.SwitchOn(relay)
			} else {
				was, err = shelly.SwitchOnWithTimer(relay, switchOffDelay)
			}

			w, err := prettyPrint(was)

			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("%v\n", w)
		},
	}

	var cmdSwitchOff = &cobra.Command{
		Use:   "off",
		Short: "power off switch",
		Long:  `power off switch`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			was, err := shelly.SwitchOff(relay)
			w, err := prettyPrint(was)

			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("%v\n", w)
		},
	}

	var cmdSwitchToggle = &cobra.Command{
		Use:   "toggle",
		Short: "toggle power for switch",
		Long:  `toggle power for switch`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			was, err := shelly.SwitchToggle(relay)
			w, err := prettyPrint(was)

			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("%v\n", w)
		},
	}

	var cmdGetStatus = &cobra.Command{
		Use:   "status",
		Short: "retrieve status of switch",
		Long:  `retrieve status of switch`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			status, err := shelly.SwitchGetStatus(relay)
			s, err := prettyPrint(status)

			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("%v\n", s)
		},
	}

	var rootCmd = &cobra.Command{
		Use: "Shelly CMD",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			var err error
			shelly, err = shelly_cmd.NewRpcApi(ip)

			if err != nil {
				log.Fatal(err)
			}
		},
	}
	rootCmd.PersistentFlags().StringVarP(&ip, "ip", "i", "", "Define the IP of the shelly for the current command")
	rootCmd.PersistentFlags().IntVarP(&relay, "relay", "r", 0, "Define the relay you want to switch on")
	cmdSwitchOn.Flags().IntVarP(&switchOffDelay, "off-delay", "o", -1, "Define the off delay in seconds. -1 for infinite")

	rootCmd.AddCommand(cmdSwitch)

	cmdSwitch.AddCommand(cmdSwitchOn, cmdSwitchOff, cmdSwitchToggle, cmdGetStatus)
	rootCmd.Execute()
}
