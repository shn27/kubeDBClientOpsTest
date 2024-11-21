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

func TestSharedBuffers() {
	pgClient, err := GetPostgresClient()
	if err != nil {
		klog.Error(err, "failed to get postgres client")
		return
	}

	var sharedBuffers string
	if err = pgClient.DB.QueryRow("SHOW shared_buffers").Scan(&sharedBuffers); err != nil {
		klog.Error(err, "failed to get shared buffers")
		return
	}

	klog.Infof("Shared buffers: %s\n", sharedBuffers)
}

func TestGetMaxAllowedMemory() {
	pgClient, err := GetPostgresClient()
	if err != nil {
		klog.Error(err, "failed to get postgres client")
		return
	}
}
