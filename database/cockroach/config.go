package cockroach

var (
	host        string = "localhost"
	port        int    = 26257
	username    string = ""
	password    string = ""
	database    string = "defaultdb"
	sslMode     string = "disable"
	options     string
	sslRootCert string
)

func Host() string        { return host }
func Port() int           { return port }
func Username() string    { return username }
func Password() string    { return password }
func Database() string    { return database }
func Options() string     { return options }
func SSLMode() string     { return sslMode }
func SSLRootCert() string { return sslRootCert }

func SetHost(newHost string)               { host = newHost }
func SetPort(newPort int)                  { port = newPort }
func SetUsername(newUsername string)       { username = newUsername }
func SetPassword(newPassword string)       { password = newPassword }
func SetDatabase(newDatabase string)       { database = newDatabase }
func SetOptions(newOptions string)         { options = newOptions }
func SetSSLMode(newSSLMode string)         { sslMode = newSSLMode }
func SetSSLRootCert(newSSLRootCert string) { sslRootCert = newSSLRootCert }
