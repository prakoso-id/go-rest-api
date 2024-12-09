package seeders

import (
	"gorm.io/gorm"
)

func RunSeeders(db *gorm.DB) error {
	// Add all seeders here
	if err := SeedAdminUser(db); err != nil {
		return err
	}

	return nil
}
