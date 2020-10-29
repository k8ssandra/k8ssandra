package structs

// CassDC a representation of the cassandra datacenter.
type CassDC struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		ClusterName string `yaml:"clusterName"`
		Size        int    `yaml:"size"`
		Racks       []struct {
			Name string `yaml:"name"`
			Zone string `yaml:"zone"`
		} `yaml:"racks"`
		Resources struct {
			Requests struct {
				Memory string `yaml:"memory"`
				CPU    string `yaml:"cpu"`
			} `yaml:"requests"`
			Limits struct {
				Memory string `yaml:"memory"`
				CPU    string `yaml:"cpu"`
			} `yaml:"limits"`
		} `yaml:"resources"`
		StorageConfig struct {
			CassandraDataVolumeClaimSpec struct {
				StorageClassName string   `yaml:"storageClassName"`
				AccessModes      []string `yaml:"accessModes"`
				Resources        struct {
					Requests struct {
						Storage string `yaml:"storage"`
					} `yaml:"requests"`
				} `yaml:"resources"`
			} `yaml:"cassandraDataVolumeClaimSpec"`
		} `yaml:"storageConfig"`
		AllowMultipleNodesPerWorker bool   `yaml:"allowMultipleNodesPerWorker"`
		Stopped                     bool   `yaml:"stopped"`
		RollingRestartRequested     bool   `yaml:"rollingRestartRequested"`
		CanaryUpgrade               bool   `yaml:"canaryUpgrade"`
		ServerType                  string `yaml:"serverType"`
		ServerVersion               string `yaml:"serverVersion"`
		ServerImage                 string `yaml:"serverImage"`
		ConfigBuilderImage          string `yaml:"configBuilderImage"`
		SuperuserSecretName         string `yaml:"superuserSecretName"`
		ManagementAPIAuth           struct {
			Insecure struct {
			} `yaml:"insecure"`
		} `yaml:"managementApiAuth"`
		ServiceAccount string        `yaml:"serviceAccount"`
		ReplaceNodes   []interface{} `yaml:"replaceNodes"`
		Config         struct {
			CassandraYaml struct {
				NumTokens         int    `yaml:"num_tokens"`
				FileCacheSizeInMb int    `yaml:"file_cache_size_in_mb"`
				Authenticator     string `yaml:"authenticator"`
				Authorizer        string `yaml:"authorizer"`
				RoleManager       string `yaml:"role_manager"`
			} `yaml:"cassandra-yaml"`
			JvmOptions struct {
				InitialHeapSize   string   `yaml:"initial_heap_size"`
				MaxHeapSize       string   `yaml:"max_heap_size"`
				AdditionalJvmOpts []string `yaml:"additional-jvm-opts"`
			} `yaml:"jvm-options"`
		} `yaml:"config"`
	} `yaml:"spec"`
}
