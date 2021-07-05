package migrations

import (
	"gorm.io/gorm"
	"../models"
	"github.com/go-gormigrate/gormigrate/v2"
)

func Run(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{}, &models.Player{}, &models.Team{})
	/*
	m := gormigrate.New(db, gormigrate.DefaultOptions, getMigrations())

	if err := m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
		return err
	}
	log.Printf("Migration did run successfully")
	return nil
	*/
}

func getMigrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID: "202104072127",
			Migrate: func(tx *gorm.DB) error {
				type User struct {
					gorm.Model
					Email           string
					PasswordHash    []byte
					PermissionLevel int
				}
				return tx.AutoMigrate(&User{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("users")
			},
		},
	}
}