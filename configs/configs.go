package configs

type Configs struct {
	MongoDBAddr     string // MongoDB address
	DBName          string // Database name
	WorkerC         string // name of worker collection
	DailyReportC    string // name of daily report collection
	ListenPort      int    // listen port of the server
	MailBoxUserName string // user name of the mailbox used to send daily report
	MailBoxPwd      string // password of the mailbox
	MailBoxSMTP     string // SMTP server of the mailbox
	MailBoxPort     int    // SMTP server port
	MailBoxSSL      bool
}

var (
	configs *Configs
)

// Instance returns the global configuration.
func Instance() *Configs {
	if configs == nil {
		configs = new(Configs)
		configs.ListenPort = 1024
		// configs.MongoDBAddr = "mongodb://192.168.1.236:27000"
		configs.MongoDBAddr = "mongodb://127.0.0.1:27017"
		configs.DBName = "main"
		configs.WorkerC = "workers"
		configs.DailyReportC = "daily_reports"
		configs.MailBoxUserName = ""
		configs.MailBoxPwd = ""
		configs.MailBoxSMTP = "smtp.exmail.qq.com"
		configs.MailBoxPort = 465
		configs.MailBoxSSL = true
	}
	return configs
}
