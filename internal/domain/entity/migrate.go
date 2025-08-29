package entity

import "gorm.io/gorm"

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Users{},
	)
}

func Drop(db *gorm.DB) error {
	return db.Migrator().DropTable(
		&Users{},
	)
}
