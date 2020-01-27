package scalers

import (
	"context"
	"fmt"
	"strconv"

	v2beta1 "k8s.io/api/autoscaling/v2beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/metrics/pkg/apis/external_metrics"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	azureMonitorMetricName = "metricName"
	targetValueName        = "targetValue"
)

type azureMonitorScaler struct {
	metadata *azureMonitorMetadata
}

type azureMonitorMetadata struct {
	resourceURI          string
	tentantID            string
	subscriptionID       string
	resourceGroupName    string
	name                 string
	filter               string
	aggregationInterval  string
	aggregationType      string
	servicePrincipalID   string
	servicePrincipalPass string
	targetValue          int
}

var azureMonitorLog = logf.Log.WithName("azure_monitor_scaler")

// NewAzureMonitorScaler stuff
func NewAzureMonitorScaler(resolvedEnv, metadata, authParams map[string]string) (Scaler, error) {
	meta, err := parseAzureMonitorMetadata(metadata, resolvedEnv, authParams)
	if err != nil {
		return nil, fmt.Errorf("error parsing azure monitor metadata: %s", err)
	}

	return &azureMonitorScaler{
		metadata: meta,
	}, nil
}

func parseAzureMonitorMetadata(metadata, resolvedEnv, authParams map[string]string) (*azureMonitorMetadata, error) {
	meta := azureMonitorMetadata{}

	if val, ok := metadata[targetValueName]; ok && val != "" {
		targetValue, err := strconv.Atoi(val)
		if err != nil {
			azureMonitorLog.Error(err, "Error parsing azure monitor metadata", "targetValue", targetValueName)
			return nil, fmt.Errorf("Error parsing azure monitor metadata %s: %s", targetValueName, err.Error())
		}

		meta.targetValue = targetValue
	}

	if val, ok := metadata["resourceURI"]; ok && val != "" {
		meta.resourceURI = val
	} else {
		return nil, fmt.Errorf("no resourceURI given")
	}

	if val, ok := metadata["tenantId"]; ok && val != "" {
		meta.tentantID = val
	} else {
		return nil, fmt.Errorf("no tenantId given")
	}

	if val, ok := metadata["subscriptionId"]; ok && val != "" {
		meta.subscriptionID = val
	} else {
		return nil, fmt.Errorf("no subscriptionId given")
	}

	if val, ok := metadata["resourceGroupName"]; ok && val != "" {
		meta.resourceGroupName = val
	} else {
		return nil, fmt.Errorf("no resourceGroupName given")
	}

	if val, ok := metadata[azureMonitorMetricName]; ok && val != "" {
		meta.name = val
	} else {
		return nil, fmt.Errorf("no metricName given")
	}

	if val, ok := metadata["metricFilter"]; ok {
		if val != "" {
			meta.filter = val
		}
	}

	if val, ok := metadata["metricAggregationInterval"]; ok {
		if val != "" {
			meta.aggregationInterval = val
		}
	}

	if val, ok := metadata["adServicePrincipleId"]; ok && val != "" {
		meta.servicePrincipalID = val
	} else {
		return nil, fmt.Errorf("no adServicePrincipleId given")
	}

	if val, ok := metadata["adServicePrinciplePassword"]; ok {
		meta.servicePrincipalPass = val
	} else {
		return nil, fmt.Errorf("no adServicePrinciplePassword given")
	}

	return &meta, nil
}

// needs to interact with azure monitor
func (s *azureMonitorScaler) IsActive(ctx context.Context) (bool, error) {
	val, err := GetAzureMetricValue(ctx, s.metadata)
	if err != nil {
		azureMonitorLog.Error(err, "error getting azure monitor metric")
		return false, err
	}

	return val > 0, nil
}

func (s *azureMonitorScaler) Close() error {
	return nil
}

func (s *azureMonitorScaler) GetMetricSpecForScaling() []v2beta1.MetricSpec {
	targetMetricVal := resource.NewQuantity(int64(s.metadata.targetValue), resource.DecimalSI)
	externalMetric := &v2beta1.ExternalMetricSource{MetricName: azureMonitorMetricName, TargetAverageValue: targetMetricVal}
	metricSpec := v2beta1.MetricSpec{External: externalMetric, Type: externalMetricType}
	return []v2beta1.MetricSpec{metricSpec}
}

// GetMetrics returns value for a supported metric and an error if there is a problem getting the metric
func (s *azureMonitorScaler) GetMetrics(ctx context.Context, metricName string, metricSelector labels.Selector) ([]external_metrics.ExternalMetricValue, error) {
	val, err := GetAzureMetricValue(ctx, s.metadata)
	if err != nil {
		azureMonitorLog.Error(err, "error getting azure monitor metric")
		return []external_metrics.ExternalMetricValue{}, err
	}

	metric := external_metrics.ExternalMetricValue{
		MetricName: metricName,
		Value:      *resource.NewQuantity(int64(val), resource.DecimalSI),
		Timestamp:  metav1.Now(),
	}

	return append([]external_metrics.ExternalMetricValue{}, metric), nil
}
