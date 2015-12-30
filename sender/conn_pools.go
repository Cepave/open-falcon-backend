package sender

import (
	"errors"
	"log"
	"github.com/Cepave/transfer/g"
	cpool "github.com/Cepave/transfer/sender/conn_pool"
	nset "github.com/toolkits/container/set"
	"strings"
)

var (
	errInvalidDSNUnescaped = errors.New("Invalid DSN: Did you forget to escape a param value?")
	errInvalidDSNAddr      = errors.New("Invalid DSN: Network Address not terminated (missing closing brace)")
	errInvalidDSNNoSlash   = errors.New("Invalid DSN: Missing the slash separating the database name")
)

func parseDSN(dsn string) (cfg *cpool.InfluxdbConnection, err error) {
	// New config
	cfg = &cpool.InfluxdbConnection{}

	// [username[:password]@][protocol[(address)]]/dbname
	// Find the last '/' (since the password or the net addr might contain a '/')
	foundSlash := false
	for i := len(dsn) - 1; i >= 0; i-- {
		if dsn[i] == '/' {
			foundSlash = true
			var j, k int

			// left part is empty if i <= 0
			if i > 0 {
				// [username[:password]@][protocol[(address)]]
				// Find the last '@' in dsn[:i]
				for j = i; j >= 0; j-- {
					if dsn[j] == '@' {
						// username[:password]
						// Find the first ':' in dsn[:j]
						for k = 0; k < j; k++ {
							if dsn[k] == ':' {
								cfg.Password = dsn[k+1 : j]
								break
							}
						}
						cfg.Username = dsn[:k]

						break
					}
				}

				// [protocol[(address)]]
				// Find the first '(' in dsn[j+1:i]
				for k = j + 1; k < i; k++ {
					if dsn[k] == '(' {
						// dsn[i-1] must be == ')' if an address is specified
						if dsn[i-1] != ')' {
							if strings.ContainsRune(dsn[k+1:i], ')') {
								return nil, errInvalidDSNUnescaped
							}
							return nil, errInvalidDSNAddr
						}
						cfg.Address = dsn[k+1 : i-1]
						break
					}
				}
				cfg.Protocol = dsn[j+1 : k]
			}

			// /dbname
			cfg.DBName = dsn[i+1 : len(dsn)]

			break
		}
	}

	if !foundSlash && len(dsn) > 0 {
		return nil, errInvalidDSNNoSlash
	}

	// Set default network if empty
	if cfg.Protocol == "" {
		cfg.Protocol = "tcp"
	}

	// Set default address if empty
	if cfg.Address == "" {
		return nil, errors.New("Empty DSN address")
	}

	return
}

func initConnPools() {
	cfg := g.Config()

	judgeInstances := nset.NewStringSet()
	for _, instance := range cfg.Judge.Cluster {
		judgeInstances.Add(instance)
	}
	JudgeConnPools = cpool.CreateSafeRpcConnPools(cfg.Judge.MaxConns, cfg.Judge.MaxIdle,
		cfg.Judge.ConnTimeout, cfg.Judge.CallTimeout, judgeInstances.ToSlice())

	// graph
	graphInstances := nset.NewSafeSet()
	for _, nitem := range cfg.Graph.Cluster2 {
		for _, addr := range nitem.Addrs {
			graphInstances.Add(addr)
		}
	}
	GraphConnPools = cpool.CreateSafeRpcConnPools(cfg.Graph.MaxConns, cfg.Graph.MaxIdle,
		cfg.Graph.ConnTimeout, cfg.Graph.CallTimeout, graphInstances.ToSlice())

	// graph migrating
	if cfg.Graph.Migrating && cfg.Graph.ClusterMigrating != nil {
		graphMigratingInstances := nset.NewSafeSet()
		for _, cnode := range cfg.Graph.ClusterMigrating2 {
			for _, addr := range cnode.Addrs {
				graphMigratingInstances.Add(addr)
			}
		}
		GraphMigratingConnPools = cpool.CreateSafeRpcConnPools(cfg.Graph.MaxConns, cfg.Graph.MaxIdle,
			cfg.Graph.ConnTimeout, cfg.Graph.CallTimeout, graphMigratingInstances.ToSlice())
	}

	influxdbInstances := make([]cpool.InfluxdbConnection, 1)
	dsn, err := parseDSN(cfg.Influxdb.Address)
	if err != nil {
		log.Print("syntax of influxdb address is wrong")
	} else {
		influxdbInstances[0] = *dsn
		InfluxdbConnPools = cpool.CreateInfluxdbConnPools(cfg.Influxdb.MaxConns, cfg.Influxdb.MaxIdle,
			cfg.Influxdb.ConnTimeout, cfg.Influxdb.CallTimeout, influxdbInstances)
	}
}

func DestroyConnPools() {
	JudgeConnPools.Destroy()
	GraphConnPools.Destroy()
	GraphMigratingConnPools.Destroy()
	InfluxdbConnPools.Destroy()
}
