package database

type Postgres struct {
	DB *sql.DB
}

func New(ctx context.Context, cfg *config.Config) (*Postgres, error) {
	
}