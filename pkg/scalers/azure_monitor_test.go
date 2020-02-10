package scalers

import (
	"context"
	"testing"
)

var testAzMonitorResolvedEnv = map[string]string{
	"XXX": "xxx",
}

type parseAzMonitorMetadataTestData struct {
	metadata    map[string]string
	isError     bool
	resolvedEnv map[string]string
	authParams  map[string]string
}

var testParseAzMonitorMetadata = []parseAzMonitorMetadataTestData{
	// nothing passed
	{map[string]string{}, true, testAzMonitorResolvedEnv, map[string]string{}},
	// properly formed
	{map[string]string{"resourceURI": "test/resource/uri", "tenantId": "123", "subscriptionId": "456", "resourceGroupName": "test", "metricName": "metric", "metricAggregationInterval": "0:15:0", "metricAggregationType": "Average", "activeDirectoryClientId": "789", "activeDirectoryClientPassword": "1234", "targetValue": "5"}, false, testAzMonitorResolvedEnv, map[string]string{}},
	// no optional parameters
	{map[string]string{"resourceURI": "test/resource/uri", "tenantId": "123", "subscriptionId": "456", "resourceGroupName": "test", "metricName": "metric", "metricAggregationType": "Average", "activeDirectoryClientId": "789", "activeDirectoryClientPassword": "1234", "targetValue": "5"}, false, testAzMonitorResolvedEnv, map[string]string{}},
	// incorrectly formatted resourceURI
	{map[string]string{"resourceURI": "bad/format", "tenantId": "123", "subscriptionId": "456", "resourceGroupName": "test", "metricName": "metric", "metricAggregationInterval": "0:15:0", "metricAggregationType": "Average", "activeDirectoryClientId": "789", "activeDirectoryClientPassword": "1234", "targetValue": "5"}, true, testAzMonitorResolvedEnv, map[string]string{}},
	// improperly formatted aggregationInterval
	{map[string]string{"resourceURI": "test/resource/uri", "tenantId": "123", "subscriptionId": "456", "resourceGroupName": "test", "metricName": "metric", "metricAggregationInterval": "0:1", "metricAggregationType": "Average", "activeDirectoryClientId": "789", "activeDirectoryClientPassword": "1234", "targetValue": "5"}, true, testAzMonitorResolvedEnv, map[string]string{}},
	// missing resourceURI
	{map[string]string{"tenantId": "123", "subscriptionId": "456", "resourceGroupName": "test", "metricName": "metric", "metricAggregationInterval": "0:15:0", "metricAggregationType": "Average", "activeDirectoryClientId": "789", "activeDirectoryClientPassword": "1234", "targetValue": "5"}, true, testAzMonitorResolvedEnv, map[string]string{}},
	// missing tenantID
	{map[string]string{"resourceURI": "test/resource/uri", "subscriptionId": "456", "resourceGroupName": "test", "metricName": "metric", "metricAggregationInterval": "0:15:0", "metricAggregationType": "Average", "activeDirectoryClientId": "789", "activeDirectoryClientPassword": "1234", "targetValue": "5"}, true, testAzMonitorResolvedEnv, map[string]string{}},
	// missing subscriptionID
	{map[string]string{"resourceURI": "test/resource/uri", "tenantId": "123", "resourceGroupName": "test", "metricName": "metric", "metricAggregationInterval": "0:15:0", "metricAggregationType": "Average", "activeDirectoryClientId": "789", "activeDirectoryClientPassword": "1234", "targetValue": "5"}, true, testAzMonitorResolvedEnv, map[string]string{}},
	// missing resourceGroupName
	{map[string]string{"resourceURI": "test/resource/uri", "tenantId": "123", "subscriptionId": "456", "metricName": "metric", "metricAggregationInterval": "0:15:0", "metricAggregationType": "Average", "activeDirectoryClientId": "789", "activeDirectoryClientPassword": "1234", "targetValue": "5"}, true, testAzMonitorResolvedEnv, map[string]string{}},
	// missing metric name
	{map[string]string{"resourceURI": "test/resource/uri", "tenantId": "123", "subscriptionId": "456", "resourceGroupName": "test", "metricAggregationInterval": "0:15:0", "metricAggregationType": "Average", "activeDirectoryClientId": "789", "activeDirectoryClientPassword": "1234", "targetValue": "5"}, true, testAzMonitorResolvedEnv, map[string]string{}},
	// filter included
	{map[string]string{"resourceURI": "test/resource/uri", "tenantId": "123", "subscriptionId": "456", "resourceGroupName": "test", "metricName": "metric", "metricFilter": "namespace eq 'default'", "metricAggregationInterval": "0:15:0", "metricAggregationType": "Average", "activeDirectoryClientId": "789", "activeDirectoryClientPassword": "1234", "targetValue": "5"}, false, testAzMonitorResolvedEnv, map[string]string{}},
	// invalid aggregation type
	{map[string]string{"resourceURI": "test/resource/uri", "tenantId": "123", "subscriptionId": "456", "resourceGroupName": "test", "metricName": "metric", "metricAggregationInterval": "0:15:0", "metricAggregationType": "Average", "activeDirectoryClientId": "789", "activeDirectoryClientPassword": "1234", "targetValue": "5"}, true, testAzMonitorResolvedEnv, map[string]string{}},
	// missing clientID
	{map[string]string{"resourceURI": "test/resource/uri", "tenantId": "123", "subscriptionId": "456", "resourceGroupName": "test", "metricName": "metric", "metricAggregationInterval": "0:15:0", "metricAggregationType": "Average", "activeDirectoryClientId": "789", "activeDirectoryClientPassword": "1234", "targetValue": "5"}, true, testAzMonitorResolvedEnv, map[string]string{}},
	// missing clientPassword
	{map[string]string{"resourceURI": "test/resource/uri", "tenantId": "123", "subscriptionId": "456", "resourceGroupName": "test", "metricName": "metric", "metricAggregationInterval": "0:15:0", "metricAggregationType": "Average", "activeDirectoryClientId": "789", "activeDirectoryClientPassword": "1234", "targetValue": "5"}, true, testAzMonitorResolvedEnv, map[string]string{}},
	// missing targetValue
	{map[string]string{"resourceURI": "test/resource/uri", "tenantId": "123", "subscriptionId": "456", "resourceGroupName": "test", "metricName": "metric", "metricAggregationInterval": "0:15:0", "metricAggregationType": "Average", "activeDirectoryClientId": "789", "activeDirectoryClientPassword": "1234", "targetValue": "5"}, true, testAzMonitorResolvedEnv, map[string]string{}},
}

func TestAzMonitorParseMetadata(t *testing.T) {
	for _, testData := range testParseAzMonitorMetadata {
		_, err := parseAzureMonitorMetadata(testData.metadata, testData.resolvedEnv, testData.authParams)
		if err != nil && !testData.isError {
			t.Error("Expected success but got error", err)
		}
		if testData.isError && err == nil {
			t.Errorf("Expected error but got success. testData: %v", testData)
		}
	}
}

type getAzMonitorMetricValueTestData struct {
	metadata    *azureMonitorMetadata
	isError     bool
	resolvedEnv map[string]string
	authParams  map[string]string
}

//test for authparams and such

/*var testAzMonitorMetadata = []getAzMonitorMetricValueTestData{
	// nothing passed
	{&azureMonitorMetadata{}, true, testAzMonitorResolvedEnv, map[string]string{}},
	// properly formed
	{&azureMonitorMetadata{resourceURI: "test/resource/uri", tenantID: "123", subscriptionID: "456", resourceGroupName: "test", name: "metric", aggregationInterval: "0:15:0", aggregationType: "Average", clientID: "789", clientPassword: "1234", targetValue: 5}, false, testAzMonitorResolvedEnv, map[string]string{}},
	// no optional parameters
	{&azureMonitorMetadata{resourceURI: "test/resource/uri", tenantID: "123", subscriptionID: "456", resourceGroupName: "test", name: "metric", aggregationType: "Average", clientID: "789", clientPassword: "1234", targetValue: 5}, false, testAzMonitorResolvedEnv, map[string]string{}},
	// incorrectly formatted resourceURI
	{&azureMonitorMetadata{resourceURI: "bad/format", tenantID: "123", subscriptionID: "456", resourceGroupName: "test", name: "metric", aggregationInterval: "0:15:0", aggregationType: "Average", clientID: "789", clientPassword: "1234", targetValue: 5}, true, testAzMonitorResolvedEnv, map[string]string{}},
	// improperly formatted aggregationInterval
	{&azureMonitorMetadata{resourceURI: "test/resource/uri", tenantID: "123", subscriptionID: "456", resourceGroupName: "test", name: "metric", aggregationInterval: "0:1", aggregationType: "Average", clientID: "789", clientPassword: "1234", targetValue: 5}, true, testAzMonitorResolvedEnv, map[string]string{}},
	// missing resourceURI
	{&azureMonitorMetadata{tenantID: "123", subscriptionID: "456", resourceGroupName: "test", name: "metric", aggregationInterval: "0:15:0", aggregationType: "Average", clientID: "789", clientPassword: "1234", targetValue: 5}, true, testAzMonitorResolvedEnv, map[string]string{}},
	// missing tenantID
	{&azureMonitorMetadata{resourceURI: "test/resource/uri", subscriptionID: "456", resourceGroupName: "test", name: "metric", aggregationInterval: "0:15:0", aggregationType: "Average", clientID: "789", clientPassword: "1234", targetValue: 5}, true, testAzMonitorResolvedEnv, map[string]string{}},
	// missing subscriptionID
	{&azureMonitorMetadata{resourceURI: "test/resource/uri", tenantID: "123", resourceGroupName: "test", name: "metric", aggregationInterval: "0:15:0", aggregationType: "Average", clientID: "789", clientPassword: "1234", targetValue: 5}, true, testAzMonitorResolvedEnv, map[string]string{}},
	// missing resourceGroupName
	{&azureMonitorMetadata{resourceURI: "test/resource/uri", tenantID: "123", subscriptionID: "456", name: "metric", aggregationInterval: "0:15:0", aggregationType: "Average", clientID: "789", clientPassword: "1234", targetValue: 5}, true, testAzMonitorResolvedEnv, map[string]string{}},
	// missing metric name
	{&azureMonitorMetadata{resourceURI: "test/resource/uri", tenantID: "123", subscriptionID: "456", resourceGroupName: "test", aggregationInterval: "0:15:0", aggregationType: "Average", clientID: "789", clientPassword: "1234", targetValue: 5}, true, testAzMonitorResolvedEnv, map[string]string{}},
	// filter included
	{&azureMonitorMetadata{resourceURI: "test/resource/uri", tenantID: "123", subscriptionID: "456", resourceGroupName: "test", name: "metric", filter: "namespace eq 'default'", aggregationInterval: "0:15:0", aggregationType: "Average", clientID: "789", clientPassword: "1234", targetValue: 5}, false, testAzMonitorResolvedEnv, map[string]string{}},
	// invalid aggregation type
	{&azureMonitorMetadata{resourceURI: "test/resource/uri", tenantID: "123", subscriptionID: "456", resourceGroupName: "test", name: "metric", aggregationInterval: "0:15:0", aggregationType: "Avg", clientID: "789", clientPassword: "1234", targetValue: 5}, true, testAzMonitorResolvedEnv, map[string]string{}},
	// missing clientID
	{&azureMonitorMetadata{resourceURI: "test/resource/uri", tenantID: "123", subscriptionID: "456", resourceGroupName: "test", name: "metric", aggregationInterval: "0:15:0", aggregationType: "Average", clientPassword: "1234", targetValue: 5}, true, testAzMonitorResolvedEnv, map[string]string{}},
	// missing clientPassword
	{&azureMonitorMetadata{resourceURI: "test/resource/uri", tenantID: "123", subscriptionID: "456", resourceGroupName: "test", name: "metric", aggregationInterval: "0:15:0", aggregationType: "Average", clientID: "789", targetValue: 5}, true, testAzMonitorResolvedEnv, map[string]string{}},
	// missing targetValue
	{&azureMonitorMetadata{resourceURI: "test/resource/uri", tenantID: "123", subscriptionID: "456", resourceGroupName: "test", name: "metric", aggregationInterval: "0:15:0", aggregationType: "Average", clientID: "789", clientPassword: "1234"}, true, testAzMonitorResolvedEnv, map[string]string{}},
}

func TestGetAzureMetricValue(t *testing.T) {
	for _, testData := range testAzMonitorMetadata {
		_, err := GetAzureMetricValue(context.TODO(), testData.metadata)
		if err != nil && !testData.isError {
			t.Error("Expected success but got error", err)
		}
		if testData.isError && err == nil {
			t.Errorf("Expected error but got success. testData: %v", testData)
		}
	}
	/*length, err := GetAzureMetricValue(context.TODO(), testData.metadata)
	if length != -1 {
		t.Error("Expected length to be -1, but got", length)
	}

	if err == nil {
		t.Error("Expected error for empty connection string, but got nil")
	}

	if !strings.Contains(err.Error(), "parse storage connection string") {
		t.Error("Expected error to contain parsing error message, but got", err.Error())
	}

	length, err = GetAzureMetricValue(context.TODO())

	if length != -1 {
		t.Error("Expected length to be -1, but got", length)
	}

	if err == nil {
		t.Error("Expected error for empty connection string, but got nil")
	}

	if !strings.Contains(err.Error(), "illegal base64") {
		t.Error("Expected error to contain base64 error message, but got", err.Error())
	}
}*/
