package spot

import (
	"reflect"
	"time"
)

type (
	// TypeMeta describes an individual object.
	TypeMeta struct {
		// Kind represents the kind of the resource this object represents.
		Kind string `json:"kind"`
	}

	// ObjectMeta is metadata that all resources must have.
	ObjectMeta struct {
		// Unique ID.
		ID string `json:"id" table:"1,id"`

		// Object name.
		Name string `json:"name" table:"2,name"`

		// CreatedAt represents the timestamp when the cluster has been created.
		CreatedAt time.Time `json:"createdAt" table:"9,created"`

		// UpdatedAt represents the timestamp when the cluster has been updated.
		UpdatedAt time.Time `json:"updatedAt" table:"10,updated"`
	}

	// Account represents a Spot account.
	Account struct {
		// Unique account ID.
		ID string `json:"accountId"`

		// Account name.
		Name string `json:"name"`

		// Organization ID.
		OrganizationID string `json:"organizationId"`

		// External ID generated by the cloud provider bounded with this account.
		ExternalID *string `json:"providerExternalId"` // nullable; used by AWS only
	}

	// OceanCluster represents an Ocean cluster.
	OceanCluster struct {
		// Type's metadata.
		TypeMeta

		// Object's metadata.
		ObjectMeta

		// Obj holds the raw object which is an orchestrator-specific implementation.
		Obj interface{} `json:"-"`
	}

	// OceanLaunchSpec represents an Ocean launch spec.
	OceanLaunchSpec struct {
		// Type's metadata.
		TypeMeta

		// Object's metadata.
		ObjectMeta

		// Obj holds the raw object which is an orchestrator-specific implementation.
		Obj interface{} `json:"-"`
	}

	// OceanRollout represents an Ocean rollout.
	OceanRollout struct {
		// Type's metadata.
		TypeMeta

		// Object's metadata.
		ObjectMeta

		// Obj holds the raw object which is an orchestrator-specific implementation.
		Obj interface{} `json:"-"`
	}

	// OceanClusterOptions represents an Ocean cluster.
	OceanClusterOptions struct {
		// Base.
		ClusterID    string
		ControllerID string
		Name         string
		Region       string

		// Strategy.
		SpotPercentage           float64
		UtilizeReservedInstances bool
		FallbackToOnDemand       bool
		DrainingTimeout          int

		// Capacity.
		MinSize    int
		MaxSize    int
		TargetSize int

		// Compute
		SubnetIDs                []string
		InstanceTypesWhitelist   []string
		InstanceTypesBlacklist   []string
		SecurityGroupIDs         []string
		ImageID                  string
		KeyPair                  string
		UserData                 string
		RootVolumeSize           int
		AssociatePublicIPAddress bool
		EnableMonitoring         bool
		EnableEBSOptimization    bool
		IAMInstanceProfileName   string
		IAMInstanceProfileARN    string

		LoadBalancerNames []string
		LoadBalancerARNs  []string
		LoadBalancerType  string // Deprecated: Inferred from name/arn.

		// Auto Scaling.
		EnableAutoScaler       bool
		EnableAutoConfig       bool
		Cooldown               int
		HeadroomCPUPerUnit     int
		HeadroomMemoryPerUnit  int
		HeadroomGPUPerUnit     int
		HeadroomNumPerUnit     int
		ResourceLimitMaxVCPU   int
		ResourceLimitMaxMemory int
		EvaluationPeriods      int
		MaxScaleDownPercentage int
	}

	// OceanLaunchSpecOptions represents an Ocean launch spec.
	OceanLaunchSpecOptions struct {
		// Base.
		SpecID    string
		ClusterID string
		Name      string

		// Compute.
		ImageID          string
		UserData         string
		SecurityGroupIDs []string
	}

	// OceanRolloutOptions represents an Ocean rollout.
	OceanRolloutOptions struct {
		// Base.
		RolloutID string
		ClusterID string
		Comment   string
		Status    string

		// Parameters.
		BatchSizePercentage int
		DisableAutoScaling  bool
		SpecIDs             []string
		InstanceIDs         []string
	}
)

// typeOf returns obj type's name using reflection.
func typeOf(obj interface{}) string { return reflect.TypeOf(obj).Name() }
