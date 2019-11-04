package ocean

const (
	FlagName                       string = "name"
	FlagRegion                     string = "region"
	FlagClusterID                  string = "cluster-id"
	FlagControllerClusterID        string = "controller-cluster-id"
	FlagSpotPercentage             string = "spot-percentage"
	FlagDrainingTimeout            string = "draining-timeout"
	FlagUtilizeReserveInstances    string = "utilize-reserved-instances"
	FlagFallbackOnDemand           string = "fallback-ondemand"
	FlagMinSize                    string = "min-size"
	FlagMaxSize                    string = "max-size"
	FlagTargetSize                 string = "target-size"
	FlagSubnetIDs                  string = "subnet-ids"
	FlagInstancesTypesWhitelist    string = "instance-types-whitelist"
	FlagInstancesTypesBlacklist    string = "instance-types-blacklist"
	FlagSecurityGroupIds           string = "security-group-ids"
	FlagImageIDs                   string = "image-id"
	FlagKeyPair                    string = "key-pair"
	FlagUserData                   string = "user-data"
	FlagRootVolumeSize             string = "root-volume-size"
	FlagAssociatePublicIPAddress   string = "associate-public-ip-address"
	FlagEnableMonitoring           string = "enable-monitoring"
	FlagEnableEBSOptimization      string = "enable-ebs-optimization"
	FlagIamInstanceProfileName     string = "iam-instance-profile-name"
	FlagIamInstanceProfileARN      string = "iam-instance-profile-arn"
	FlagLoadBalancerName           string = "load-balancer-name"
	FlagLoadBalancerARN            string = "load-balancer-arn"
	FlagLoadBalancerType           string = "load-balancer-type"
	FlagEnableAutoScaler           string = "enable-auto-scaler"
	FlagEnableAutoScalerAutoConfig string = "enable-auto-scaler-autoconfig"
	FlagCooldown                   string = "cooldown"
	FlagHeadroomCPUPerUnit         string = "headroom-cpu-per-unit"
	FlagHeadroomMemoryPerUnit      string = "headroom-memory-per-unit"
	FlagHeadroomGPUPerUnit         string = "headroom-gpu-per-unit"
	FlagHeadroomNumPerUnit         string = "headroom-num-per-unit"
	FlagResourceLimitMaxVCPU       string = "resource-limit-max-vcpu"
	FlagResourceLimitMaxMemory     string = "resource-limit-max-memory"
	FlagEvaluationPeriods          string = "evaluation-periods"
	FlagMaxScaleDownPercentage     string = "max-scale-down-percentage"
)