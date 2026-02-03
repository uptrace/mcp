package bootstrap

import (
	"context"
	"net/http"

	"github.com/uptrace/mcp/appconf"
	"github.com/uptrace/mcp/uptraceapi"
	"github.com/yorunikakeru4/oapi-codegen-dd/v3/pkg/runtime"
)

// httpClient wraps http.Client to implement runtime.HttpRequestDoer.
type httpClient struct {
	*http.Client
}

func (c *httpClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	return c.Client.Do(req.WithContext(ctx))
}

// NewUptraceClient creates a new Uptrace API client from config.
func NewUptraceClient(conf *appconf.Config) (*uptraceapi.Client, error) {
	return uptraceapi.NewDefaultClient(
		conf.Uptrace.APIURL,
		runtime.WithHTTPClient(&httpClient{http.DefaultClient}),
		runtime.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+conf.Uptrace.APIToken)
			return nil
		}),
	)
}
