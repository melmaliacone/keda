package scalers

import (
	"context"
	"testing"
)

var testAzMonitorResolvedEnv = map[string]string{
	"CONNECTION": "SAMPLE",
}

type parseAzMonitorMetadataTestData struct {
	metadata    map[string]string
	isError     bool
	resolvedEnv map[string]string
	authParams  map[string]string
}

var testParseAzMonitorMetadata = []parseAzQueueMetadataTestData{
	// nothing passed
	{map[string]string{}, true, testAzQueueResolvedEnv, map[string]string{}, ""},
	// properly formed
	{map[string]string{"resourceURI": "CONNECTION", "tenantID": "sample", "subscriptionID": "5"}, false, testAzQueueResolvedEnv, map[string]string{}, ""},
	// Empty queueName
	{map[string]string{"connection": "CONNECTION", "queueName": ""}, true, testAzQueueResolvedEnv, map[string]string{}, ""},
	// improperly formed target
	{map[string]string{"connection": "CONNECTION", "queueName": "sample", "queueLength": "AA"}, true, testAzQueueResolvedEnv, map[string]string{}, ""},
	// connection from authParams
	{map[string]string{"queueName": "sample", "queueLength": "5"}, false, testAzQueueResolvedEnv, map[string]string{"connection": "value"}, "none"},
}

func TestAzMonitorParseMetadata(t *testing.T) {
	for _, testData := range parseAzQueueMetadataTestData {
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

var testAzMonitorMetadata = []testGetAzMetricValueTestData{
	// nothing passed
	{&azureMonitorMetadata{}, true, testAzMonitorResolvedEnv, map[string]string{}},
	// properly formed
	{&azureMonitorMetadata{resourceURI: "CONNECTION", tenantID: "sample", subscriptionID: "5", resourceGroupName: "", name: "", filter: "", aggregationInterval: "", aggregationType: "", servicePrincipalID: "", servicePrincipalPass: "", targetValue: 5}, false, testAzMonitorResolvedEnv, map[string]string{}},
	// Empty queueName
	{&azureMonitorMetadata{resourceURI: "CONNECTION", tenantID: ""}, true, testAzMonitorResolvedEnv, map[string]string{}},
	// improperly formed queueLength
	{&azureMonitorMetadata{resourceURI: "CONNECTION", tenantID: "sample", subscriptionID: "5"}, true, testAzMonitorResolvedEnv, map[string]string{}},
	// connection from authParams
	{&azureMonitorMetadata{resourceURI: "sample", tenantID: "5"}, false, testAzMonitorResolvedEnv, map[string]string{"connection": "value"}},
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
	}*/
}
