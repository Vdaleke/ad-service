package tests

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"testing"
)

import (
	"context"
	"net"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"homework10/internal/adapters/repo"
	"homework10/internal/app"
	grpcPort "homework10/internal/ports/grpc"
)

func TestGRPCCreateAd(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(repo.New(), repo.New()))
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	user, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg"})
	assert.NoError(t, err, "client.GetUser")

	assert.Equal(t, "Oleg", user.Name)

	ad, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{
		Title:  "hello",
		Text:   "world",
		UserId: user.Id,
	})
	assert.NoError(t, err)
	assert.Zero(t, ad.Id)
	assert.Equal(t, ad.Title, "hello")
	assert.Equal(t, ad.Text, "world")
	assert.Equal(t, ad.AuthorId, user.Id)
	assert.False(t, ad.Published)
	assert.False(t, ad.CreatedAt.AsTime().IsZero())
	assert.True(t, ad.UpdatedAt.AsTime().IsZero())
}

func TestGRPCChangeAdStatus(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(repo.New(), repo.New()))
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	user, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg"})
	assert.NoError(t, err, "client.GetUser")

	assert.Equal(t, "Oleg", user.Name)

	ad, _ := client.CreateAd(ctx, &grpcPort.CreateAdRequest{
		Title:  "hello",
		Text:   "world",
		UserId: user.Id,
	})
	response, err := client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{UserId: user.Id, AdId: ad.Id, Published: true})
	assert.NoError(t, err)
	assert.True(t, response.Published)

	response, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{UserId: user.Id, AdId: ad.Id, Published: false})
	assert.NoError(t, err)
	assert.False(t, response.Published)
}

func TestGRPCUpdateAd(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(repo.New(), repo.New()))
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	user, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg"})
	assert.NoError(t, err, "client.GetUser")

	assert.Equal(t, "Oleg", user.Name)

	ad, _ := client.CreateAd(ctx, &grpcPort.CreateAdRequest{
		Title:  "hello",
		Text:   "world",
		UserId: user.Id,
	})

	response, err := client.UpdateAd(ctx, &grpcPort.UpdateAdRequest{
		UserId: user.Id, AdId: ad.Id, Title: "привет", Text: "мир"})
	assert.NoError(t, err)
	assert.Equal(t, response.Title, "привет")
	assert.Equal(t, response.Text, "мир")
	assert.False(t, response.CreatedAt.AsTime().IsZero())
	assert.False(t, response.UpdatedAt.AsTime().IsZero())
}

func TestGRPCGetAd(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(repo.New(), repo.New()))
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	user, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg"})
	assert.NoError(t, err, "client.GetUser")

	assert.Equal(t, "Oleg", user.Name)

	ad, _ := client.CreateAd(ctx, &grpcPort.CreateAdRequest{
		Title:  "hello",
		Text:   "world",
		UserId: user.Id,
	})

	response, err := client.GetAd(ctx, &grpcPort.GetAdRequest{Id: ad.Id})
	assert.NoError(t, err)
	assert.Zero(t, response.Id)
	assert.Equal(t, response.Title, "hello")
	assert.Equal(t, response.Text, "world")
	assert.Equal(t, response.AuthorId, user.Id)
	assert.False(t, response.Published)
	assert.False(t, response.CreatedAt.AsTime().IsZero())
	assert.True(t, response.UpdatedAt.AsTime().IsZero())
}

func TestGRPCDeleteAd(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(repo.New(), repo.New()))
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	user, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg"})
	assert.NoError(t, err, "client.GetUser")

	assert.Equal(t, "Oleg", user.Name)

	ad, _ := client.CreateAd(ctx, &grpcPort.CreateAdRequest{
		Title:  "hello",
		Text:   "world",
		UserId: user.Id,
	})

	_, err = client.DeleteAd(ctx, &grpcPort.DeleteAdRequest{AdId: ad.Id, AuthorId: user.Id})
	assert.NoError(t, err)

	_, err = client.GetAd(ctx, &grpcPort.GetAdRequest{Id: ad.Id})
	assert.Equal(t, status.Code(err), codes.InvalidArgument)
}

func TestGRPCListAds(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(repo.New(), repo.New()))
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	user, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg"})
	assert.NoError(t, err, "client.GetUser")

	assert.Equal(t, "Oleg", user.Name)

	ad1, _ := client.CreateAd(ctx, &grpcPort.CreateAdRequest{
		Title:  "hello",
		Text:   "world",
		UserId: user.Id,
	})

	response, err := client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{UserId: user.Id, AdId: ad1.Id, Published: true})
	assert.NoError(t, err)
	assert.True(t, response.Published)

	_, _ = client.CreateAd(ctx, &grpcPort.CreateAdRequest{
		Title:  "best cat",
		Text:   "not for sale",
		UserId: user.Id,
	})

	ads, err := client.ListAds(ctx, &grpcPort.ListAdsRequest{UserId: user.Id, Published: true})
	assert.NoError(t, err)
	assert.Len(t, ads.List, 1)
	assert.Equal(t, ads.List[0].Id, ad1.Id)
	assert.Equal(t, ads.List[0].Title, ad1.Title)
	assert.Equal(t, ads.List[0].Text, ad1.Text)
	assert.Equal(t, ads.List[0].AuthorId, ad1.AuthorId)
	assert.True(t, ads.List[0].Published)
}

func TestGRPCSearchAds(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(repo.New(), repo.New()))
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	user, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg"})
	assert.NoError(t, err, "client.GetUser")

	assert.Equal(t, "Oleg", user.Name)

	ad1, _ := client.CreateAd(ctx, &grpcPort.CreateAdRequest{
		Title:  "hello",
		Text:   "world",
		UserId: user.Id,
	})

	response, err := client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{UserId: user.Id, AdId: ad1.Id, Published: true})
	assert.NoError(t, err)
	assert.True(t, response.Published)

	_, _ = client.CreateAd(ctx, &grpcPort.CreateAdRequest{
		Title:  "best cat",
		Text:   "not for sale",
		UserId: user.Id,
	})

	ads, err := client.SearchAds(ctx, &grpcPort.SearchAdsRequest{Pattern: "ell"})
	assert.NoError(t, err)
	assert.Len(t, ads.List, 1)
	assert.Equal(t, ads.List[0].Id, ad1.Id)
	assert.Equal(t, ads.List[0].Title, ad1.Title)
	assert.Equal(t, ads.List[0].Text, ad1.Text)
	assert.Equal(t, ads.List[0].AuthorId, ad1.AuthorId)
	assert.True(t, ads.List[0].Published)
}

func TestGRPCCreateUser(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(repo.New(), repo.New()))
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	res, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg"})
	assert.NoError(t, err, "client.GetUser")

	assert.Equal(t, "Oleg", res.Name)
}

func TestGRPCUpdateUser(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(repo.New(), repo.New()))
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	user, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg"})
	assert.NoError(t, err, "client.GetUser")

	response, err := client.UpdateUser(ctx, &grpcPort.UpdateUserRequest{Id: user.Id, Name: "Test User 2", Email: "test2@testing.ru"})
	assert.NoError(t, err)
	assert.Zero(t, response.Id)
	assert.Equal(t, response.Name, "Test User 2")
	assert.Equal(t, response.Email, "test2@testing.ru")
}

func TestGRPCGetUser(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(repo.New(), repo.New()))
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	user, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg"})
	assert.NoError(t, err, "client.GetUser")

	response, err := client.GetUser(ctx, &grpcPort.GetUserRequest{Id: user.Id})
	assert.NoError(t, err)
	assert.Equal(t, user.Id, response.Id)
	assert.Equal(t, user.Name, response.Name)
	assert.Equal(t, user.Email, response.Email)
}

func TestGRPCDeleteUser(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(repo.New(), repo.New()))
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	user, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg"})
	assert.NoError(t, err, "client.GetUser")

	_, err = client.DeleteUser(ctx, &grpcPort.DeleteUserRequest{Id: user.Id})
	assert.NoError(t, err)

	_, err = client.GetUser(ctx, &grpcPort.GetUserRequest{Id: user.Id})
	assert.Equal(t, status.Code(err), codes.InvalidArgument)
}

func TestGRPCServerInterceptor(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(
		logging.UnaryServerInterceptor(grpcPort.InterceptorLogger()),
		recovery.UnaryServerInterceptor([]recovery.Option{
			recovery.WithRecoveryHandler(grpcPort.PanicInterceptor),
		}...),
	))
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(repo.New(), repo.New()))
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	user, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg"})
	assert.NoError(t, err, "client.GetUser")

	ad, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{
		Title:  "hello",
		Text:   "world",
		UserId: user.Id,
	})
	assert.NoError(t, err)
	assert.Zero(t, ad.Id)

	response, err := client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{UserId: user.Id, AdId: ad.Id, Published: true})
	assert.NoError(t, err)
	assert.True(t, response.Published)

	_, err = client.UpdateAd(ctx, &grpcPort.UpdateAdRequest{
		UserId: user.Id, AdId: ad.Id, Title: "привет", Text: "мир"})
	assert.NoError(t, err)

	response, err = client.GetAd(ctx, &grpcPort.GetAdRequest{Id: ad.Id})
	assert.NoError(t, err)
	assert.Zero(t, response.Id)

	_, err = client.DeleteUser(ctx, &grpcPort.DeleteUserRequest{Id: user.Id})
	assert.NoError(t, err)

	_, err = client.GetUser(ctx, &grpcPort.GetUserRequest{Id: user.Id})
	assert.Equal(t, status.Code(err), codes.InvalidArgument)
}
