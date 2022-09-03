package main

import (
	"context"
	"es4gophers/logic"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/keploy/go-sdk/integrations/kgin/v1"
	"github.com/keploy/go-sdk/keploy"
)

var ctx = context.Background()

func putURL(c *gin.Context) {
	var m map[string]string

	err := c.ShouldBindJSON(&m)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to decode req"})
		return
	}
	u := m["indexName"]
	if u == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing indexName param"})
		return
	}
	ctx = logic.LoadMoviesFromFile(ctx)
	logic.IndexMoviesAsDocuments(ctx, u)
	t := time.Now()
	c.JSON(http.StatusOK, gin.H{
		"ts":    t.UnixNano(),
		"index": u + " indexed !",
	})
}

func getURL(c *gin.Context) {
	var m map[string]string
	err := c.ShouldBindJSON(&m)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to decode req"})
		return
	}
	indexName := m["indexName"]
	docId := m["docId"]
	fmt.Println(indexName, " ", docId)
	if indexName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing indexName param"})
		return
	}
	if docId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing docId param"})
		return
	}
	movieName := logic.QueryMovieByDocumentID(ctx, indexName, docId)
	// c.Redirect(http.StatusSeeOther, u.URL)
	t := time.Now()
	c.JSON(http.StatusOK, gin.H{
		"ts":          t.UnixNano(),
		"Movie Name ": movieName,
	})
}

func main() {
	port := "8080"
	// initialize keploy
	k := keploy.New(keploy.Config{
		App: keploy.AppConfig{
			Name: "sample-elastic-search",
			Port: port,
		},
		Server: keploy.ServerConfig{
			URL: "http://localhost:8081/api",
		},
	})

	r := gin.Default()
	// ctx := context.Background()
	// integrate keploy with gin router
	kgin.GinV1(k, r)
	ctx = logic.ConnectWithElasticsearch(ctx)
	r.GET("/param", getURL)
	r.POST("/indexName", putURL)

	// router.POST("/user", controllers.CreateAUser())
	// router.GET("/user/:userId", controllers.GetAUser())
	// router.PUT("/user/:userId", controllers.EditAUser())
	// router.DELETE("/user/:userId", controllers.DeleteAUser())
	// router.GET("/users", controllers.GetAllUsers())

	r.Run(":" + port)

	// ctx = logic.LoadMoviesFromFile(ctx)
	// ctx = logic.ConnectWithElasticsearch(ctx)
	// logic.IndexMoviesAsDocuments(ctx)
	// logic.QueryMovieByDocumentID(ctx)
	// logic.BestKeanuActionMovies(ctx)
	// logic.MovieCountPerGenreAgg(ctx)

}
