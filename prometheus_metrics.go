//
// Copyright (c) 2015-2024 Hanzo AI, Inc.
//
// This file is part of Hanzo S3 stack
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.
//

package madmin

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/dustin/go-humanize"
	"github.com/prometheus/common/expfmt"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prom2json"
)

// MetricsRespBodyLimit sets the top level limit to the size of the
// metrics results supported by this library.
var (
	MetricsRespBodyLimit = int64(humanize.GiByte)
)

// NodeMetrics - returns Node Metrics in Prometheus format
//
//	The client needs to be configured with the endpoint of the desired node
func (client *MetricsClient) NodeMetrics(ctx context.Context) ([]*prom2json.Family, error) {
	return client.GetMetrics(ctx, "node")
}

// ClusterMetrics - returns Cluster Metrics in Prometheus format
func (client *MetricsClient) ClusterMetrics(ctx context.Context) ([]*prom2json.Family, error) {
	return client.GetMetrics(ctx, "cluster")
}

// BucketMetrics - returns Bucket Metrics in Prometheus format
func (client *MetricsClient) BucketMetrics(ctx context.Context) ([]*prom2json.Family, error) {
	return client.GetMetrics(ctx, "bucket")
}

// ResourceMetrics - returns Resource Metrics in Prometheus format
func (client *MetricsClient) ResourceMetrics(ctx context.Context) ([]*prom2json.Family, error) {
	return client.GetMetrics(ctx, "resource")
}

// GetMetrics - returns Metrics of given subsystem in Prometheus format
func (client *MetricsClient) GetMetrics(ctx context.Context, subSystem string) ([]*prom2json.Family, error) {
	reqData := metricsRequestData{
		relativePath: "/v2/metrics/" + subSystem,
	}

	// Execute GET on /minio/v2/metrics/<subSys>
	resp, err := client.executeGetRequest(ctx, reqData)
	if err != nil {
		return nil, err
	}
	defer closeResponse(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, httpRespToErrorResponse(resp)
	}

	return ParsePrometheusResults(io.LimitReader(resp.Body, MetricsRespBodyLimit))
}

func ParsePrometheusResults(reader io.Reader) (results []*prom2json.Family, err error) {
	// We could do further content-type checks here, but the
	// fallback for now will anyway be the text format
	// version 0.0.4, so just go for it and see if it works.
	parser := expfmt.NewTextParser(model.UTF8Validation)
	metricFamilies, err := parser.TextToMetricFamilies(reader)
	if err != nil {
		return nil, fmt.Errorf("reading text format failed: %v", err)
	}
	results = make([]*prom2json.Family, 0, len(metricFamilies))
	for _, mf := range metricFamilies {
		results = append(results, prom2json.NewFamily(mf))
	}
	return results, nil
}

// Prometheus v2 metrics
var (
	ClusterV2Metrics = []string{
		// Cluster capacity metrics
		"s3_cluster_bucket_total",
		"s3_cluster_capacity_raw_free_bytes",
		"s3_cluster_capacity_raw_total_bytes",
		"s3_cluster_capacity_usable_free_bytes",
		"s3_cluster_capacity_usable_total_bytes",
		"s3_cluster_objects_size_distribution",
		"s3_cluster_objects_version_distribution",
		"s3_cluster_usage_object_total",
		"s3_cluster_usage_total_bytes",
		"s3_cluster_usage_version_total",
		"s3_cluster_usage_deletemarker_total",
		// Cluster drive metrics
		"s3_cluster_drive_offline_total",
		"s3_cluster_drive_online_total",
		"s3_cluster_drive_total",
		// Cluster health metrics
		"s3_cluster_nodes_offline_total",
		"s3_cluster_nodes_online_total",
		"s3_cluster_write_quorum",
		"s3_cluster_health_status",
		"s3_cluster_health_erasure_set_healing_drives",
		"s3_cluster_health_erasure_set_online_drives",
		"s3_cluster_health_erasure_set_read_quorum",
		"s3_cluster_health_erasure_set_write_quorum",
		"s3_cluster_health_erasure_set_status",
		// S3 API requests metrics
		"s3_s3_requests_incoming_total",
		"s3_s3_requests_inflight_total",
		"s3_s3_requests_rejected_auth_total",
		"s3_s3_requests_rejected_header_total",
		"s3_s3_requests_rejected_invalid_total",
		"s3_s3_requests_rejected_timestamp_total",
		"s3_s3_requests_total",
		"s3_s3_requests_waiting_total",
		"s3_s3_requests_ttfb_seconds_distribution",
		"s3_s3_traffic_received_bytes",
		"s3_s3_traffic_sent_bytes",
		// Scanner metrics
		"s3_node_scanner_bucket_scans_finished",
		"s3_node_scanner_bucket_scans_started",
		"s3_node_scanner_directories_scanned",
		"s3_node_scanner_objects_scanned",
		"s3_node_scanner_versions_scanned",
		"s3_node_syscall_read_total",
		"s3_node_syscall_write_total",
		"s3_usage_last_activity_nano_seconds",
		// Inter node metrics
		"s3_inter_node_traffic_dial_avg_time",
		"s3_inter_node_traffic_received_bytes",
		"s3_inter_node_traffic_sent_bytes",
		// Process metrics
		"s3_node_process_cpu_total_seconds",
		"s3_node_process_resident_memory_bytes",
		"s3_node_process_starttime_seconds",
		"s3_node_process_uptime_seconds",
		// File descriptor metrics
		"s3_node_file_descriptor_limit_total",
		"s3_node_file_descriptor_open_total",
		// Node metrics
		"s3_node_go_routine_total",
		"s3_node_io_rchar_bytes",
		"s3_node_io_read_bytes",
		"s3_node_io_wchar_bytes",
		"s3_node_io_write_bytes",
	}
	ReplicationV2Metrics = []string{
		// Cluster replication metrics
		"s3_cluster_replication_last_hour_failed_bytes",
		"s3_cluster_replication_last_hour_failed_count",
		"s3_cluster_replication_last_minute_failed_bytes",
		"s3_cluster_replication_last_minute_failed_count",
		"s3_cluster_replication_total_failed_bytes",
		"s3_cluster_replication_total_failed_count",
		"s3_cluster_replication_received_bytes",
		"s3_cluster_replication_received_count",
		"s3_cluster_replication_sent_bytes",
		"s3_cluster_replication_sent_count",
		"s3_cluster_replication_proxied_get_requests_total",
		"s3_cluster_replication_proxied_head_requests_total",
		"s3_cluster_replication_proxied_delete_tagging_requests_total",
		"s3_cluster_replication_proxied_get_tagging_requests_total",
		"s3_cluster_replication_proxied_put_tagging_requests_total",
		"s3_cluster_replication_proxied_get_requests_failures",
		"s3_cluster_replication_proxied_head_requests_failures",
		"s3_cluster_replication_proxied_delete_tagging_requests_failures",
		"s3_cluster_replication_proxied_get_tagging_requests_failures",
		"s3_cluster_replication_proxied_put_tagging_requests_failures",
		// Node replication metrics
		"s3_node_replication_current_active_workers",
		"s3_node_replication_average_active_workers",
		"s3_node_replication_max_active_workers",
		"s3_node_replication_link_online",
		"s3_node_replication_link_offline_duration_seconds",
		"s3_node_replication_link_downtime_duration_seconds",
		"s3_node_replication_average_link_latency_ms",
		"s3_node_replication_max_link_latency_ms",
		"s3_node_replication_current_link_latency_ms",
		"s3_node_replication_current_transfer_rate",
		"s3_node_replication_average_transfer_rate",
		"s3_node_replication_max_transfer_rate",
		"s3_node_replication_last_minute_queued_count",
		"s3_node_replication_last_minute_queued_bytes",
		"s3_node_replication_average_queued_count",
		"s3_node_replication_average_queued_bytes",
		"s3_node_replication_max_queued_bytes",
		"s3_node_replication_max_queued_count",
		"s3_node_replication_recent_backlog_count",
	}
	BucketV2Metrics = []string{
		// Bucket metrics
		"s3_bucket_objects_size_distribution",
		"s3_bucket_objects_version_distribution",
		"s3_bucket_traffic_received_bytes",
		"s3_bucket_traffic_sent_bytes",
		"s3_bucket_usage_object_total",
		"s3_bucket_usage_version_total",
		"s3_bucket_usage_deletemarker_total",
		"s3_bucket_usage_total_bytes",
		"s3_bucket_requests_inflight_total",
		"s3_bucket_requests_total",
		"s3_bucket_requests_ttfb_seconds_distribution",
		// Bucket replication metrics
		"s3_bucket_replication_last_minute_failed_bytes",
		"s3_bucket_replication_last_minute_failed_count",
		"s3_bucket_replication_last_hour_failed_bytes",
		"s3_bucket_replication_last_hour_failed_count",
		"s3_bucket_replication_total_failed_bytes",
		"s3_bucket_replication_total_failed_count",
		"s3_bucket_replication_latency_ms",
		"s3_bucket_replication_received_bytes",
		"s3_bucket_replication_received_count",
		"s3_bucket_replication_sent_bytes",
		"s3_bucket_replication_sent_count",
		"s3_bucket_replication_proxied_get_requests_total",
		"s3_bucket_replication_proxied_head_requests_total",
		"s3_bucket_replication_proxied_delete_tagging_requests_total",
		"s3_bucket_replication_proxied_get_tagging_requests_total",
		"s3_bucket_replication_proxied_put_tagging_requests_total",
		"s3_bucket_replication_proxied_get_requests_failures",
		"s3_bucket_replication_proxied_head_requests_failures",
		"s3_bucket_replication_proxied_delete_tagging_requests_failures",
		"s3_bucket_replication_proxied_get_tagging_requests_failures",
		"s3_bucket_replication_proxied_put_tagging_requests_failures",
	}
	NodeV2Metrics = []string{
		"s3_node_drive_free_bytes",
		"s3_node_drive_free_inodes",
		"s3_node_drive_latency_us",
		"s3_node_drive_offline_total",
		"s3_node_drive_online_total",
		"s3_node_drive_total",
		"s3_node_drive_total_bytes",
		"s3_node_drive_used_bytes",
		"s3_node_drive_errors_timeout",
		"s3_node_drive_errors_ioerror",
		"s3_node_drive_errors_availability",
		"s3_node_drive_io_waiting",
	}
	ResourceV2Metrics = []string{
		"s3_node_drive_total_bytes",
		"s3_node_drive_used_bytes",
		"s3_node_drive_total_inodes  ",
		"s3_node_drive_used_inodes",
		"s3_node_drive_reads_per_sec",
		"s3_node_drive_reads_kb_per_sec",
		"s3_node_drive_reads_await",
		"s3_node_drive_writes_per_sec",
		"s3_node_drive_writes_kb_per_sec",
		"s3_node_drive_writes_await",
		"s3_node_drive_perc_util",
		"s3_node_if_rx_bytes",
		"s3_node_if_rx_bytes_avg",
		"s3_node_if_rx_bytes_max",
		"s3_node_if_rx_errors",
		"s3_node_if_rx_errors_avg",
		"s3_node_if_rx_errors_max",
		"s3_node_if_tx_bytes",
		"s3_node_if_tx_bytes_avg",
		"s3_node_if_tx_bytes_max",
		"s3_node_if_tx_errors",
		"s3_node_if_tx_errors_avg",
		"s3_node_if_tx_errors_max",
		"s3_node_cpu_avg_user",
		"s3_node_cpu_avg_user_avg",
		"s3_node_cpu_avg_user_max",
		"s3_node_cpu_avg_system",
		"s3_node_cpu_avg_system_avg",
		"s3_node_cpu_avg_system_max",
		"s3_node_cpu_avg_idle",
		"s3_node_cpu_avg_idle_avg",
		"s3_node_cpu_avg_idle_max",
		"s3_node_cpu_avg_iowait",
		"s3_node_cpu_avg_iowait_avg",
		"s3_node_cpu_avg_iowait_max",
		"s3_node_cpu_avg_nice",
		"s3_node_cpu_avg_nice_avg",
		"s3_node_cpu_avg_nice_max",
		"s3_node_cpu_avg_steal",
		"s3_node_cpu_avg_steal_avg",
		"s3_node_cpu_avg_steal_max",
		"s3_node_cpu_avg_load1",
		"s3_node_cpu_avg_load1_avg",
		"s3_node_cpu_avg_load1_max",
		"s3_node_cpu_avg_load1_perc",
		"s3_node_cpu_avg_load1_perc_avg",
		"s3_node_cpu_avg_load1_perc_max",
		"s3_node_cpu_avg_load5",
		"s3_node_cpu_avg_load5_avg",
		"s3_node_cpu_avg_load5_max",
		"s3_node_cpu_avg_load5_perc",
		"s3_node_cpu_avg_load5_perc_avg",
		"s3_node_cpu_avg_load5_perc_max",
		"s3_node_cpu_avg_load15",
		"s3_node_cpu_avg_load15_avg",
		"s3_node_cpu_avg_load15_max",
		"s3_node_cpu_avg_load15_perc",
		"s3_node_cpu_avg_load15_perc_avg",
		"s3_node_cpu_avg_load15_perc_max",
		"s3_node_mem_available",
		"s3_node_mem_available_avg",
		"s3_node_mem_available_max",
		"s3_node_mem_buffers",
		"s3_node_mem_buffers_avg",
		"s3_node_mem_buffers_max",
		"s3_node_mem_cache",
		"s3_node_mem_cache_avg",
		"s3_node_mem_cache_max",
		"s3_node_mem_free",
		"s3_node_mem_free_avg",
		"s3_node_mem_free_max",
		"s3_node_mem_shared",
		"s3_node_mem_shared_avg",
		"s3_node_mem_shared_max",
		"s3_node_mem_total",
		"s3_node_mem_total_avg",
		"s3_node_mem_total_max",
		"s3_node_mem_used",
		"s3_node_mem_used_avg",
		"s3_node_mem_used_max",
		"s3_node_mem_used_perc",
		"s3_node_mem_used_perc_avg",
		"s3_node_mem_used_perc_max",
	}
)
