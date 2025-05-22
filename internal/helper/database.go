package helper

import (
	"fmt"

	"github.com/tubagusmf/tbwallet-user-auth/internal/config"
)

func GetConnectionString() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.GetDbHost(),
		config.GetDbUser(),
		config.GetDbPassword(),
		config.GetDbName(),
		config.GetDbPort(),
	)
}
