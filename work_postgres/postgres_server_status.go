package work_postgres

import (
	"encoding/json"

	gohumanize "github.com/dustin/go-humanize"
	utils "github.com/shn27/Test/utils"
	"k8s.io/klog/v2"
)

func TestPostgresServerStatus() {
	kubeClient, err := utils.GetKBClient()
	if err != nil {
		klog.Error(err, "failed to get kube client")
		return
	}

	db, err := GetPostgresDB(kubeClient)
	if err != nil {
		klog.Error(err, "failed to get postgres db")
		return
	}

	pgClient, err := GetPostgresClient(kubeClient, db)
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
	kubeClient, err := utils.GetKBClient()
	if err != nil {
		klog.Error(err, "failed to get kube client")
		return
	}

	db, err := GetPostgresDB(kubeClient)
	if err != nil {
		klog.Error(err, "failed to get postgres db")
		return
	}

	pgClient, err := GetPostgresClient(kubeClient, db)
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
	kubeClient, err := utils.GetKBClient()
	if err != nil {
		klog.Error(err, "failed to get kube client")
		return
	}

	db, err := GetPostgresDB(kubeClient)
	if err != nil {
		klog.Error(err, "failed to get postgres db")
		return
	}

	pgClient, err := GetPostgresClient(kubeClient, db)
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
	kubeClient, err := utils.GetKBClient()
	if err != nil {
		klog.Error(err, "failed to get kube client")
		return
	}

	db, err := GetPostgresDB(kubeClient)
	if err != nil {
		klog.Error(err, "failed to get postgres db")
		return
	}

	pgClient, err := GetPostgresClient(kubeClient, db)
	if err != nil {
		klog.Error(err, "failed to get postgres client")
		return
	}

	totalMemory, err := GetTotalMemory(pgClient, db)
	if err != nil {
		klog.Error(err, "failed to get total memory")
		return
	}

	klog.Infof("Total memory: %s\n", gohumanize.IBytes(uint64(totalMemory)))
}
