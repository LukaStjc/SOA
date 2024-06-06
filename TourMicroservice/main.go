package main

import (
	"go-tourm/controllers"
	"go-tourm/initializers"
	"go-tourm/middleware"
	configurations "go-tourm/startup"

	"github.com/gin-gonic/gin"

	"log"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

//const serviceName = "tour-service"

func init() {
	configuration := configurations.NewConfigurations()
	// initializers.LoadEnvVariables()
	initializers.ConnectToDb(configuration)
	initializers.SyncDatabase()
	initializers.PreloadTours()

}

func main() {

	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	var err error
	tp, err = initTracer()
	if err != nil {
		log.Fatalf("Failed to initialize tracer: %v", err)
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	r.POST("/create-tour", controllers.CreateTour)
	r.GET("/guide/:id/tours", controllers.GetToursByUser)
	r.Run()

}

func httpErrorInternalServerError(err error, span trace.Span, ctx *gin.Context) {
	httpError(err, span, ctx, http.StatusInternalServerError)
}

func httpError(err error, span trace.Span, ctx *gin.Context, status int) {
	log.Println(err.Error())
	span.RecordError(err) // Record the error in the span
	span.SetStatus(codes.Error, err.Error())
	ctx.String(status, err.Error())
}
