package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	cipTypes "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

func (app *application) createCognitoUserPool(ctx context.Context, name string) (*cip.CreateUserPoolOutput, error) {
	return app.cognitoClient.CreateUserPool(ctx, &cip.CreateUserPoolInput{
		PoolName: aws.String(name),
		AccountRecoverySetting: &cipTypes.AccountRecoverySettingType{
			RecoveryMechanisms: []cipTypes.RecoveryOptionType{{
				Name:     cipTypes.RecoveryOptionNameTypeVerifiedEmail,
				Priority: 1,
			}},
		},
		AdminCreateUserConfig: &cipTypes.AdminCreateUserConfigType{
			AllowAdminCreateUserOnly: true,
		},
		AutoVerifiedAttributes: []cipTypes.VerifiedAttributeType{
			cipTypes.VerifiedAttributeTypeEmail,
		},
		MfaConfiguration: "",
		Policies: &cipTypes.UserPoolPolicyType{
			PasswordPolicy: &cipTypes.PasswordPolicyType{
				MinimumLength:                 12,
				RequireLowercase:              true,
				RequireNumbers:                true,
				RequireSymbols:                true,
				RequireUppercase:              true,
				TemporaryPasswordValidityDays: 1,
			},
		},
		UsernameAttributes: []cipTypes.UsernameAttributeType{
			cipTypes.UsernameAttributeTypeEmail,
		},
		UsernameConfiguration: &cipTypes.UsernameConfigurationType{
			CaseSensitive: aws.Bool(false),
		},
	})
}

func (app *application) createCognitoUserPoolClient(ctx context.Context, name string, poolID string) (*cip.CreateUserPoolClientOutput, error) {
	return app.cognitoClient.CreateUserPoolClient(ctx, &cip.CreateUserPoolClientInput{
		ClientName:     aws.String(name),
		UserPoolId:     aws.String(poolID),
		GenerateSecret: false,
	})
}

func (app *application) deleteCognitoUserPool(ctx context.Context, poolID string) error {
	_, err := app.cognitoClient.DeleteUserPool(ctx, &cip.DeleteUserPoolInput{
		UserPoolId: aws.String(poolID),
	})
	return err
}
