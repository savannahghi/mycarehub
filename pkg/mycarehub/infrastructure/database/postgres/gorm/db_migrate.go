package gorm

import "log"

// Migrate updates tables
func (db *PGInstance) Migrate() {
	for _, t := range allTables() {
		if !db.DB.Migrator().HasTable(t) {
			if err := db.DB.Migrator().CreateTable(t); err != nil {
				log.Fatalf("error occurred while creating table %v: %v", t, err)
			}
		}

		if err := db.DB.AutoMigrate(t); err != nil {
			log.Fatalf("error occurred while performing migration: %v", err)
		}
	}
}
