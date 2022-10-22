package config

type ServerConfig struct {
	Port string `mapstructure:"PORT"`
}

// type PostgresConn struct {
// 	Username string
// 	Password string
// 	Endpoint string
// 	Port     string
// 	DBName   string
// }

type config struct {
	ServerConfig ServerConfig
	// Postgres     PostgresConn
}
