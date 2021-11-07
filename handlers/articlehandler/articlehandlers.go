package article

import (
	proto "article/articleshare/proto"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var NewArticleConnection proto.ArticleStoreClient

func InitialiseServiceHandlers() {
	connection, connErr := grpc.Dial("localhost:4040", grpc.WithInsecure())
	if connErr != nil {
		panic(connErr)
	}

	NewArticleConnection = proto.NewArticleStoreClient(connection)
}

// ListArticlesHandler  handler function
func ListArticlesHandler(c *gin.Context) {
	ListArticlesBody(c)
}

// ListArticlesBody  handler function
func ListArticlesBody(c *gin.Context) error {

	// extract data from json values
	var listArticleReq proto.ListArticlesBody
	jsonErr := c.BindJSON(&listArticleReq)

	if jsonErr != nil {
		logrus.Error("Error parsing JSON: ", jsonErr)
		return jsonErr
	}

	// Invoke gRPC handler
	res, listError := NewArticleConnection.ListArticles(c, &listArticleReq)
	if listError != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "false",
			"error":  listError.Error(),
		})
		return nil
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "OK",
		"response": res.ArticleLists,
	})

	return nil
}

// CreateArticlesHandler  handler function
func CreateArticlesHandler(c *gin.Context) {
	CreateArticlesBody(c)
}

// CreateArticlesBody  handler function
func CreateArticlesBody(c *gin.Context) error {

	var createArticleReq proto.CreateArticleRequest
	jsonErr := c.BindJSON(&createArticleReq)

	if jsonErr != nil {
		logrus.Error("Error parsing JSON: ", jsonErr)
		return jsonErr
	}

	// Invoke gRPC handler
	_, listError := NewArticleConnection.CreateArticles(c, &createArticleReq)
	if listError != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "false",
			"error":  listError.Error(),
		})
		return nil
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})

	return nil
}
