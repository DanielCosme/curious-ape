package main

import (
	"encoding/json"
	"fmt"
	"github.com/danielcosme/curious-ape/internal/client"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

var authCmd = &cobra.Command{
	Use:   "auth -username <username> -password <password>",
	Short: "Enter credentials to authenticate requests",
	Long:  "Authenticates the client with basic auth credentials to communicate with the API",
	Run:   login,
}

func login(cmd *cobra.Command, args []string) {
	username, _ := cmd.Flags().GetString("username")
	password, _ := cmd.Flags().GetString("password")
	configPath, _ := cmd.Flags().GetString("config-dir-path")
	configFilePath := fmt.Sprintf("%s/credentials.json", configPath)

	credentials := client.Credentials{
		Username: username,
		Password: password,
	}

	client.DefaultService.Auth(credentials.Username, credentials.Password)
	err := client.Ping()
	CheckErr(err)
	SaveCredentials(configFilePath, credentials)
	fmt.Println("Successfully logged-in")
}

func SaveCredentials(fileName string, c client.Credentials) {
	body, err := json.Marshal(c)
	CheckErr(err)
	err = os.WriteFile(fileName, body, os.ModePerm)
	CheckErr(err)
}

func ReadCredentials(fileName string) client.Credentials {
	body, err := os.ReadFile(fileName)
	CheckErr(err)

	var credentials client.Credentials
	err = json.Unmarshal(body, &credentials)
	CheckErr(err)
	return credentials
}

func loadCredentials(cmd *cobra.Command, args []string) error {
	host, _ := cmd.Flags().GetString("host")
	if strings.Contains(host, "localhost") {
		// localhost environment has no password protection
		return nil
	}

	configPath, _ := cmd.Flags().GetString("config-dir-path")
	configFilePath := fmt.Sprintf("%s/credentials.json", configPath)
	credentials := ReadCredentials(configFilePath)
	client.DefaultService.Auth(credentials.Username, credentials.Password)
	client.DefaultService.Host(host)

	return nil
}

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
