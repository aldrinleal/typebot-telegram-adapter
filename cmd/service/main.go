package main

import (
	"github.com/aldrinleal/typebot-telegram-adapter/util"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	xff "github.com/sebest/xff"
	"github.com/shurcooL/go-goon"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"net/http"
)

var (
	ginLambda *ginadapter.GinLambda
	engine    *gin.Engine
)

func init() {
	// TODO Investigate better integration logrus / lambda / gin
	log.Infof("Initializing")

	engine = BuildEngine()

	ginLambda = ginadapter.New(engine)
}

func main() {
	if util.IsRunningOnLambda() {
		lambda.Start(func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
			log.Infof("req: %s", goon.Sdump(req))

			return ginLambda.ProxyWithContext(ctx, req)
		})
	} else {
		endpoint := ":" + util.EnvIf("PORT", "8000")

		xffmw, _ := xff.Default()

		log.Fatalf("Oops: %s", http.ListenAndServe(endpoint, xffmw.Handler(engine)))
	}
}
