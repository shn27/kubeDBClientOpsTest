package work_postgres

import (
	"encoding/json"

	"k8s.io/klog/v2"
)

func TestPostgresServerStatus() {
	pgClient, err := GetPostgresClient()
	if err != nil {
		klog.Error(err, "failed to get postgres client")
		return
	}

	stats := pgClient.DB.Stats()
	klog.Info("=======Test postgres server stats=======")

	prettyData, err := json.MarshalIndent(stats, "  ", "   ")
	if err != nil {
		klog.Error(err, "failed to marshal db stats")
	}

	klog.Info(string(prettyData))
}

func TestClientFuncs() {
	pgClient, err := GetPostgresClient()
	if err != nil {
		klog.Error(err, "failed to get postgres client")
		return
	}

	err = pgClient.DB.Ping()
	if err != nil {
		klog.Error(err, "failed to ping postgres")
		return
	}
	klog.Info("Pinged postgres\n")
	klog.Infof("pgClient.DB.Stats().InUse : %d", pgClient.DB.Stats().InUse)
}

/*
{
     "MaxOpenConnections": 0,
     "OpenConnections": 1,
     "InUse": 0,
     "Idle": 1,
     "WaitCount": 0,
     "WaitDuration": 0,
     "MaxIdleClosed": 0,
     "MaxIdleTimeClosed": 0,
     "MaxLifetimeClosed": 0
  }
*/
