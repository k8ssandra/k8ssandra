---
title: "CassandraConfig CRD"
linkTitle: "CassandraConfig CRD"
no_list: true
toc_hide: true
simple_list: false
weight: 6
description: >
  CassandraConfig Custom Resource Definition (CRD) reference for use with K8ssandra Operator.
---

### Custom Resources



* [AuditLogOptions](#auditlogoptions)
* [CassandraConfig](#cassandraconfig)
* [CassandraYaml](#cassandrayaml)
* [FullQueryLoggerOptions](#fullqueryloggeroptions)
* [Group](#group)
* [JvmOptions](#jvmoptions)
* [ParameterizedClass](#parameterizedclass)
* [ReplicaFilteringProtectionOptions](#replicafilteringprotectionoptions)
* [RequestSchedulerOptions](#requestscheduleroptions)
* [SubnetGroups](#subnetgroups)
* [TrackWarnings](#trackwarnings)

#### AuditLogOptions



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| enabled |  | bool | true |
| logger |  | *[ParameterizedClass](#parameterizedclass) | false |
| included_keyspaces |  | *string | false |
| excluded_keyspaces |  | string | false |
| included_categories |  | *string | false |
| excluded_categories |  | string | false |
| included_users |  | *string | false |
| excluded_users |  | *string | false |
| roll_cycle |  | *string | false |
| block |  | *bool | false |
| max_queue_weight |  | *int | false |
| max_log_size |  | *int | false |
| archive_command |  | *string | false |
| max_archive_retries |  | *int | false |

[Back to Custom Resources](#custom-resources)

#### CassandraConfig



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| cassandraYaml |  | [CassandraYaml](#cassandrayaml) | false |
| jvmOptions |  | [JvmOptions](#jvmoptions) | false |

[Back to Custom Resources](#custom-resources)

#### CassandraYaml

CassandraYaml defines the contents of the cassandra.yaml file. For more info see: https://cassandra.apache.org/doc/latest/cassandra/configuration/cass_yaml_file.html

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| allocate_tokens_for_keyspace | Exists in 3.11, 4.0, trunk | *string | false |
| allocate_tokens_for_local_replication_factor | Exists in: 4.0, trunk | *int | false |
| audit_logging_options | Exists in: 4.0, trunk | *[AuditLogOptions](#auditlogoptions) | false |
| auth_read_consistency_level | Exists in trunk | *string | false |
| auth_write_consistency_level | Exists in trunk | *string | false |
| authenticator | Exists in 3.11, 4.0, trunk | *string | false |
| authorizer | Exists in 3.11, 4.0, trunk | *string | false |
| auto_hints_cleanup_enabled | Exists in trunk | *bool | false |
| auto_optimise_full_repair_streams | Exists in: 4.0, trunk | *bool | false |
| auto_optimise_inc_repair_streams | Exists in: 4.0, trunk | *bool | false |
| auto_optimise_preview_repair_streams | Exists in: 4.0, trunk | *bool | false |
| auto_snapshot | Exists in 3.11, 4.0, trunk | *bool | false |
| autocompaction_on_startup_enabled | Exists in: 4.0, trunk | *bool | false |
| automatic_sstable_upgrade | Exists in: 4.0, trunk | *bool | false |
| available_processors | Exists in trunk | *int | false |
| back_pressure_enabled | Exists in 3.11, 4.0, trunk | *bool | false |
| back_pressure_strategy | Exists in 3.11, 4.0, trunk | *[ParameterizedClass](#parameterizedclass) | false |
| batch_size_fail_threshold_in_kb | Exists in 3.11, 4.0, trunk | *int | false |
| batch_size_warn_threshold_in_kb | Exists in 3.11, 4.0, trunk | *int | false |
| batchlog_replay_throttle_in_kb | Exists in 3.11, 4.0, trunk | *int | false |
| block_for_peers_in_remote_dcs | Exists in: 4.0, trunk | *bool | false |
| block_for_peers_timeout_in_secs | Exists in: 4.0, trunk | *int | false |
| buffer_pool_use_heap_if_exhausted | Exists in 3.11, 4.0, trunk | *bool | false |
| cas_contention_timeout_in_ms | Exists in 3.11, 4.0, trunk | *int | false |
| cdc_enabled | Exists in 3.11, 4.0, trunk | *bool | false |
| cdc_free_space_check_interval_ms | Exists in 3.11, 4.0, trunk | *int | false |
| cdc_raw_directory | Exists in 3.11, 4.0, trunk | *string | false |
| cdc_total_space_in_mb | Exists in 3.11, 4.0, trunk | *int | false |
| check_for_duplicate_rows_during_compaction | Exists in 3.11, 4.0, trunk | *bool | false |
| check_for_duplicate_rows_during_reads | Exists in 3.11, 4.0, trunk | *bool | false |
| client_encryption_options | Exists in 3.11, 4.0, trunk | *encryption.ClientEncryptionOptions | false |
| client_error_reporting_exclusions | Exists in trunk | *[SubnetGroups](#subnetgroups) | false |
| column_index_cache_size_in_kb | Exists in 3.11, 4.0, trunk | *int | false |
| column_index_size_in_kb | Exists in 3.11, 4.0, trunk | *int | false |
| commitlog_compression | Exists in 3.11, 4.0, trunk | *[ParameterizedClass](#parameterizedclass) | false |
| commitlog_max_compression_buffers_in_pool | Exists in 3.11, 4.0, trunk | *int | false |
| commitlog_periodic_queue_size | Exists in 3.11, 4.0, trunk | *int | false |
| commitlog_segment_size_in_mb | Exists in 3.11, 4.0, trunk | *int | false |
| commitlog_sync | Exists in 3.11, 4.0, trunk | *string | false |
| commitlog_sync_batch_window_in_ms | Exists in 3.11, 4.0, trunk | *string | false |
| commitlog_sync_group_window_in_ms | Exists in: 4.0, trunk | *int | false |
| commitlog_sync_period_in_ms | Exists in 3.11, 4.0, trunk | *int | false |
| commitlog_total_space_in_mb | Exists in 3.11, 4.0, trunk | *int | false |
| compaction_large_partition_warning_threshold_mb | Exists in 3.11, 4.0, trunk | *int | false |
| compaction_throughput_mb_per_sec | Exists in 3.11, 4.0, trunk | *int | false |
| compaction_tombstone_warning_threshold | Exists in trunk | *int | false |
| concurrent_compactors | Exists in 3.11, 4.0, trunk | *int | false |
| concurrent_counter_writes | Exists in 3.11, 4.0, trunk | *int | false |
| concurrent_materialized_view_builders | Exists in: 4.0, trunk | *int | false |
| concurrent_materialized_view_writes | Exists in 3.11, 4.0, trunk | *int | false |
| concurrent_reads | Exists in 3.11, 4.0, trunk | *int | false |
| concurrent_replicates | Exists in 3.11, 4.0, trunk | *int | false |
| concurrent_validations | Exists in: 4.0, trunk | *int | false |
| concurrent_writes | Exists in 3.11, 4.0, trunk | *int | false |
| consecutive_message_errors_threshold | Exists in: 4.0, trunk | *int | false |
| corrupted_tombstone_strategy | Exists in: 4.0, trunk | *string | false |
| counter_cache_keys_to_save | Exists in 3.11, 4.0, trunk | *int | false |
| counter_cache_save_period | Exists in 3.11, 4.0, trunk | *int | false |
| counter_cache_size_in_mb | Exists in 3.11, 4.0, trunk | *int | false |
| counter_write_request_timeout_in_ms | Exists in 3.11, 4.0, trunk | *int | false |
| credentials_cache_max_entries | Exists in 3.11, 4.0, trunk | *int | false |
| credentials_update_interval_in_ms | Exists in 3.11, 4.0, trunk | *int | false |
| credentials_validity_in_ms | Exists in 3.11, 4.0, trunk | *int | false |
| cross_node_timeout | Exists in 3.11, 4.0, trunk | *bool | false |
| default_keyspace_rf | Exists in trunk | *int | false |
| denylist_consistency_level | Exists in trunk | *string | false |
| denylist_initial_load_retry_seconds | Exists in trunk | *int | false |
| denylist_max_keys_per_table | Exists in trunk | *int | false |
| denylist_max_keys_total | Exists in trunk | *int | false |
| denylist_refresh_seconds | Exists in trunk | *int | false |
| diagnostic_events_enabled | Exists in: 4.0, trunk | *bool | false |
| disk_access_mode | Exists in 3.11, 4.0, trunk | *string | false |
| disk_optimization_estimate_percentile | Exists in 3.11, 4.0, trunk | *string | false |
| disk_optimization_page_cross_chance | Exists in 3.11, 4.0, trunk | *string | false |
| disk_optimization_strategy | Exists in 3.11, 4.0, trunk | *string | false |
| dynamic_snitch | Exists in 3.11, 4.0, trunk | *bool | false |
| dynamic_snitch_badness_threshold | Exists in 3.11, 4.0, trunk | *string | false |
| dynamic_snitch_reset_interval_in_ms | Exists in 3.11, 4.0, trunk | *int | false |
| dynamic_snitch_update_interval_in_ms | Exists in 3.11, 4.0, trunk | *int | false |
| enable_denylist_range_reads | Exists in trunk | *bool | false |
| enable_denylist_reads | Exists in trunk | *bool | false |
| enable_denylist_writes | Exists in trunk | *bool | false |
| enable_drop_compact_storage | Exists in 3.11, 4.0, trunk | *bool | false |
| enable_materialized_views | Exists in 3.11, 4.0, trunk | *bool | false |
| enable_partition_denylist | Exists in trunk | *bool | false |
| enable_sasi_indexes | Exists in 3.11, 4.0, trunk | *bool | false |
| enable_scripted_user_defined_functions | Exists in 3.11, 4.0, trunk | *bool | false |
| enable_transient_replication | Exists in: 4.0, trunk | *bool | false |
| enable_user_defined_functions | Exists in 3.11, 4.0, trunk | *bool | false |
| enable_user_defined_functions_threads | Exists in 3.11, 4.0, trunk | *bool | false |
| endpoint_snitch | Exists in 3.11, 4.0, trunk | *string | false |
| failure_detector | Exists in trunk | *string | false |
| file_cache_enabled | Exists in: 4.0, trunk | *bool | false |
| file_cache_round_up | Exists in 3.11, 4.0, trunk | *bool | false |
| file_cache_size_in_mb | Exists in 3.11, 4.0, trunk | *int | false |
| flush_compression | Exists in: 4.0, trunk | *string | false |
| full_query_logging_options | Exists in: 4.0, trunk | *[FullQueryLoggerOptions](#fullqueryloggeroptions) | false |
| gc_log_threshold_in_ms | Exists in 3.11, 4.0, trunk | *int | false |
| gc_warn_threshold_in_ms | Exists in 3.11, 4.0, trunk | *int | false |
| hint_window_persistent_enabled | Exists in trunk | *bool | false |
| hinted_handoff_disabled_datacenters | Exists in 3.11, 4.0, trunk | *[]string | false |
| hinted_handoff_enabled | Exists in 3.11, 4.0, trunk | *bool | false |
| hinted_handoff_throttle_in_kb | Exists in 3.11, 4.0, trunk | *int | false |
| hints_compression | Exists in 3.11, 4.0, trunk | *[ParameterizedClass](#parameterizedclass) | false |
| hints_flush_period_in_ms | Exists in 3.11, 4.0, trunk | *int | false |
| ideal_consistency_level | Exists in: 4.0, trunk | *string | false |
| index_interval | Exists in 3.11 | *int | false |
| index_summary_capacity_in_mb | Exists in 3.11, 4.0, trunk | *int | false |
| index_summary_resize_interval_in_minutes | Exists in 3.11, 4.0, trunk | *int | false |
| initial_range_tombstone_list_allocation_size | Exists in: 4.0, trunk | *int | false |
| inter_dc_stream_throughput_outbound_megabits_per_sec | Exists in 3.11, 4.0, trunk | *int | false |
| inter_dc_tcp_nodelay | Exists in 3.11, 4.0, trunk | *bool | false |
| internode_application_receive_queue_capacity_in_bytes | Exists in: 4.0, trunk | *int | false |
| internode_application_receive_queue_reserve_endpoint_capacity_in_bytes | Exists in: 4.0, trunk | *int | false |
| internode_application_receive_queue_reserve_global_capacity_in_bytes | Exists in: 4.0, trunk | *int | false |
| internode_application_send_queue_capacity_in_bytes | Exists in: 4.0, trunk | *int | false |
| internode_application_send_queue_reserve_endpoint_capacity_in_bytes | Exists in: 4.0, trunk | *int | false |
| internode_application_send_queue_reserve_global_capacity_in_bytes | Exists in: 4.0, trunk | *int | false |
| internode_authenticator | Exists in 3.11, 4.0, trunk | *string | false |
| internode_compression | Exists in 3.11, 4.0, trunk | *string | false |
| internode_error_reporting_exclusions | Exists in trunk | *[SubnetGroups](#subnetgroups) | false |
| internode_max_message_size_in_bytes | Exists in: 4.0, trunk | *int | false |
| internode_recv_buff_size_in_bytes | Exists in 3.11 | *int | false |
| internode_send_buff_size_in_bytes | Exists in 3.11 | *int | false |
| internode_socket_receive_buffer_size_in_bytes | Exists in: 4.0, trunk | *int | false |
| internode_socket_send_buffer_size_in_bytes | Exists in: 4.0, trunk | *int | false |
| internode_streaming_tcp_user_timeout_in_ms | Exists in: 4.0, trunk | *int | false |
| internode_tcp_connect_timeout_in_ms | Exists in: 4.0, trunk | *int | false |
| internode_tcp_user_timeout_in_ms | Exists in: 4.0, trunk | *int | false |
| key_cache_keys_to_save | Exists in 3.11, 4.0, trunk | *int | false |
| key_cache_migrate_during_compaction | Exists in: 4.0, trunk | *bool | false |
| key_cache_save_period | Exists in 3.11, 4.0, trunk | *int | false |
| key_cache_size_in_mb | Exists in 3.11, 4.0, trunk | *int | false |
| keyspace_count_warn_threshold | Exists in: 4.0, trunk | *int | false |
| max_concurrent_automatic_sstable_upgrades | Exists in: 4.0, trunk | *int | false |
| max_hint_window_in_ms | Exists in 3.11, 4.0, trunk | *int | false |
| max_hints_delivery_threads | Exists in 3.11, 4.0, trunk | *int | false |
| max_hints_file_size_in_mb | Exists in 3.11, 4.0, trunk | *int | false |
| max_mutation_size_in_kb | Exists in 3.11, 4.0, trunk | *int | false |
| max_streaming_retries | Exists in 3.11, 4.0, trunk | *int | false |
| max_value_size_in_mb | Exists in 3.11, 4.0, trunk | *int | false |
| memtable_allocation_type | Exists in 3.11, 4.0, trunk | *string | false |
| memtable_cleanup_threshold | Exists in 3.11, 4.0, trunk | *string | false |
| memtable_flush_writers | Exists in 3.11, 4.0, trunk | *int | false |
| memtable_heap_space_in_mb | Exists in 3.11, 4.0, trunk | *int | false |
| memtable_offheap_space_in_mb | Exists in 3.11, 4.0, trunk | *int | false |
| min_free_space_per_drive_in_mb | Exists in 3.11, 4.0, trunk | *int | false |
| minimum_keyspace_rf | Exists in trunk | *int | false |
| native_transport_allow_older_protocols | Exists in: 4.0, trunk | *bool | false |
| native_transport_flush_in_batches_legacy | Exists in 3.11, 4.0, trunk | *bool | false |
| native_transport_idle_timeout_in_ms | Exists in: 4.0, trunk | *int | false |
| native_transport_max_concurrent_connections | Exists in 3.11, 4.0, trunk | *int | false |
| native_transport_max_concurrent_connections_per_ip | Exists in 3.11, 4.0, trunk | *int | false |
| native_transport_max_concurrent_requests_in_bytes | Exists in 3.11, 4.0, trunk | *int | false |
| native_transport_max_concurrent_requests_in_bytes_per_ip | Exists in 3.11, 4.0, trunk | *int | false |
| native_transport_max_frame_size_in_mb | Exists in 3.11, 4.0, trunk | *int | false |
| native_transport_max_negotiable_protocol_version | Exists in 3.11, 4.0, trunk | *int | false |
| native_transport_max_requests_per_second | Exists in trunk | *int | false |
| native_transport_max_threads | Exists in 3.11, 4.0, trunk | *int | false |
| native_transport_rate_limiting_enabled | Exists in trunk | *bool | false |
| native_transport_receive_queue_capacity_in_bytes | Exists in: 4.0, trunk | *int | false |
| network_authorizer | Exists in: 4.0, trunk | *string | false |
| networking_cache_size_in_mb | Exists in: 4.0, trunk | *int | false |
| num_tokens | Exists in 3.11, 4.0, trunk | *int | false |
| otc_backlog_expiration_interval_ms | Exists in 3.11 | *int | false |
| otc_coalescing_enough_coalesced_messages | Exists in 3.11, 4.0, trunk | *int | false |
| otc_coalescing_strategy | Exists in 3.11, 4.0, trunk | *string | false |
| otc_coalescing_window_us | Exists in 3.11, 4.0, trunk | *int | false |
| paxos_cache_size_in_mb | Exists in trunk | *int | false |
| periodic_commitlog_sync_lag_block_in_ms | Exists in: 4.0, trunk | *int | false |
| permissions_cache_max_entries | Exists in 3.11, 4.0, trunk | *int | false |
| permissions_update_interval_in_ms | Exists in 3.11, 4.0, trunk | *int | false |
| permissions_validity_in_ms | Exists in 3.11, 4.0, trunk | *int | false |
| phi_convict_threshold | Exists in 3.11, 4.0, trunk | *string | false |
| prepared_statements_cache_size_mb | Exists in 3.11, 4.0, trunk | *int | false |
| range_request_timeout_in_ms | Exists in 3.11, 4.0, trunk | *int | false |
| range_tombstone_list_growth_factor | Exists in: 4.0, trunk | *string | false |
| read_request_timeout_in_ms | Exists in 3.11, 4.0, trunk | *int | false |
| reject_repair_compaction_threshold | Exists in: 4.0, trunk | *int | false |
| repair_command_pool_full_strategy | Exists in: 4.0, trunk | *string | false |
| repair_command_pool_size | Exists in: 4.0, trunk | *int | false |
| repair_session_max_tree_depth | Exists in 3.11, 4.0, trunk | *int | false |
| repair_session_space_in_mb | Exists in: 4.0, trunk | *int | false |
| repaired_data_tracking_for_partition_reads_enabled | Exists in: 4.0, trunk | *bool | false |
| repaired_data_tracking_for_range_reads_enabled | Exists in: 4.0, trunk | *bool | false |
| replica_filtering_protection | Exists in 3.11, 4.0, trunk | *[ReplicaFilteringProtectionOptions](#replicafilteringprotectionoptions) | false |
| report_unconfirmed_repaired_data_mismatches | Exists in: 4.0, trunk | *bool | false |
| request_scheduler | Exists in 3.11 | *string | false |
| request_scheduler_id | Exists in 3.11 | *string | false |
| request_scheduler_options | Exists in 3.11 | *[RequestSchedulerOptions](#requestscheduleroptions) | false |
| request_timeout_in_ms | Exists in 3.11, 4.0, trunk | *int | false |
| role_manager | Exists in 3.11, 4.0, trunk | *string | false |
| roles_cache_max_entries | Exists in 3.11, 4.0, trunk | *int | false |
| roles_update_interval_in_ms | Exists in 3.11, 4.0, trunk | *int | false |
| roles_validity_in_ms | Exists in 3.11, 4.0, trunk | *int | false |
| row_cache_class_name | Exists in 3.11, 4.0, trunk | *string | false |
| row_cache_keys_to_save | Exists in 3.11, 4.0, trunk | *int | false |
| row_cache_save_period | Exists in 3.11, 4.0, trunk | *int | false |
| row_cache_size_in_mb | Exists in 3.11, 4.0, trunk | *int | false |
| server_encryption_options | Exists in 3.11, 4.0, trunk | *encryption.ServerEncryptionOptions | false |
| slow_query_log_timeout_in_ms | Exists in 3.11, 4.0, trunk | *int | false |
| snapshot_before_compaction | Exists in 3.11, 4.0, trunk | *bool | false |
| snapshot_links_per_second | Exists in: 4.0, trunk | *int | false |
| snapshot_on_duplicate_row_detection | Exists in 3.11, 4.0, trunk | *bool | false |
| snapshot_on_repaired_data_mismatch | Exists in: 4.0, trunk | *bool | false |
| sstable_preemptive_open_interval_in_mb | Exists in 3.11, 4.0, trunk | *int | false |
| stream_entire_sstables | Exists in: 4.0, trunk | *bool | false |
| stream_throughput_outbound_megabits_per_sec | Exists in 3.11, 4.0, trunk | *int | false |
| streaming_connections_per_host | Exists in: 4.0, trunk | *int | false |
| streaming_keep_alive_period_in_secs | Exists in 3.11, 4.0, trunk | *int | false |
| streaming_socket_timeout_in_ms | Exists in 3.11 | *int | false |
| table_count_warn_threshold | Exists in: 4.0, trunk | *int | false |
| thrift_framed_transport_size_in_mb | Exists in 3.11 | *int | false |
| thrift_max_message_length_in_mb | Exists in 3.11 | *int | false |
| thrift_prepared_statements_cache_size_mb | Exists in 3.11 | *int | false |
| tombstone_failure_threshold | Exists in 3.11, 4.0, trunk | *int | false |
| tombstone_warn_threshold | Exists in 3.11, 4.0, trunk | *int | false |
| tracetype_query_ttl | Exists in 3.11, 4.0, trunk | *int | false |
| tracetype_repair_ttl | Exists in 3.11, 4.0, trunk | *int | false |
| track_warnings | Exists in trunk | *[TrackWarnings](#trackwarnings) | false |
| traverse_auth_from_root | Exists in trunk | *bool | false |
| trickle_fsync | Exists in 3.11, 4.0, trunk | *bool | false |
| trickle_fsync_interval_in_kb | Exists in 3.11, 4.0, trunk | *int | false |
| truncate_request_timeout_in_ms | Exists in 3.11, 4.0, trunk | *int | false |
| unlogged_batch_across_partitions_warn_threshold | Exists in 3.11, 4.0, trunk | *int | false |
| use_deterministic_table_id | Exists in trunk | *bool | false |
| use_offheap_merkle_trees | Exists in: 4.0, trunk | *bool | false |
| user_defined_function_fail_timeout | Exists in 3.11, 4.0, trunk | *int | false |
| user_defined_function_warn_timeout | Exists in 3.11, 4.0, trunk | *int | false |
| user_function_timeout_policy | Exists in 3.11, 4.0, trunk | *string | false |
| validation_preview_purge_head_start_in_sec | Exists in: 4.0, trunk | *int | false |
| windows_timer_interval | Exists in 3.11, 4.0, trunk | *int | false |
| write_request_timeout_in_ms | Exists in 3.11, 4.0, trunk | *int | false |

[Back to Custom Resources](#custom-resources)

#### FullQueryLoggerOptions



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| archive_command |  | *string | false |
| roll_cycle |  | *string | false |
| block |  | *bool | false |
| max_queue_weight |  | *int | false |
| max_log_size |  | *int | false |
| max_archive_retries |  | *int | false |
| log_dir |  | *string | false |

[Back to Custom Resources](#custom-resources)

#### Group



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| subnet |  | string | true |

[Back to Custom Resources](#custom-resources)

#### JvmOptions



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| heapSize |  | *resource.Quantity | false |
| heapNewGenSize |  | *resource.Quantity | false |
| additionalOptions |  | []string | false |

[Back to Custom Resources](#custom-resources)

#### ParameterizedClass



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| class_name |  | string | true |
| parameters |  | *map[string]string | false |

[Back to Custom Resources](#custom-resources)

#### ReplicaFilteringProtectionOptions



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| cached_rows_warn_threshold |  | *int | false |
| cached_rows_fail_threshold |  | *int | false |

[Back to Custom Resources](#custom-resources)

#### RequestSchedulerOptions



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| throttle_limit |  | *int | false |
| default_weight |  | *int | false |
| weights |  | *map[string]int | false |

[Back to Custom Resources](#custom-resources)

#### SubnetGroups



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| subnets |  | [][Group](#group) | true |

[Back to Custom Resources](#custom-resources)

#### TrackWarnings



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| enabled |  | bool | true |
| coordinator_read_size |  | *int | false |
| local_read_size |  | *int | false |
| row_index_size |  | *int | false |

[Back to Custom Resources](#custom-resources)
