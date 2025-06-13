package console

import (
	"log"
	"net"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/tubagusmf/tbwallet-user-auth/database"
	"github.com/tubagusmf/tbwallet-user-auth/internal/config"
	grpcHandler "github.com/tubagusmf/tbwallet-user-auth/internal/delivery/grpc"
	handlerHttp "github.com/tubagusmf/tbwallet-user-auth/internal/delivery/http"
	"github.com/tubagusmf/tbwallet-user-auth/internal/repository"
	"github.com/tubagusmf/tbwallet-user-auth/internal/usecase"

	kycPb "github.com/tubagusmf/tbwallet-user-auth/pb/kycdoc"
	userPb "github.com/tubagusmf/tbwallet-user-auth/pb/user"
)

var startServeCmd = &cobra.Command{
	Use:   "httpsrv",
	Short: "Start HTTP and gRPC servers",
	Run: func(cmd *cobra.Command, args []string) {
		config.LoadWithViper()

		dbConn := database.NewPostgres()
		sqlDB, err := dbConn.DB()
		if err != nil {
			log.Fatalf("Failed to get SQL DB from Gorm: %v", err)
		}
		defer sqlDB.Close()

		userRepo := repository.NewUserRepo(dbConn)
		kycRepo := repository.NewKycRepo(dbConn)

		userUsecase := usecase.NewUserUsecase(userRepo, nil)
		kycUsecase := usecase.NewKycUsecase(kycRepo, nil)

		quitCh := make(chan bool, 1)

		go func() {
			e := echo.New()
			e.GET("/ping", func(c echo.Context) error {
				return c.String(http.StatusOK, "pong")
			})

			handlerHttp.NewUserHandler(e, userUsecase)
			handlerHttp.NewKycHandler(e, kycUsecase)

			log.Println("HTTP server running on :3000")
			if err := e.Start(":3000"); err != nil && err != http.ErrServerClosed {
				logrus.Fatalf("Failed to start HTTP server: %v", err)
			}
		}()

		go func() {
			grpcServer := grpc.NewServer()
			userHandler := grpcHandler.NewUsergRPCHandler(userUsecase)
			kycHandler := grpcHandler.NewKycdocGRPCHandler(kycUsecase)

			userPb.RegisterUserServiceServer(grpcServer, userHandler)
			kycPb.RegisterKycdocServiceServer(grpcServer, kycHandler)

			lis, err := net.Listen("tcp", ":4001")
			if err != nil {
				log.Fatalf("Failed to listen on gRPC port: %v", err)
			}

			log.Println("gRPC server running on :4001")
			if err := grpcServer.Serve(lis); err != nil {
				log.Fatalf("Failed to serve gRPC: %v", err)
			}
		}()

		<-quitCh
	},
}

func init() {
	rootCmd.AddCommand(startServeCmd)
}
