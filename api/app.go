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

package api

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/felipejfc/k8s-oidc-auth-builder/oidc"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
)

// API the api struct that will keep config
type API struct {
	Port          int
	ClientID      string
	ClientSecret  string
	APIVersion    string
	KubeAPI       string
	KubeCA        string
	Environment   string
	ClusterConfig *oidc.KubectlCluster
	HTTP          *echo.Echo
}

// Start starts the api
func (a *API) Start() {
	a.HTTP = echo.New()
	h := a.HTTP

	h.GET("/oidcurl", a.GetGoogleLoginURL)
	h.POST("/kubeconfig", a.GetKubeConfig(a.APIVersion, a.KubeCA, a.KubeAPI, a.Environment, a.ClientID, a.ClientSecret))

	log.Infof("API listening at port %d", a.Port)
	h.Run(standard.New(fmt.Sprintf(":%d", a.Port)))
}
