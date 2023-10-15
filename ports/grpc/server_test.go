package grpc

import (
	"context"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"homework10/internal/ads"
	"homework10/internal/app"
	"homework10/internal/mocks"
	"homework10/internal/users"
	"net"
	"testing"
	"time"
)

func BenchmarkAdService_CreateUser(b *testing.B) {
	lis := bufconn.Listen(1024 * 1024)
	b.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	b.Cleanup(func() {
		srv.Stop()
	})

	mockedApp := &mocks.App{}
	mockedApp.On("CreateUser", mock.Anything, mock.Anything).
		Return(users.User{}, nil)

	svc := NewService(mockedApp)
	RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(b, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	b.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(b, err, "grpc.DialContext")

	b.Cleanup(func() {
		conn.Close()
	})

	client := NewAdServiceClient(conn)

	for i := 0; i < b.N; i++ {
		_, _ = client.CreateUser(ctx, &CreateUserRequest{Name: "Oleg"})
	}
}

func TestAdService_DefunctUser(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	mockedApp := &mocks.App{}

	mockedApp.On("CreateUser", mock.Anything).
		Return(users.User{}, nil)

	svc := NewService(mockedApp)
	RegisterAdServiceServer(srv, svc)

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

	client := NewAdServiceClient(conn)

	mockedApp.On("CreateAd", mock.Anything, mock.Anything, mock.Anything).
		Return(ads.Ad{}, app.DefunctUser)
	mockedApp.On("UpdateAd", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(ads.Ad{}, app.DefunctUser)
	mockedApp.On("ChangeAdStatus", mock.Anything, mock.Anything, mock.Anything).
		Return(ads.Ad{}, app.DefunctUser)
	mockedApp.On("DeleteAd", mock.Anything, mock.Anything).
		Return(app.DefunctUser)
	mockedApp.On("UpdateUser", mock.Anything, mock.Anything, mock.Anything).
		Return(users.User{}, app.DefunctUser)
	mockedApp.On("GetUser", mock.Anything).
		Return(users.User{}, app.DefunctUser)
	mockedApp.On("DeleteUser", mock.Anything).
		Return(app.DefunctUser)

	_, err = client.CreateAd(ctx, &CreateAdRequest{
		Title:  "best cat",
		Text:   "not for sale",
		UserId: 12321,
	})

	assert.ErrorIs(t, err, status.New(codes.InvalidArgument, "invalid information received").Err())

	_, err = client.UpdateAd(ctx, &UpdateAdRequest{
		Title:  "best cat",
		Text:   "not for sale",
		UserId: 12321,
		AdId:   0,
	})

	assert.ErrorIs(t, err, status.New(codes.InvalidArgument, "invalid information received").Err())

	_, err = client.ChangeAdStatus(ctx, &ChangeAdStatusRequest{
		Published: false,
		UserId:    0,
	})

	assert.ErrorIs(t, err, status.New(codes.InvalidArgument, "invalid information received").Err())

	_, err = client.DeleteAd(ctx, &DeleteAdRequest{
		AdId:     0,
		AuthorId: 0,
	})

	assert.ErrorIs(t, err, status.New(codes.InvalidArgument, "invalid information received").Err())

	_, err = client.UpdateUser(ctx, &UpdateUserRequest{
		Id: 0,
	})

	assert.ErrorIs(t, err, status.New(codes.InvalidArgument, "invalid information received").Err())

	_, err = client.GetUser(ctx, &GetUserRequest{
		Id: 0,
	})

	assert.ErrorIs(t, err, status.New(codes.InvalidArgument, "invalid information received").Err())

	_, err = client.DeleteUser(ctx, &DeleteUserRequest{
		Id: 0,
	})

	assert.ErrorIs(t, err, status.New(codes.InvalidArgument, "invalid information received").Err())

}

func TestAdService_UnknownError(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	mockedApp := &mocks.App{}

	mockedApp.On("CreateUser", mock.Anything).
		Return(users.User{}, nil)

	svc := NewService(mockedApp)
	RegisterAdServiceServer(srv, svc)

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

	client := NewAdServiceClient(conn)

	mockedApp.On("CreateAd", mock.Anything, mock.Anything, mock.Anything).
		Return(ads.Ad{}, errors.New("Unknown error"))
	mockedApp.On("GetAd", mock.Anything).
		Return(ads.Ad{}, errors.New("Unknown error"))
	mockedApp.On("UpdateAd", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(ads.Ad{}, errors.New("Unknown error"))
	mockedApp.On("ChangeAdStatus", mock.Anything, mock.Anything, mock.Anything).
		Return(ads.Ad{}, errors.New("Unknown error"))
	mockedApp.On("DeleteAd", mock.Anything, mock.Anything).
		Return(errors.New("Unknown error"))
	mockedApp.On("CreateUser", mock.Anything, mock.Anything).
		Return(users.User{}, errors.New("Unknown error"))
	mockedApp.On("UpdateUser", mock.Anything, mock.Anything, mock.Anything).
		Return(users.User{}, errors.New("Unknown error"))
	mockedApp.On("GetUser", mock.Anything).
		Return(users.User{}, errors.New("Unknown error"))
	mockedApp.On("DeleteUser", mock.Anything).
		Return(errors.New("Unknown error"))

	_, err = client.CreateAd(ctx, &CreateAdRequest{
		Title:  "best cat",
		Text:   "not for sale",
		UserId: 12321,
	})

	assert.ErrorIs(t, err, status.New(codes.Unknown, "an unknown error has occurred").Err())

	_, err = client.GetAd(ctx, &GetAdRequest{
		Id: 0,
	})

	assert.ErrorIs(t, err, status.New(codes.Unknown, "an unknown error has occurred").Err())

	_, err = client.UpdateAd(ctx, &UpdateAdRequest{
		Title:  "best cat",
		Text:   "not for sale",
		UserId: 12321,
		AdId:   0,
	})

	assert.ErrorIs(t, err, status.New(codes.Unknown, "an unknown error has occurred").Err())

	_, err = client.ChangeAdStatus(ctx, &ChangeAdStatusRequest{
		Published: false,
		UserId:    0,
	})

	assert.ErrorIs(t, err, status.New(codes.Unknown, "an unknown error has occurred").Err())

	_, err = client.DeleteAd(ctx, &DeleteAdRequest{
		AdId:     0,
		AuthorId: 0,
	})

	assert.ErrorIs(t, err, status.New(codes.Unknown, "an unknown error has occurred").Err())

	_, err = client.CreateUser(ctx, &CreateUserRequest{})

	assert.ErrorIs(t, err, status.New(codes.Unknown, "an unknown error has occurred").Err())

	_, err = client.UpdateUser(ctx, &UpdateUserRequest{
		Id: 0,
	})

	assert.ErrorIs(t, err, status.New(codes.Unknown, "an unknown error has occurred").Err())

	_, err = client.GetUser(ctx, &GetUserRequest{
		Id: 0,
	})

	assert.ErrorIs(t, err, status.New(codes.Unknown, "an unknown error has occurred").Err())

	_, err = client.DeleteUser(ctx, &DeleteUserRequest{
		Id: 0,
	})

	assert.ErrorIs(t, err, status.New(codes.Unknown, "an unknown error has occurred").Err())
}
