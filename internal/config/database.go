package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"minyjae/go-starter/internal/adapters/presistance/models"
)

func SetupDatabase(config *Config) *gorm.DB {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort, config.DBSSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")

	if shouldRunMigration() {
		runMigration(db)
	} else {
		autoMigrate := os.Getenv("AUTO_MIGRATE")
		appEnv := os.Getenv("APP_ENV")

		if autoMigrate == "false" {
			log.Printf("Skipping database migration (AUTO_MIGRATE=false)")
		} else if appEnv == "production" && autoMigrate != "true" {
			log.Printf("Skipping database migration (production environment, set AUTO_MIGRATE=true to enable)")
		} else {
			log.Printf("Skipping database migration (set AUTO_MIGRATE=true to enable)")
		}
	}

	return db
}

func shouldRunMigration() bool {
	if os.Getenv("AUTO_MIGRATE") == "false" {
		return false
	}

	if os.Getenv("AUTO_MIGRATE") == "true" {
		return true
	}

	if os.Getenv("APP_ENV") == "development" {
		return true
	}

	return false
}

// runMigration: เพิ่ม model ใหม่ใน list ของ AutoMigrate
// ตัวอย่าง: db.AutoMigrate(&models.User{}, &models.List{})
func runMigration(db *gorm.DB) {
	log.Println("Starting database migration...")

	err := db.AutoMigrate(
		&models.User{},
		&models.LineUser{},
		&models.MessageLog{},
		&models.AssistantIntent{},
		&models.Todo{},
		&models.Expense{},
		&models.CalendarEvent{},
		&models.Reminder{},
		&models.Note{},
		&models.ConversationSession{},
		&models.IntegrationAccount{},
		&models.EmbeddingDocument{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migrated successfully")
}

// RunMigrationManual: เรียก migration ด้วยตัวเองโดยไม่ผ่านกลไก auto ของ SetupDatabase
// ใช้ผ่าน CLI / job แยกได้ ถ้าอยาก opt-out จาก migration ตอน boot server
func RunMigrationManual(config *Config) error {
	db := SetupDatabase(config)

	log.Println("Running manual migration...")

	err := db.AutoMigrate(
		&models.User{},
		&models.LineUser{},
		&models.MessageLog{},
		&models.AssistantIntent{},
		&models.Todo{},
		&models.Expense{},
		&models.CalendarEvent{},
		&models.Reminder{},
		&models.Note{},
		&models.ConversationSession{},
		&models.IntegrationAccount{},
		&models.EmbeddingDocument{},
	)
	if err != nil {
		return fmt.Errorf("migration failed: %v", err)
	}

	log.Println("Manual migration completed successfully")
	return nil
}
