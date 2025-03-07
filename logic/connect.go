package logic

import (
	"context"
	"es4gophers/domain"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/keploy/go-sdk/integrations/khttpclient"
)

func ConnectWithElasticsearch(ctx context.Context) context.Context {
	interceptor := khttpclient.NewInterceptor(http.DefaultTransport)
	newClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
			
		},
		Transport: interceptor,
	})
	if err != nil {
		panic(err)
	}

	return context.WithValue(ctx, domain.ClientKey, newClient)

}
