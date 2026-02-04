package storage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v2"
)

// TestAzureConfig_validate tests the validation logic for AzureConfig
func TestAzureConfig_validate(t *testing.T) {
	tests := []struct {
		name    string
		config  AzureConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config with storage account key",
			config: AzureConfig{
				StorageAccountName: "testaccount",
				StorageAccountKey:  "testkey",
				ContainerName:      "testcontainer",
			},
			wantErr: false,
		},
		{
			name: "valid config with connection string",
			config: AzureConfig{
				StorageAccountName:      "testaccount",
				StorageConnectionString: "DefaultEndpointsProtocol=https;AccountName=test;AccountKey=key;EndpointSuffix=core.windows.net",
				ContainerName:           "testcontainer",
			},
			wantErr: false,
		},
		{
			name: "valid config with user assigned ID",
			config: AzureConfig{
				StorageAccountName: "testaccount",
				UserAssignedID:     "test-msi-id",
				ContainerName:      "testcontainer",
			},
			wantErr: false,
		},
		{
			name: "valid config with client secret",
			config: AzureConfig{
				StorageAccountName: "testaccount",
				TenantID:           "tenant-id",
				ClientID:           "client-id",
				ClientSecret:       "client-secret",
				ContainerName:      "testcontainer",
			},
			wantErr: false,
		},
		{
			name: "missing storage account name",
			config: AzureConfig{
				StorageAccountKey: "testkey",
				ContainerName:     "testcontainer",
			},
			wantErr: true,
			errMsg:  "storage_account_name is required but not configured",
		},
		{
			name: "missing container name",
			config: AzureConfig{
				StorageAccountName: "testaccount",
				StorageAccountKey:  "testkey",
			},
			wantErr: true,
			errMsg:  "no container specified",
		},
		{
			name: "user_assigned_id with storage_account_key",
			config: AzureConfig{
				StorageAccountName: "testaccount",
				StorageAccountKey:  "testkey",
				UserAssignedID:     "test-msi-id",
				ContainerName:      "testcontainer",
			},
			wantErr: true,
			errMsg:  "user_assigned_id cannot be set when using storage_account_key authentication",
		},
		{
			name: "user_assigned_id with storage_connection_string",
			config: AzureConfig{
				StorageAccountName:      "testaccount",
				StorageConnectionString: "connection-string",
				UserAssignedID:          "test-msi-id",
				ContainerName:           "testcontainer",
			},
			wantErr: true,
			errMsg:  "user_assigned_id cannot be set when using storage_connection_string authentication",
		},
		{
			name: "storage_account_key with storage_connection_string",
			config: AzureConfig{
				StorageAccountName:      "testaccount",
				StorageAccountKey:       "testkey",
				StorageConnectionString: "connection-string",
				ContainerName:           "testcontainer",
			},
			wantErr: true,
			errMsg:  "storage_account_key and storage_connection_string cannot both be set",
		},
		{
			name: "tenant_id with storage_account_key",
			config: AzureConfig{
				StorageAccountName: "testaccount",
				StorageAccountKey:  "testkey",
				TenantID:           "tenant-id",
				ClientID:           "client-id",
				ClientSecret:       "client-secret",
				ContainerName:      "testcontainer",
			},
			wantErr: true,
			errMsg:  "tenant_id cannot be set when using storage_account_key authentication",
		},
		{
			name: "tenant_id with storage_connection_string",
			config: AzureConfig{
				StorageAccountName:      "testaccount",
				StorageConnectionString: "connection-string",
				TenantID:                "tenant-id",
				ClientID:                "client-id",
				ClientSecret:            "client-secret",
				ContainerName:           "testcontainer",
			},
			wantErr: true,
			errMsg:  "tenant_id cannot be set when using storage_connection_string authentication",
		},
		{
			name: "tenant_id with user_assigned_id",
			config: AzureConfig{
				StorageAccountName: "testaccount",
				UserAssignedID:     "test-msi-id",
				TenantID:           "tenant-id",
				ClientID:           "client-id",
				ClientSecret:       "client-secret",
				ContainerName:      "testcontainer",
			},
			wantErr: true,
			errMsg:  "tenant_id cannot be set when using user_assigned_id authentication",
		},
		{
			name: "negative max_tries",
			config: AzureConfig{
				StorageAccountName: "testaccount",
				StorageAccountKey:  "testkey",
				ContainerName:      "testcontainer",
				PipelineConfig: PipelineConfig{
					MaxTries: -1,
				},
			},
			wantErr: true,
			errMsg:  "The value of max_tries must be greater than or equal to 0 in the config file",
		},
		{
			name: "negative max_retry_requests",
			config: AzureConfig{
				StorageAccountName: "testaccount",
				StorageAccountKey:  "testkey",
				ContainerName:      "testcontainer",
				ReaderConfig: ReaderConfig{
					MaxRetryRequests: -1,
				},
			},
			wantErr: true,
			errMsg:  "The value of max_retry_requests must be greater than or equal to 0 in the config file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AzureConfig.validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("AzureConfig.validate() error = %v, want error containing %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

// TestAzureConfig_validate_multipleErrors tests that multiple validation errors are combined
func TestAzureConfig_validate_multipleErrors(t *testing.T) {
	config := AzureConfig{
		// Missing StorageAccountName and ContainerName
		StorageAccountKey:       "testkey",
		StorageConnectionString: "connection-string",
		PipelineConfig: PipelineConfig{
			MaxTries: -1,
		},
	}

	err := config.validate()
	if err == nil {
		t.Fatal("AzureConfig.validate() expected error, got nil")
	}

	errMsg := err.Error()
	expectedErrors := []string{
		"storage_account_key and storage_connection_string cannot both be set",
		"storage_account_name is required but not configured",
		"no container specified",
		"The value of max_tries must be greater than or equal to 0 in the config file",
	}

	for _, expected := range expectedErrors {
		if !strings.Contains(errMsg, expected) {
			t.Errorf("AzureConfig.validate() error missing expected message: %v\nGot: %v", expected, errMsg)
		}
	}
}

// TestParseAzureConfig tests parsing full federated-store YAML configuration files
// This test validates YAML parsing and configuration validation without attempting actual Azure connections
func TestParseAzureConfig(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		validate func(*testing.T, AzureConfig)
	}{
		{
			name:     "client secret configuration",
			filename: "azure-client-secret.yaml",
			validate: func(t *testing.T, config AzureConfig) {
				if config.StorageAccountName != "myaccount" {
					t.Errorf("StorageAccountName = %v, want myaccount", config.StorageAccountName)
				}
				if config.ContainerName != "opencost-data" {
					t.Errorf("ContainerName = %v, want opencost-data", config.ContainerName)
				}
				if config.TenantID != "12345678-1234-1234-1234-123456789012" {
					t.Errorf("TenantID = %v, want 12345678-1234-1234-1234-123456789012", config.TenantID)
				}
				if config.ClientID != "87654321-4321-4321-4321-210987654321" {
					t.Errorf("ClientID = %v, want 87654321-4321-4321-4321-210987654321", config.ClientID)
				}
				if config.ClientSecret != "my-secret-value" {
					t.Errorf("ClientSecret = %v, want my-secret-value", config.ClientSecret)
				}
				if config.Endpoint != "blob.core.windows.net" {
					t.Errorf("Endpoint = %v, want blob.core.windows.net", config.Endpoint)
				}
			},
		},
		{
			name:     "shared key configuration",
			filename: "azure-shared-key.yaml",
			validate: func(t *testing.T, config AzureConfig) {
				if config.StorageAccountName != "myaccount" {
					t.Errorf("StorageAccountName = %v, want myaccount", config.StorageAccountName)
				}
				if config.StorageAccountKey != "base64encodedkey==" {
					t.Errorf("StorageAccountKey = %v, want base64encodedkey==", config.StorageAccountKey)
				}
				if config.ContainerName != "opencost-data" {
					t.Errorf("ContainerName = %v, want opencost-data", config.ContainerName)
				}
			},
		},
		{
			name:     "connection string configuration",
			filename: "azure-connection-string.yaml",
			validate: func(t *testing.T, config AzureConfig) {
				if config.StorageAccountName != "myaccount" {
					t.Errorf("StorageAccountName = %v, want myaccount", config.StorageAccountName)
				}
				if !strings.Contains(config.StorageConnectionString, "DefaultEndpointsProtocol=https") {
					t.Errorf("StorageConnectionString missing expected content")
				}
				if config.ContainerName != "opencost-data" {
					t.Errorf("ContainerName = %v, want opencost-data", config.ContainerName)
				}
			},
		},
		{
			name:     "managed identity configuration",
			filename: "azure-managed-identity.yaml",
			validate: func(t *testing.T, config AzureConfig) {
				if config.StorageAccountName != "myaccount" {
					t.Errorf("StorageAccountName = %v, want myaccount", config.StorageAccountName)
				}
				if config.UserAssignedID != "12345678-1234-1234-1234-123456789012" {
					t.Errorf("UserAssignedID = %v, want 12345678-1234-1234-1234-123456789012", config.UserAssignedID)
				}
				if config.ContainerName != "opencost-data" {
					t.Errorf("ContainerName = %v, want opencost-data", config.ContainerName)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Read the test file
			path := filepath.Join("testdata", tt.filename)
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("Failed to read test file %s: %v", path, err)
			}

			// Parse just the storage config structure to extract the Azure config
			var storageConfig StorageConfig
			if err := yaml.Unmarshal(data, &storageConfig); err != nil {
				t.Fatalf("Failed to parse storage config: %v", err)
			}

			// Marshal the config section back to YAML
			configData, err := yaml.Marshal(storageConfig.Config)
			if err != nil {
				t.Fatalf("Failed to marshal config: %v", err)
			}

			// Parse the Azure-specific configuration
			azureConfig, err := parseAzureConfig(configData)
			if err != nil {
				t.Fatalf("Failed to parse Azure config: %v", err)
			}

			// Validate the parsed configuration
			if tt.validate != nil {
				tt.validate(t, azureConfig)
			}

			// Ensure the configuration passes validation
			if err := azureConfig.validate(); err != nil {
				t.Errorf("Parsed config failed validation: %v", err)
			}
		})
	}
}
