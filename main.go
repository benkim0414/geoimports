package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/benkim0414/geoauth"
	"github.com/benkim0414/geoimports/imports"
)

const (
	clientID     = "ben.kim@greenenergytrading.com.au"
	clientSecret = "AQICAHhVrOaNnK50PWMuS5OKJu776vi3df3OYAVa5vStP5DRGAEv8Q8TClCGDlk/ovIQ6tlcAAAAbDBqBgkqhkiG9w0BBwagXTBbAgEAMFYGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMG/oyqN3cLwXuv5q6AgEQgCmbQnTiB8FSW+YVyh2TsxctymAxBvAuXRdmP9mMIuF5rCpbhEjKJa3v9Q=="
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

	_, err = client.GetImportTypes(true)
	if err != nil {
		log.Fatal(err)
		return err
	}

	types := [...]string{
		imports.TypeRETSWHOver700L,
		imports.TypeRETSWH700LOrLess,
		imports.TypeRETPVModule,
		imports.TypeRETInverter,
		imports.TypeRETASHP,
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
