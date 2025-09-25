package intializers

import (
	"fmt"
	"main/models"
)

// SyncDataBase runs AutoMigrate for your models and returns an error if anything fails.
func SyncDataBase() error {
	if DB == nil {
		return fmt.Errorf("DB is nil — connect to DB before calling SyncDataBase")
	}

	if err := DB.AutoMigrate(&models.User{}, &models.Note{}, &models.RefreshToken{}); err != nil {
		return fmt.Errorf("AutoMigrate failed: %w", err)
	}

	return nil
}
