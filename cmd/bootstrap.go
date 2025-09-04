package cmd

import (
	"fmt"
	"time"

	"github.com/Temisaputra/warOnk/internal/infrastructure/config"
	repository "github.com/Temisaputra/warOnk/internal/infrastructure/db"
	"github.com/Temisaputra/warOnk/internal/infrastructure/logger"
	"github.com/Temisaputra/warOnk/pkg/auth"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Dependencies struct {
	DB         *gorm.DB
	Logger     *zap.Logger
	Cfg        *config.Config
	JwtService auth.JwtService
}

func InitDependencies() *Dependencies {
	cfg := config.Get()

	// Init logger
	log := logger.NewLogger(cfg)

	// Init DB
	dsn := fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBUsername, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		log.Fatal("failed to connect database", zap.Error(err))
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Minute * time.Duration(cfg.DBMaxLifetime))
	sqlDB.SetConnMaxIdleTime(time.Minute * time.Duration(cfg.DBMaxIdleTime))

	userRepo := repository.NewUserRepo(db)

	jwtSvc := auth.NewJwtService(*cfg, *log, userRepo)

	return &Dependencies{
		DB:         db,
		Logger:     log,
		Cfg:        cfg,
		JwtService: jwtSvc,
	}
}
