package config

type Config struct {
	Port      int      `toml:"port"`
	Routines  int      `toml:"routines"`
	Db        string   `toml:"db"` // database name
	Dsn       string   `toml:"dsn"`
	Url       string   `toml:"url"`
	CacheSize int      `toml:"cachesize"`
	Addrs     []string `toml:"addrs"`
	LogConfig *Log     `toml:"log"`
}

type Log struct {
	Level  string `toml:"level"`
	Prefix string `toml:"prefix"`
}
