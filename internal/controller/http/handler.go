package http

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/romik1505/userDetailsService/graph"
	"github.com/romik1505/userDetailsService/graph/resolver"
	v1 "github.com/romik1505/userDetailsService/internal/controller/http/v1"
	"github.com/romik1505/userDetailsService/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	PersonService service.Persons
}

func NewHandler(ps service.Persons) *Handler {
	return &Handler{
		PersonService: ps,
	}
}

func (h *Handler) NewRouter() *gin.Engine {
	r := gin.Default()

	handlerV1 := v1.NewHandler(h.PersonService)
	api := r.Group("/api")
	{
		handlerV1.Init(api)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/query", h.graphqlHandler())
	r.GET("/", playgroundHandler())

	return r
}

// Defining the Graphql handler
func (h *Handler) graphqlHandler() gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	s := handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{Resolvers: resolver.NewResolver(h.PersonService)},
		),
	)

	return func(c *gin.Context) {
		s.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
