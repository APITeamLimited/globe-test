package worker

import (
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)

func init() {
	funcframework.RegisterHTTPFunction("worker", RunGoogleCloud)
}
