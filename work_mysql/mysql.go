package work_mysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/shn27/Test/utils"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation"
	kmapi "kmodules.xyz/client-go/api/v1"
	api "kubedb.dev/apimachinery/apis/kubedb/v1"
	"kubedb.dev/db-client-go/mysql"
	"strings"
	"xorm.io/xorm"
)

var (
	slowQueryLogPathQuery      = "SHOW VARIABLES LIKE 'slow_query_log_file';"
	innoDBLogBufferSizeQuery   = "SHOW VARIABLES LIKE 'innodb_log_buffer_size';"
	innoDBLogFileSizeQuery     = "SHOW VARIABLES LIKE 'innodb_log_file_size';"
	innoDBLogWaitsQuery        = "SHOW STATUS LIKE 'Innodb_log_waits';"
	innoDBLogWritesQuery       = "SHOW STATUS LIKE 'Innodb_log_writes';"
	Innodb_os_log_fsyncsQuery  = "SHOW STATUS LIKE 'Innodb_os_log_fsyncs';"
	QueriesCausedInnoDBLogWait = "SELECT EVENT_ID, EVENT_NAME, SQL_TEXT, TIMER_WAIT FROM performance_schema.events_statements_history_long WHERE EVENT_NAME LIKE 'wait/io/file/innodb%';"
	tailLines                  = "50"
)

func mysqlQueryTest(client *mysql.XormClient) error {
	result, err := client.Query(innoDBLogBufferSizeQuery)
	if err != nil {
		return err
	}
	fmt.Println("=========================innoDBLogBufferSizeQuery=========================")

	for _, row := range result {
		for _, col := range row {
			fmt.Println(string(col))
		}
	}

	return nil
}

func primaryServiceDNSMySql(db *api.MySQL) string {
	return "localhost" //kc port-forward -n demo svc/mysql 3306
	//return fmt.Sprintf("%v.%v.svc", db.ServiceName(), db.Namespace)
}

func GetMysqlClient() (*mysql.XormClient, error) {
	kbClient, err := utils.GetKBClient()
	if err != nil {
		fmt.Println("failed to get k8s client", err)
		return nil, err
	}
	ref := kmapi.ObjectReference{
		Name:      "mysql",
		Namespace: "demo",
	}
	gvk := schema.GroupVersionKind{
		Version: "v1",
		Group:   "kubedb.com",
		Kind:    "MySQL",
	}

	obj, err := utils.GetK8sObject(gvk, ref, kbClient)
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s object : %v", err)
	}

	db := &api.MySQL{}
	err = runtime.DefaultUnstructuredConverter.
		FromUnstructured(obj.UnstructuredContent(), db)
	if err != nil {
		return nil, fmt.Errorf("failed to convert unstructured object to a concrete type: %w", err)
	}

	mysqlClient, err := mysql.NewKubeDBClientBuilder(kbClient, db).
		WithURL(primaryServiceDNSMySql(db)).
		//WithContext(context.Background()).
		//WithPod("mysql-0").
		GetMySQLXormClient()
	if err != nil {
		fmt.Println("failed to get kube db client: %w", err)
		return nil, err
	}

	err = mysqlClient.Ping()
	if err != nil {
		return nil, err
	}

	return mysqlClient, nil
}

func getMySqlClientUsingCred() (*xorm.Engine, error) {
	user := "root"
	pass := "7pP)IxKy*)qHcA2B"

	dns := getMySQLHostDNS("mysql")
	connectionString := fmt.Sprintf("%v:%v@tcp(%s:%d)/%s?", user, pass, dns, 3306, "mysql")
	engine, err := xorm.NewEngine(api.ResourceSingularMySQL, connectionString)
	if err != nil {
		return nil, err
	}
	engine.SetDefaultContext(context.Background())

	return engine, nil
}

func isPrimary(dsn string) (bool, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return false, err
	}
	defer db.Close()

	var role string
	err = db.QueryRow("SELECT @@global.read_only").Scan(&role)
	if err != nil {
		return false, err
	}

	return role == "0", nil // `read_only = 0` means primary
}

func getMySQLHostDNS(dbname string) string {
	return fmt.Sprintf("%v.%v.%v.svc", "mysql", getMysqlGoverningServiceName(dbname), "demo")
}

func getMysqlGoverningServiceName(dbname string) string {
	return NameWithSuffix(dbname, "pods")
}
func NameWithSuffix(name, suffix string, customLength ...int) string {
	maxLength := validation.DNS1123LabelMaxLength
	if len(customLength) != 0 {
		maxLength = customLength[0]
	}
	if len(suffix) >= maxLength {
		return strings.Trim(suffix[max(0, len(suffix)-maxLength):], "-")
	}
	out := fmt.Sprintf("%s-%s", name[:min(len(name), maxLength-len(suffix)-1)], suffix)
	return strings.Trim(out, "-")
}
