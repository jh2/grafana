package metrics

import (
	"runtime"

	"github.com/grafana/grafana/pkg/setting"

	"github.com/prometheus/client_golang/prometheus"
)

const exporterName = "grafana"

var (
	M_Instance_Start       prometheus.Counter
	M_Page_Status          *prometheus.CounterVec
	M_Api_Status           *prometheus.CounterVec
	M_Proxy_Status         *prometheus.CounterVec
	M_Http_Request_Total   *prometheus.CounterVec
	M_Http_Request_Summary *prometheus.SummaryVec

	M_Api_User_SignUpStarted   prometheus.Counter
	M_Api_User_SignUpCompleted prometheus.Counter
	M_Api_User_SignUpInvite    prometheus.Counter
	M_Api_Dashboard_Save       prometheus.Summary
	M_Api_Dashboard_Get        prometheus.Summary
	M_Api_Dashboard_Search     prometheus.Summary
	M_Api_Admin_User_Create    prometheus.Counter
	M_Api_Login_Post           prometheus.Counter
	M_Api_Login_OAuth          prometheus.Counter
	M_Api_Org_Create           prometheus.Counter

	M_Api_Dashboard_Snapshot_Create      prometheus.Counter
	M_Api_Dashboard_Snapshot_External    prometheus.Counter
	M_Api_Dashboard_Snapshot_Get         prometheus.Counter
	M_Api_Dashboard_Insert               prometheus.Counter
	M_Alerting_Result_State              *prometheus.CounterVec
	M_Alerting_Notification_Sent         *prometheus.CounterVec
	M_Aws_CloudWatch_GetMetricStatistics prometheus.Counter
	M_Aws_CloudWatch_ListMetrics         prometheus.Counter
	M_Aws_CloudWatch_GetMetricData       prometheus.Counter
	M_DB_DataSource_QueryById            prometheus.Counter

	// LDAPUsersSyncExecutionTime is a metric for
	// how much time it took to sync the LDAP users
	LDAPUsersSyncExecutionTime prometheus.Summary

	// Timers
	M_DataSource_ProxyReq_Timer prometheus.Summary
	M_Alerting_Execution_Time   prometheus.Summary
)

// StatTotals
var (
	M_Alerting_Active_Alerts prometheus.Gauge
	M_StatTotal_Dashboards   prometheus.Gauge
	M_StatTotal_Users        prometheus.Gauge
	M_StatActive_Users       prometheus.Gauge
	M_StatTotal_Orgs         prometheus.Gauge
	M_StatTotal_Playlists    prometheus.Gauge

	StatsTotalViewers       prometheus.Gauge
	StatsTotalEditors       prometheus.Gauge
	StatsTotalAdmins        prometheus.Gauge
	StatsTotalActiveViewers prometheus.Gauge
	StatsTotalActiveEditors prometheus.Gauge
	StatsTotalActiveAdmins  prometheus.Gauge

	// M_Grafana_Version is a gauge that contains build info about this binary
	//
	// Deprecated: use M_Grafana_Build_Version instead.
	M_Grafana_Version *prometheus.GaugeVec

	// grafanaBuildVersion is a gauge that contains build info about this binary
	grafanaBuildVersion *prometheus.GaugeVec
)

func init() {
	M_Instance_Start = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "instance_start_total",
		Help:      "counter for started instances",
		Namespace: exporterName,
	})

	httpStatusCodes := []string{"200", "404", "500", "unknown"}
	M_Page_Status = newCounterVecStartingAtZero(
		prometheus.CounterOpts{
			Name:      "page_response_status_total",
			Help:      "page http response status",
			Namespace: exporterName,
		}, []string{"code"}, httpStatusCodes...)

	M_Api_Status = newCounterVecStartingAtZero(
		prometheus.CounterOpts{
			Name:      "api_response_status_total",
			Help:      "api http response status",
			Namespace: exporterName,
		}, []string{"code"}, httpStatusCodes...)

	M_Proxy_Status = newCounterVecStartingAtZero(
		prometheus.CounterOpts{
			Name:      "proxy_response_status_total",
			Help:      "proxy http response status",
			Namespace: exporterName,
		}, []string{"code"}, httpStatusCodes...)

	M_Http_Request_Total = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_total",
			Help: "http request counter",
		},
		[]string{"handler", "statuscode", "method"},
	)

	M_Http_Request_Summary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_request_duration_milliseconds",
			Help: "http request summary",
		},
		[]string{"handler", "statuscode", "method"},
	)

	M_Api_User_SignUpStarted = newCounterStartingAtZero(prometheus.CounterOpts{
		Name:      "api_user_signup_started_total",
		Help:      "amount of users who started the signup flow",
		Namespace: exporterName,
	})

	M_Api_User_SignUpCompleted = newCounterStartingAtZero(prometheus.CounterOpts{
		Name:      "api_user_signup_completed_total",
		Help:      "amount of users who completed the signup flow",
		Namespace: exporterName,
	})

	M_Api_User_SignUpInvite = newCounterStartingAtZero(prometheus.CounterOpts{
		Name:      "api_user_signup_invite_total",
		Help:      "amount of users who have been invited",
		Namespace: exporterName,
	})

	M_Api_Dashboard_Save = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:      "api_dashboard_save_milliseconds",
		Help:      "summary for dashboard save duration",
		Namespace: exporterName,
	})

	M_Api_Dashboard_Get = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:      "api_dashboard_get_milliseconds",
		Help:      "summary for dashboard get duration",
		Namespace: exporterName,
	})

	M_Api_Dashboard_Search = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:      "api_dashboard_search_milliseconds",
		Help:      "summary for dashboard search duration",
		Namespace: exporterName,
	})

	M_Api_Admin_User_Create = newCounterStartingAtZero(prometheus.CounterOpts{
		Name:      "api_admin_user_created_total",
		Help:      "api admin user created counter",
		Namespace: exporterName,
	})

	M_Api_Login_Post = newCounterStartingAtZero(prometheus.CounterOpts{
		Name:      "api_login_post_total",
		Help:      "api login post counter",
		Namespace: exporterName,
	})

	M_Api_Login_OAuth = newCounterStartingAtZero(prometheus.CounterOpts{
		Name:      "api_login_oauth_total",
		Help:      "api login oauth counter",
		Namespace: exporterName,
	})

	M_Api_Org_Create = newCounterStartingAtZero(prometheus.CounterOpts{
		Name:      "api_org_create_total",
		Help:      "api org created counter",
		Namespace: exporterName,
	})

	M_Api_Dashboard_Snapshot_Create = newCounterStartingAtZero(prometheus.CounterOpts{
		Name:      "api_dashboard_snapshot_create_total",
		Help:      "dashboard snapshots created",
		Namespace: exporterName,
	})

	M_Api_Dashboard_Snapshot_External = newCounterStartingAtZero(prometheus.CounterOpts{
		Name:      "api_dashboard_snapshot_external_total",
		Help:      "external dashboard snapshots created",
		Namespace: exporterName,
	})

	M_Api_Dashboard_Snapshot_Get = newCounterStartingAtZero(prometheus.CounterOpts{
		Name:      "api_dashboard_snapshot_get_total",
		Help:      "loaded dashboards",
		Namespace: exporterName,
	})

	M_Api_Dashboard_Insert = newCounterStartingAtZero(prometheus.CounterOpts{
		Name:      "api_models_dashboard_insert_total",
		Help:      "dashboards inserted ",
		Namespace: exporterName,
	})

	M_Alerting_Result_State = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "alerting_result_total",
		Help:      "alert execution result counter",
		Namespace: exporterName,
	}, []string{"state"})

	M_Alerting_Notification_Sent = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "alerting_notification_sent_total",
		Help:      "counter for how many alert notifications been sent",
		Namespace: exporterName,
	}, []string{"type"})

	M_Aws_CloudWatch_GetMetricStatistics = newCounterStartingAtZero(prometheus.CounterOpts{
		Name:      "aws_cloudwatch_get_metric_statistics_total",
		Help:      "counter for getting metric statistics from aws",
		Namespace: exporterName,
	})

	M_Aws_CloudWatch_ListMetrics = newCounterStartingAtZero(prometheus.CounterOpts{
		Name:      "aws_cloudwatch_list_metrics_total",
		Help:      "counter for getting list of metrics from aws",
		Namespace: exporterName,
	})

	M_Aws_CloudWatch_GetMetricData = newCounterStartingAtZero(prometheus.CounterOpts{
		Name:      "aws_cloudwatch_get_metric_data_total",
		Help:      "counter for getting metric data time series from aws",
		Namespace: exporterName,
	})

	M_DB_DataSource_QueryById = newCounterStartingAtZero(prometheus.CounterOpts{
		Name:      "db_datasource_query_by_id_total",
		Help:      "counter for getting datasource by id",
		Namespace: exporterName,
	})

	LDAPUsersSyncExecutionTime = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:      "ldap_users_sync_execution_time",
		Help:      "summary for LDAP users sync execution duration",
		Namespace: exporterName,
	})

	M_DataSource_ProxyReq_Timer = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:      "api_dataproxy_request_all_milliseconds",
		Help:      "summary for dataproxy request duration",
		Namespace: exporterName,
	})

	M_Alerting_Execution_Time = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:      "alerting_execution_time_milliseconds",
		Help:      "summary of alert exeuction duration",
		Namespace: exporterName,
	})

	M_Alerting_Active_Alerts = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "alerting_active_alerts",
		Help:      "amount of active alerts",
		Namespace: exporterName,
	})

	M_StatTotal_Dashboards = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "stat_totals_dashboard",
		Help:      "total amount of dashboards",
		Namespace: exporterName,
	})

	M_StatTotal_Users = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "stat_total_users",
		Help:      "total amount of users",
		Namespace: exporterName,
	})

	M_StatActive_Users = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "stat_active_users",
		Help:      "number of active users",
		Namespace: exporterName,
	})

	M_StatTotal_Orgs = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "stat_total_orgs",
		Help:      "total amount of orgs",
		Namespace: exporterName,
	})

	M_StatTotal_Playlists = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "stat_total_playlists",
		Help:      "total amount of playlists",
		Namespace: exporterName,
	})

	StatsTotalViewers = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "stat_totals_viewers",
		Help:      "total amount of viewers",
		Namespace: exporterName,
	})

	StatsTotalEditors = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "stat_totals_editors",
		Help:      "total amount of editors",
		Namespace: exporterName,
	})

	StatsTotalAdmins = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "stat_totals_admins",
		Help:      "total amount of admins",
		Namespace: exporterName,
	})

	StatsTotalActiveViewers = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "stat_totals_active_viewers",
		Help:      "total amount of viewers",
		Namespace: exporterName,
	})

	StatsTotalActiveEditors = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "stat_totals_active_editors",
		Help:      "total amount of active editors",
		Namespace: exporterName,
	})

	StatsTotalActiveAdmins = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "stat_totals_active_admins",
		Help:      "total amount of active admins",
		Namespace: exporterName,
	})

	M_Grafana_Version = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "info",
		Help:      "Information about the Grafana. This metric is deprecated. please use `grafana_build_info`",
		Namespace: exporterName,
	}, []string{"version"})

	grafanaBuildVersion = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "build_info",
		Help:      "A metric with a constant '1' value labeled by version, revision, branch, and goversion from which Grafana was built.",
		Namespace: exporterName,
	}, []string{"version", "revision", "branch", "goversion", "edition"})
}

// SetBuildInformation sets the build information for this binary
func SetBuildInformation(version, revision, branch string) {
	// We export this info twice for backwards compatibility.
	// Once this have been released for some time we should be able to remote `M_Grafana_Version`
	// The reason we added a new one is that its common practice in the prometheus community
	// to name this metric `*_build_info` so its easy to do aggregation on all programs.
	edition := "oss"
	if setting.IsEnterprise {
		edition = "enterprise"
	}

	M_Grafana_Version.WithLabelValues(version).Set(1)
	grafanaBuildVersion.WithLabelValues(version, revision, branch, runtime.Version(), edition).Set(1)
}

func initMetricVars() {
	prometheus.MustRegister(
		M_Instance_Start,
		M_Page_Status,
		M_Api_Status,
		M_Proxy_Status,
		M_Http_Request_Total,
		M_Http_Request_Summary,
		M_Api_User_SignUpStarted,
		M_Api_User_SignUpCompleted,
		M_Api_User_SignUpInvite,
		M_Api_Dashboard_Save,
		M_Api_Dashboard_Get,
		M_Api_Dashboard_Search,
		M_DataSource_ProxyReq_Timer,
		M_Alerting_Execution_Time,
		M_Api_Admin_User_Create,
		M_Api_Login_Post,
		M_Api_Login_OAuth,
		M_Api_Org_Create,
		M_Api_Dashboard_Snapshot_Create,
		M_Api_Dashboard_Snapshot_External,
		M_Api_Dashboard_Snapshot_Get,
		M_Api_Dashboard_Insert,
		M_Alerting_Result_State,
		M_Alerting_Notification_Sent,
		M_Aws_CloudWatch_GetMetricStatistics,
		M_Aws_CloudWatch_ListMetrics,
		M_Aws_CloudWatch_GetMetricData,
		M_DB_DataSource_QueryById,
		LDAPUsersSyncExecutionTime,
		M_Alerting_Active_Alerts,
		M_StatTotal_Dashboards,
		M_StatTotal_Users,
		M_StatActive_Users,
		M_StatTotal_Orgs,
		M_StatTotal_Playlists,
		M_Grafana_Version,
		StatsTotalViewers,
		StatsTotalEditors,
		StatsTotalAdmins,
		StatsTotalActiveViewers,
		StatsTotalActiveEditors,
		StatsTotalActiveAdmins,
		grafanaBuildVersion,
	)

}

func newCounterVecStartingAtZero(opts prometheus.CounterOpts, labels []string, labelValues ...string) *prometheus.CounterVec {
	counter := prometheus.NewCounterVec(opts, labels)

	for _, label := range labelValues {
		counter.WithLabelValues(label).Add(0)
	}

	return counter
}

func newCounterStartingAtZero(opts prometheus.CounterOpts, labelValues ...string) prometheus.Counter {
	counter := prometheus.NewCounter(opts)
	counter.Add(0)

	return counter
}
