package aws

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/go-ini/ini"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/spf13/viper"
)

func selectTagVal(Tags []types.Tag, TagName string) string {
	for _, tag := range Tags {
		if *tag.Key == TagName {
			return *tag.Value
		}
	}
	return ""
}

func GetProfile() string {
	var profile string
	sources := []func() string{
		func() string { return viper.GetString("GLOBAL.AWS.AWS_PROFILE") },
		func() string { return os.Getenv("AWS_PROFILE") },
	}

	for _, source := range sources {
		if p := source(); p != "" {
			profile = p
			break
		}
	}

	if profile == "" {
		fmt.Println("AWS_PROFILE not set")
		os.Exit(1)
	}
	return profile
}

func GetRegion() string {
	region := ""
	sources := []func() string{
		func() string { return viper.GetString("GLOBAL.AWS.AWS_PROFILE.REGION") },
		func() string { return os.Getenv("AWS_REGION") },
	}

	for _, source := range sources {
		if p := source(); p != "" {
			region = p
			break
		}
	}
	return region
}

type AWSProfile struct {
	Name       string `mapstructure:"NAME"`
	Region     string `mapstructure:"REGION"`
	PemKeyPath string `mapstructure:"PEM_KEY_PATH"`
	Port       string `mapstructure:"SERVER_PORT"`
	// Endpoint_ID string `mapstructure:"ENDPOINT_ID"`
}

func GetProfileConfig() AWSProfile {
	var awsProfiles []AWSProfile

	// Unmarshal the config into AWSProfile struct
	err := viper.UnmarshalKey("AWS_PROFILES", &awsProfiles)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	// Print out the profiles to verify
	// for _, profile := range awsProfiles {
	// 	fmt.Printf("Profile: %+v\n", profile)
	// }

	var profileConfig AWSProfile
	for _, profile := range awsProfiles {
		if profile.Name == Profile {
			profileConfig = profile
			break
		}
	}

	return profileConfig
}

func GetAWSProfiles() ([]string, error) {
	var profiles []string

	credsFilePath := config.DefaultSharedCredentialsFilename()

	credsFile, err := ini.Load(credsFilePath)
	if err != nil {
		return nil, err
	}
	for _, section := range credsFile.Sections() {
		profiles = append(profiles, section.Name())
	}

	return profiles, nil
}
