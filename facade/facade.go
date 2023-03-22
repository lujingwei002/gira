package facade

import (
	"context"

	"github.com/lujingwei/gira"
)

func Context() context.Context {
	return gira.Facade().Context()
}

func GetResourceDbClient() gira.MongoClient {
	return gira.Facade().GetResourceDbClient()
}
