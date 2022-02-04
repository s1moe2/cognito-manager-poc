# AWS Cognito Manager with Go

This sample project demonstrates how to manage AWS Cognito resources with AWS's SDK for Go.

The application is a simple REST API. See all its definitions in `application.go`.

Specific AWS interactions are kept in the `aws.go` file.
To keep things simple, all these functions are part of the `application` struct, where dependencies are kept.
In a real scenario one would eventually want to have this as a separate package.