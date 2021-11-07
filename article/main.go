package article

import (
	"article/articleshare/database"
	proto "article/articleshare/proto"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type article struct{}

func (a *article) ListArticles(ctx context.Context, req *proto.ListArticlesBody) (*proto.ListArticlesResponse, error) {
	listResults := make([]*proto.ListArticlesBody, 0)
	var listError error
	listResults, listError = database.FetchArticleList(req.GetTitle(), req.GetDescription())
	if listError != nil {
		log.Fatalf("Error while listing articles %v", listError)
	}
	return &proto.ListArticlesResponse{
		Status:       true,
		ArticleLists: listResults,
	}, nil

}

func (a *article) CreateArticles(ctx context.Context, req *proto.CreateArticleRequest) (*proto.CreateArticleResponse, error) {
	if createError := database.CreateArticle(req.Title, req.Description, req.AuthorDetails.Name,
		req.AuthorDetails.Emailid); createError != nil {
		log.Fatalf("Couldn't create article %v", createError)
	}

	return nil, nil
}

func main() {
	database.Connect()
	listener, listenErr := net.Listen("tcp", "4040")
	if listenErr != nil {
		panic(listenErr)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterArticleStoreServer(grpcServer, &article{})
	reflection.Register(grpcServer)

	if serveError := grpcServer.Serve(listener); serveError != nil {
		panic(serveError)
	}
}
