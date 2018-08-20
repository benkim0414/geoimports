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
	clientSecret = "AQICAHhVrOaNnK50PWMuS5OKJu776vi3df3OYAVa5vStP5DRGAGkwt9bQN1V4INyDl9upR8SAAAAajBoBgkqhkiG9w0BBwagWzBZAgEAMFQGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMKc48x5mwgj6vkha8AgEQgCcrA4uNwmKAOqnNpNtPVL2qqF6+BaMqYARW6nhx2W0Rh/MG+xc3WuY="
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
