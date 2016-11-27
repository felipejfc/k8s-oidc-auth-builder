/* Code in this file was grabbed from https://github.com/micahhausler/k8s-oidc-helper

The MIT License (MIT)

Copyright (c) 2016 Micah Hausler

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/

package oidc

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// TokenResponse is the result of GetTokens method
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
}

// KubectlUser is a struct to help build kubeconfig file
type KubectlUser struct {
	Name         string        `yaml:"name"`
	KubeUserInfo *KubeUserInfo `yaml:"user"`
}

// KubeUserInfo is a struct to help build kubeconfig file
type KubeUserInfo struct {
	AuthProvider *AuthProvider `yaml:"auth-provider"`
}

// AuthProvider is a struct to help build kubeconfig file
type AuthProvider struct {
	APConfig *APConfig `yaml:"config"`
	Name     string    `yaml:"name"`
}

// APConfig is a struct to help build kubeconfig file
type APConfig struct {
	ClientID     string `yaml:"client-id"`
	ClientSecret string `yaml:"client-secret"`
	IDToken      string `yaml:"id-token"`
	IdpIssuerURL string `yaml:"idp-issuer-url"`
	RefreshToken string `yaml:"refresh-token"`
}

// KubectlCluster is a struct to help build kubeconfig file
type KubectlCluster struct {
	Cluster *KubectlClusterData `yaml:"cluster"`
	Name    string              `yaml:"name"`
}

// KubectlClusterData is a struct to help build kubeconfig file
type KubectlClusterData struct {
	CertificateAuthorityData string `yaml:"certificate-authority-data"`
	Server                   string `yaml:"server"`
}

// KubectlContext is a struct to help buil kubeconfig file
type KubectlContext struct {
	Context *KubectlContextData `yaml:"context"`
	Name    string              `yaml:"name"`
}

// KubectlContextData is a struct to help buil kubeconfig file
type KubectlContextData struct {
	Cluster string `yaml:"cluster"`
	User    string `yaml:"user"`
}

// KubectlConfig is a struct to help build kubeconfig file
type KubectlConfig struct {
	APIVersion     string            `yaml:"apiVersion"`
	Clusters       []*KubectlCluster `yaml:"clusters"`
	Contexts       []*KubectlContext `yaml:"contexts"`
	CurrentContext string            `yaml:"current-context"`
	Kind           string            `yaml:"kind"`
	Users          []*KubectlUser    `yaml:"users"`
}

// UserInfo  is the result of GetEmail
type UserInfo struct {
	Email string `json:"email"`
}

// GetTokens get the id_token and refresh_token from google
func GetTokens(clientID, clientSecret, code string) (*TokenResponse, error) {
	val := url.Values{}
	val.Add("grant_type", "authorization_code")
	val.Add("redirect_uri", "urn:ietf:wg:oauth:2.0:oob")
	val.Add("client_id", clientID)
	val.Add("client_secret", clientSecret)
	val.Add("code", code)

	resp, err := http.PostForm("https://www.googleapis.com/oauth2/v3/token", val)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	tr := &TokenResponse{}
	err = json.NewDecoder(resp.Body).Decode(tr)
	if err != nil {
		return nil, err
	}
	return tr, nil
}

// GenerateKubectlConfig generates a kubeconfig with cluster, context and user configured
func GenerateKubectlConfig(apiVersion string, kubeAPI string, kubeCA string, environment string, user *KubectlUser) *KubectlConfig {
	cluster := &KubectlCluster{
		Name: environment,
		Cluster: &KubectlClusterData{
			Server: kubeAPI,
			CertificateAuthorityData: kubeCA,
		},
	}
	context := &KubectlContext{
		Name: cluster.Name,
		Context: &KubectlContextData{
			Cluster: cluster.Name,
			User:    user.Name,
		},
	}
	return &KubectlConfig{
		APIVersion:     apiVersion,
		Clusters:       []*KubectlCluster{cluster},
		Contexts:       []*KubectlContext{context},
		CurrentContext: cluster.Name,
		Kind:           "Config",
		Users:          []*KubectlUser{user},
	}
}

// GenerateUser generate a kubeconfig user
func GenerateUser(email, clientID, clientSecret, idToken, refreshToken string) *KubectlUser {
	return &KubectlUser{
		Name: email,
		KubeUserInfo: &KubeUserInfo{
			AuthProvider: &AuthProvider{
				APConfig: &APConfig{
					ClientID:     clientID,
					ClientSecret: clientSecret,
					IDToken:      idToken,
					IdpIssuerURL: "https://accounts.google.com",
					RefreshToken: refreshToken,
				},
				Name: "oidc",
			},
		},
	}
}

// GetEmail gets the user email linked with accessToken
func GetEmail(accessToken string) (string, error) {
	uri, _ := url.Parse("https://www.googleapis.com/oauth2/v1/userinfo")
	q := uri.Query()
	q.Set("alt", "json")
	q.Set("access_token", accessToken)
	uri.RawQuery = q.Encode()
	resp, err := http.Get(uri.String())
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	ui := &UserInfo{}
	err = json.NewDecoder(resp.Body).Decode(ui)
	if err != nil {
		return "", err
	}
	return ui.Email, nil
}
