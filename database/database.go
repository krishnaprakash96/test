package database

import (
	"article/articleshare/constants"
	"article/articleshare/proto"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// BaseDAO implements the common database interatction functions
type BaseDAO struct {
	articleDB *sql.DB
}

var baseDAOObj BaseDAO

// Connect connects to the database.
func Connect() {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		constants.User,
		constants.Password,
		constants.Host,
		constants.Port,
		constants.Database,
	)

	baseDAOObj = BaseDAO{}
	logrus.Info("Connecting to article database")
	var dbError error
	baseDAOObj.articleDB, dbError = sql.Open("mysql", connectionString)
	if dbError != nil {
		log.Fatalln("Connection to virtual reservation database failed - " + dbError.Error())
	}
}

// ExecQuery executes the query against the database.
func ExecQuery(query string, args ...interface{}) error {
	executeStatement, err := baseDAOObj.articleDB.Prepare(query)
	if err != nil {
		logrus.Error("Error preparing create statement: ", query, " -> ", err)
		return status.Error(codes.Internal, "Error preparing create statement")
	}

	defer executeStatement.Close()
	_, err = executeStatement.Exec(args...)
	if err != nil {
		logrus.Error("Error executing query: ", query, " -> ", err)
		return status.Error(codes.Internal, "Error executing create query")
	}
	return nil
}

// CreateArticle database function
func CreateArticle(title, description, name, email string) error {
	insertQuery := "INSERT INTO Articles (Title, ArticleDescription, Author, Email, CreatedTime) VALUES(?, ?, ?, ?, ?, ?, ?)"

	startTime := time.Now().UTC()
	executeError := ExecQuery(insertQuery, title, description, name, email, startTime)
	if executeError != nil {
		logrus.Error("Error inserting article to database")
		return executeError
	}
	return nil
}

// FetchArticleList database function
func FetchArticleList(title, description string) ([]*proto.ListArticlesBody, error) {
	selectQuery := ""
	if len(title) == 0 && len(description) == 0 {
		selectQuery = "SELECT Title, ArticleDescription FROM Articles WHERE Status=Approved"
		return getArticleLists(selectQuery, "")
	} else if len(title) == 0 {
		selectQuery = "SELECT Title, ArticleDescription FROM Articles WHERE Title=? AND Status=Approved"
		return getArticleLists(selectQuery, description)
	} else if len(description) == 0 {
		selectQuery = "SELECT Title, ArticleDescription FROM Articles WHERE Description=? AND Status=Approved"
		return getArticleLists(selectQuery, title)
	}
	return nil, nil
}

func getArticleLists(selectQuery, param string) ([]*proto.ListArticlesBody, error) {
	var rows *sql.Rows
	var queryError error
	if param == "" {
		rows, queryError = baseDAOObj.articleDB.Query(selectQuery)
	} else {
		rows, queryError = baseDAOObj.articleDB.Query(selectQuery, param)
	}
	if queryError != nil {
		return nil, queryError
	}

	defer rows.Close()

	lists := make([]*proto.ListArticlesBody, 0)
	for rows.Next() {
		var articleList proto.ListArticlesBody
		rows.Scan(&articleList.Title, &articleList.Description)
		lists = append(lists, &articleList)
	}
	return lists, nil
}
