# AWS Cognito Manager with Go

This sample project demonstrates how to manage AWS Cognito resources with AWS's SDK for Go.

The application is a simple REST API. See all its definitions in `application.go`.

Specific AWS interactions are kept in the `aws.go` file.
To keep things simple, all these functions are part of the `application` struct, where dependencies are kept.
In a real scenario one would eventually want to have this as a separate package.

You'll need an AWS account key/secret with some Cognito policies that allow creating/removing pools and other resources.
For the sake of keeping things simple, I created a user with the policy `arn:aws:iam::aws:policy/AmazonCognitoPowerUser`.
Remember, this may not be the policy to choose in a real scenario.
Make sure you follow the [principle of least privilege](https://docs.aws.amazon.com/IAM/latest/UserGuide/best-practices.html#grant-least-privilege).


Configuration requires some environment variables available at runtime.
The following are read by AWS SDK when we run `LoadDefaultConfig`. More details [here](https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/).
The key 

```
AWS_ACCESS_KEY_ID=xxxxx
AWS_SECRET_ACCESS_KEY=zzzzzzzzzz
AWS_ACCOUNT_ID=nnnnnnnn
```