package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string `env:"POSTGRES_HOST" env-default:"db"`
	Port     string `env:"POSTGRES_PORT" env-default:"5432"`
	User     string `env:"POSTGRES_USER" env-default:"postgres"`
	Password string `env:"POSTGRES_PASSWORD" env-default:"postgres"`
	DBName   string `env:"POSTGRES_DB" env-default:"pr_reviewer"`
	DSN      string
}

func (c *Config) FormatConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}

type Repo struct {
	DB  *sql.DB
	DSN string

	UserRepository  *UserRepository
	TeamRepository  *TeamRepository
	PrRepository    *PrRepository
	StatsRepository *StatsRepository
}

func NewDB(cfg Config) (*Repo, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * 60)

	return &Repo{
		DB:             db,
		DSN:            cfg.DSN,
		UserRepository: NewUserRepository(db),
		TeamRepository: NewTeamRepository(db),
		PrRepository:   NewPrRepository(db),
        StatsRepository: NewStatsRepository(db),}, nil
}
