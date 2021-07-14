package migrations

import (
	"github.com/go-gormigrate/gormigrate"
	"gorm.io/gorm"
	"log"
)

func Run(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, getMigrations())

	if err := m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
		return err
	}
	log.Printf("Migration did run successfully")
	return nil
}

func getMigrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID: "202104072127",
			Migrate: func(tx *gorm.DB) error {
				type User struct {
					gorm.Model
					Email           string `gorm:"primary_key"`
					PasswordHash    []byte
					PermissionLevel int
				}

				return tx.AutoMigrate(&User{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("users")
			},
		},
		{
			ID: "202104072129",
			Migrate: func(tx *gorm.DB) error {
				type Team struct {
					gorm.Model
					Name    string
					Country string
					Budget  int
					UserID  uint
				}

				return tx.AutoMigrate(&Team{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("teams")
			},
		},
		{
			ID: "202104072130",
			Migrate: func(tx *gorm.DB) error {
				type Player struct {
					gorm.Model
					FirstName   string
					LastName    string
					Country     string
					Age         int
					MarketValue int32
					Position    int
					TeamID      uint
				}

				return tx.AutoMigrate(&Player{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("players")
			},
		},
		{
			ID: "202104072131",
			Migrate: func(tx *gorm.DB) error {
				type Transfer struct {
					gorm.Model
					PlayerID uint
					Ask      int
				}

				return tx.AutoMigrate(&Transfer{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("transfers")
			},
		},
	}
}
