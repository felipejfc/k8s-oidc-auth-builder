/*
 * Copyright (c) 2016 Felipe Cavalcanti <fjfcavalcanti@gmail.com>
 * Author: Felipe Cavalcanti <fjfcavalcanti@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "Software"), to deal in
 * the Software without restriction, including without limitation the rights to
 * use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 * the Software, and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 * FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 * IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package cmd

import (
	"fmt"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/felipejfc/k8s-oidc-auth-builder/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var port int
var clientID string
var clientSecret string
var apiVersion string
var kubeAPI string
var kubeCA string
var environment string

func configureLogger(debug bool) {
	log.SetOutput(os.Stderr)
	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the api",
	Long:  "Starts the api",
	Run: func(cmd *cobra.Command, args []string) {
		loadConfiguration()
		configureLogger(Debug)
		log.WithFields(log.Fields{
			"clientID":     clientID,
			"clientSecret": clientSecret,
			"kubeAPI":      kubeAPI,
			"kubeCA":       kubeCA,
			"port":         port,
		}).Debug("starting app...")
		api := &api.API{
			Port:         port,
			ClientID:     clientID,
			ClientSecret: clientSecret,
			APIVersion:   apiVersion,
			KubeAPI:      kubeAPI,
			KubeCA:       kubeCA,
			Environment:  environment,
		}
		api.Start()
	},
}

func loadConfiguration() {
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("koh")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	clientID = viper.GetString("oidc.clientId")
	clientSecret = viper.GetString("oidc.clientSecret")
	kubeAPI = viper.GetString("kubernetes.api")
	kubeCA = viper.GetString("kubernetes.ca")
	apiVersion = viper.GetString("kubernetes.apiVersion")
}

func init() {
	RootCmd.AddCommand(startCmd)
	startCmd.Flags().IntVarP(&port, "port", "p", 8080, "The port that the API will bind to")
	startCmd.Flags().StringVarP(&environment, "environment", "e", "staging", "The environment that the program will run")
}
