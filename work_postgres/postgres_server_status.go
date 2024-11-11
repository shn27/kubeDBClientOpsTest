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
