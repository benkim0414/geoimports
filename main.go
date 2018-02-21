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

// Handler is a Lambda function handler to import RET inverters and PV modules to the GEO.
func Handler(ctx context.Context, snsEvent events.SNSEvent) error {
	conf := &geoauth.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		AuthURL:      geoauth.URL,
	}
	client, err := imports.NewClient(ctx, conf)
	if err != nil {
		return err
	}

	types := [...]string{
		imports.TypeRETPVModule,
		imports.TypeRETInverter,
	}
	for _, typ := range types {
		m, err := client.GetRecentImport(typ)
		if err != nil {
			return err
		}
		if m.Status == imports.StatusNew {
			if _, err := client.PutImport(m.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	lambda.Start(Handler)
}
