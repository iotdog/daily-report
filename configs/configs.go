package configs

type Configs struct {
	MongoDBAddr  string // MongoDB address
	DBName       string // Database name
	WorkerC      string // name of worker collection
	DailyReportC string // name of daily report collection
	ListenPort   int    // listen port of the server
}

var (
	configs *Configs
)

// Instance returns the global configuration.
func Instance() *Configs {
	if configs == nil {
		configs = new(Configs)
		configs.ListenPort = 1024
		configs.MongoDBAddr = "mongodb://192.168.1.236:27000"
		configs.DBName = "main"
		configs.WorkerC = "workers"
		configs.DailyReportC = "daily_reports"
	}
	return configs
}
