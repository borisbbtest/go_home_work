package postgres

const pluginName = "Postgres"

type Plugin struct {
	connMgr *connManager
}

type requestHandler func(conn *postgresConn, key string, params []string) (res interface{}, err error)

// impl is the pointer to the plugin implementation.
var impl Plugin

func (p *Plugin) Stop() {
	p.connMgr.stop()
	p.connMgr = nil
}

// whereToConnect builds a session based on key's parameters and a configuration file.
func whereToConnect(params []string) (u *URI, err error) {
	var uri string
	user := ""
	if len(params) > 1 {
		user = params[1]
	}

	password := ""
	if len(params) > 2 {
		password = params[2]
	}

	database := ""
	if len(params) > 3 && len(params) < 5 {
		database = params[3]
	}

	// The first param can be either a URI or a session identifier
	if len(params) > 0 && len(params[0]) > 0 {
		if isLooksLikeURI(params[0]) {
			// Use the URI defined as key's parameter
			uri = params[0]
		}
	}

	if len(user) > 0 || len(password) > 0 || len(database) > 0 {
		return newURIWithCreds(uri, user, password, database)
	}

	return parseURI(uri)
}

// Export implements the Exporter interface.
func (p *Plugin) NewDBConn(key string, params []string) (result interface{}, err error) {
	var (
		handler       requestHandler
		handlerParams []string
	)

	u, err := whereToConnect(params)
	if err != nil {
		return nil, err
	}
	// get connection string for PostgreSQL
	connString := u.URI()
	switch key {
	case keyPostgresConnections:
		handler = p.connectionsHandler // postgres.connections[[connString]]
	case keyPostgresPing:
		handler = p.pingHandler // postgres.ping[[connString]]
		return nil, errorUnsupportedQuery
	}

	conn, err := p.connMgr.GetPostgresConnection(connString)
	if err != nil {
		// Here is another logic of processing connection errors if postgres.ping is requested
		if key == keyPostgresPing {
			return postgresPingFailed, nil
		}
		log.Errorf("connection error: %s", err)
		log.Debugf("parameters: %+v", params)
		return nil, err
	}

	return handler(conn, key, handlerParams)
}
