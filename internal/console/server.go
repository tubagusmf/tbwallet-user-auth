package console

import (
	"log"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tubagusmf/tbwallet-user-auth/database"
	"github.com/tubagusmf/tbwallet-user-auth/internal/config"
	"github.com/tubagusmf/tbwallet-user-auth/internal/repository"
	"github.com/tubagusmf/tbwallet-user-auth/internal/usecase"

	handlerHttp "github.com/tubagusmf/tbwallet-user-auth/internal/delivery/http"
)

func init() {
	rootCmd.AddCommand(serverCMD)
}

var serverCMD = &cobra.Command{
	Use:   "httpsrv",
	Short: "Start HTTP server",
	Long:  "Start the HTTP server to handle incoming requests for the to-do list application.",
	Run:   httpServer,
}

func httpServer(cmd *cobra.Command, args []string) {
	config.LoadWithViper()

	postgresDB := database.NewPostgres()
	sqlDB, err := postgresDB.DB()
	if err != nil {
		log.Fatalf("Failed to get SQL DB from Gorm: %v", err)
	}
	defer sqlDB.Close()

	userRepo := repository.NewUserRepo(postgresDB)
	userUsecase := usecase.NewUserUsecase(userRepo)
	kycRepo := repository.NewKycRepo(postgresDB)
	kycUsecase := usecase.NewKycUsecase(kycRepo)

	e := echo.New()

	handlerHttp.NewUserHandler(e, userUsecase)
	handlerHttp.NewKycHandler(e, kycUsecase)

	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		errCh <- e.Start(":3000")
	}()

	go func() {
		wg.Wait()
		close(errCh)
	}()

	if err := <-errCh; err != nil {
		if err != http.ErrServerClosed {
			logrus.Errorf("HTTP server error: %v", err)
		}
	}
}
