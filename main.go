package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/benkim0414/geoauth"
	"github.com/benkim0414/geoimports/imports"
)

const (
	clientID     = "ben.kim@greenenergytrading.com.au"
	clientSecret = "AQICAHjEkFlMfByWgCktJWRFfuVMhkFaCtDynoodWngDmNQ14gFceiCwAnMkyiCpxY54/YXvAAAAajBoBgkqhkiG9w0BBwagWzBZAgEAMFQGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMEylGqoZ1CXyzNEE6AgEQgCfnvXYZM7iDXNGqaEsUIn7dN3daDY1pZS4imkUZlvNie2eK9q1orC4="
)

func Handler(ctx context.Context, snsEvent events.SNSEvent) error {
	conf := &geoauth.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		AuthURL:      geoauth.GEOAuthURL,
	}
	client, err := imports.NewClient(ctx, conf)
	if err != nil {
		return err
	}

	const pull = true
	types, err := client.GetImportTypes(pull)
	if err != nil {
		return err
	}
	for _, t := range types {
		if t.Available() {
			if _, err := client.PutImport(t.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	lambda.Start(Handler)
}
