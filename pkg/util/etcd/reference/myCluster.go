package reference

type MyCluster struct {
	ClusterName        string
	IP                 string
	PORT               string
	OpenEOEMasterIP    string
	isEtcdBackupServer bool
}
