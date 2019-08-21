package engine

import (
	"fmt"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	appcat_util "kmodules.xyz/custom-resources/client/clientset/versioned/typed/appcatalog/v1alpha1/util"
	api "kubevault.dev/operator/apis/engine/v1alpha1"
	"kubevault.dev/operator/pkg/vault"
)

func (secretEngineClient *SecretEngine) CreateConfig() error {
	vAppRef := &appcat.AppReference{
		Namespace: secretEngineClient.secretEngine.Namespace,
		Name:      secretEngineClient.secretEngine.Spec.VaultRef.Name,
	}

	// Update vault client so that it has the permission
	// to create config
	vClient, err2 := vault.NewClient(secretEngineClient.kubeClient, secretEngineClient.appClient, vAppRef)
	if err2 != nil {
		return errors.Wrap(err2, "failed to create vault api client")
	}
	secretEngineClient.vaultClient = vClient

	var err error
	engSpec := secretEngineClient.secretEngine.Spec
	if engSpec.GCP != nil {
		err = secretEngineClient.CreateGCPConfig()
	} else if engSpec.Azure != nil {
		err = secretEngineClient.CreateAzureConfig()
	} else if engSpec.AWS != nil {
		err = secretEngineClient.CreateAWSConfig()
	} else if engSpec.MySQL != nil {
		err = secretEngineClient.CreateMySQLConfig()
	} else if engSpec.Postgres != nil {
		err = secretEngineClient.CreatePostgresConfig()
	} else if engSpec.MongoDB != nil {
		err = secretEngineClient.CreateMongoDBConfig()
	} else {
		return errors.New("failed to create config: unknown secret engine type")
	}
	return err
}

// https://www.vaultproject.io/api/secret/databases/index.html#configure-connection
// https:https://www.vaultproject.io/api/secret/databases/mysql-maria.html#configure-connection
//
// CreateMySQLConfig creates MySQL database configuration
func (secretEngineClient *SecretEngine) CreateMySQLConfig() error {
	config := secretEngineClient.secretEngine.Spec.MySQL
	if config == nil {
		return errors.New("MySQL database config is nil")
	}

	// Set Default plugin name, if config.PluginName is empty
	config.SetDefaults()

	dbAppRef := config.DatabaseRef
	dbApp, err := secretEngineClient.appClient.AppBindings(dbAppRef.Namespace).Get(dbAppRef.Name, metav1.GetOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to get DatabaseAppBinding for MySQL database config")
	}

	connURL, err := dbApp.URLTemplate()
	if err != nil {
		return errors.Wrap(err, "failed to get MySQL database connection url")
	}

	path := fmt.Sprintf("/v1/%s/config/%s", secretEngineClient.path, api.GetDBNameFromAppBindingRef(&dbAppRef))
	req := secretEngineClient.vaultClient.NewRequest("POST", path)
	payload := map[string]interface{}{
		"plugin_name":    config.PluginName,
		"allowed_roles":  config.AllowedRoles,
		"connection_url": connURL,
	}

	if dbApp.Spec.Secret != nil {
		secret, err := secretEngineClient.kubeClient.CoreV1().Secrets(dbAppRef.Namespace).Get(dbApp.Spec.Secret.Name, metav1.GetOptions{})
		if err != nil {
			return errors.Wrap(err, "failed to get secret for MySQL database config")
		}

		data := make(map[string]interface{}, len(secret.Data))
		for k, v := range secret.Data {
			data[k] = string(v)
		}

		err = appcat_util.TransformCredentials(secretEngineClient.kubeClient, dbApp.Spec.SecretTransforms, data)
		if err != nil {
			return err
		}
		if v, ok := data[appcat.KeyUsername]; ok {
			payload[appcat.KeyUsername] = v
		}
		if v, ok := data[appcat.KeyPassword]; ok {
			payload[appcat.KeyPassword] = v
		}
	}

	if config.MaxOpenConnections > 0 {
		payload["max_open_connections"] = config.MaxOpenConnections
	}
	if config.MaxIdleConnections > 0 {
		payload["max_idle_connections"] = config.MaxIdleConnections
	}
	if config.MaxConnectionLifetime != "" {
		payload["max_connection_lifetime"] = config.MaxConnectionLifetime
	}

	err = req.SetJSONBody(payload)
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = secretEngineClient.vaultClient.RawRequest(req)
	return err
}

// https://www.vaultproject.io/api/secret/databases/index.html#configure-connection
// https://www.vaultproject.io/api/secret/databases/mongodb.html#configure-connection
//
// CreateMongoDBConfig creates MongoDB database configuration
func (secretEngineClient *SecretEngine) CreateMongoDBConfig() error {
	config := secretEngineClient.secretEngine.Spec.MongoDB
	if config == nil {
		return errors.New("MongoDB database config is nil")
	}

	// Set Default plugin name, if config.PluginName is empty
	config.SetDefaults()

	dbAppRef := config.DatabaseRef
	dbApp, err := secretEngineClient.appClient.AppBindings(dbAppRef.Namespace).Get(dbAppRef.Name, metav1.GetOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to get DatabaseAppBinding for MongoDB database config")
	}

	connURL, err := dbApp.URLTemplate()
	if err != nil {
		return errors.Wrap(err, "failed to get MongoDB database connection url")
	}

	path := fmt.Sprintf("/v1/%s/config/%s", secretEngineClient.path, api.GetDBNameFromAppBindingRef(&dbAppRef))
	req := secretEngineClient.vaultClient.NewRequest("POST", path)

	payload := map[string]interface{}{
		"plugin_name":    config.PluginName,
		"allowed_roles":  config.AllowedRoles,
		"connection_url": connURL,
	}

	if dbApp.Spec.Secret != nil {
		secret, err := secretEngineClient.kubeClient.CoreV1().Secrets(dbAppRef.Namespace).Get(dbApp.Spec.Secret.Name, metav1.GetOptions{})
		if err != nil {
			return errors.Wrap(err, "Failed to get secret for MongoDB database config")
		}

		data := make(map[string]interface{}, len(secret.Data))
		for k, v := range secret.Data {
			data[k] = string(v)
		}

		err = appcat_util.TransformCredentials(secretEngineClient.kubeClient, dbApp.Spec.SecretTransforms, data)
		if err != nil {
			return err
		}
		if v, ok := data[appcat.KeyUsername]; ok {
			payload[appcat.KeyUsername] = v
		}
		if v, ok := data[appcat.KeyPassword]; ok {
			payload[appcat.KeyPassword] = v
		}
	}

	if config.WriteConcern != "" {
		payload["write_concern"] = config.WriteConcern
	}

	err = req.SetJSONBody(payload)
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = secretEngineClient.vaultClient.RawRequest(req)
	if err != nil {
		return errors.Wrap(err, "failed to create database config")
	}
	return nil
}

// https://www.vaultproject.io/api/secret/databases/index.html#configure-connection
// https://www.vaultproject.io/api/secret/databases/postgresql.html#configure-connection
//
// CreatePostgresConfig creates database configuration
func (secretEngineClient *SecretEngine) CreatePostgresConfig() error {
	config := secretEngineClient.secretEngine.Spec.Postgres
	if config == nil {
		return errors.New("Postgres database config is nil")
	}

	// Set Default plugin name, if config.PluginName is empty
	config.SetDefaults()

	dbAppRef := config.DatabaseRef
	dbApp, err := secretEngineClient.appClient.AppBindings(dbAppRef.Namespace).Get(dbAppRef.Name, metav1.GetOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to get DatabaseAppBinding for Postgres database config")
	}

	connURL, err := dbApp.URLTemplate()
	if err != nil {
		return errors.Wrap(err, "failed to get Postgres database connection url")
	}

	path := fmt.Sprintf("/v1/%s/config/%s", secretEngineClient.path, api.GetDBNameFromAppBindingRef(&dbAppRef))
	req := secretEngineClient.vaultClient.NewRequest("POST", path)

	payload := map[string]interface{}{
		"plugin_name":    config.PluginName,
		"allowed_roles":  config.AllowedRoles,
		"connection_url": connURL,
	}

	if dbApp.Spec.Secret != nil {
		secret, err := secretEngineClient.kubeClient.CoreV1().Secrets(dbAppRef.Namespace).Get(dbApp.Spec.Secret.Name, metav1.GetOptions{})
		if err != nil {
			return errors.Wrap(err, "Failed to get secret for Postgres database config")
		}

		data := make(map[string]interface{}, len(secret.Data))
		for k, v := range secret.Data {
			data[k] = string(v)
		}

		err = appcat_util.TransformCredentials(secretEngineClient.kubeClient, dbApp.Spec.SecretTransforms, data)
		if err != nil {
			return err
		}
		if v, ok := data[appcat.KeyUsername]; ok {
			payload[appcat.KeyUsername] = v
		}
		if v, ok := data[appcat.KeyPassword]; ok {
			payload[appcat.KeyPassword] = v
		}
	}

	if config.MaxOpenConnections > 0 {
		payload["max_open_connections"] = config.MaxOpenConnections
	}
	if config.MaxIdleConnections > 0 {
		payload["max_idle_connections"] = config.MaxIdleConnections
	}
	if config.MaxConnectionLifetime != "" {
		payload["max_connection_lifetime"] = config.MaxConnectionLifetime
	}

	err = req.SetJSONBody(payload)
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = secretEngineClient.vaultClient.RawRequest(req)
	if err != nil {
		return errors.Wrap(err, "failed to create database config")
	}
	return nil
}

// ref:
// - https://www.vaultproject.io/api/secret/aws/index.html#configure-root-iam-credentials

// Configures AWS secret engine at specified path
func (secretEngineClient *SecretEngine) CreateAWSConfig() error {
	config := secretEngineClient.secretEngine.Spec.AWS
	if config == nil {
		return errors.New("AWS config is nil")
	}

	if secretEngineClient.vaultClient == nil {
		return errors.New("vault client is nil")
	}

	path := fmt.Sprintf("/v1/%s/config/root", secretEngineClient.path)
	req := secretEngineClient.vaultClient.NewRequest("POST", path)

	payload := map[string]interface{}{}
	if config.MaxRetries != nil {
		payload["max_retries"] = *config.MaxRetries
	}
	if config.Region != "" {
		payload["region"] = config.Region
	}
	if config.IAMEndpoint != "" {
		payload["iam_endpoint"] = config.IAMEndpoint
	}
	if config.STSEndpoint != "" {
		payload["sts_endpoint"] = config.STSEndpoint
	}

	if config.CredentialSecret != "" {
		sr, err := secretEngineClient.kubeClient.CoreV1().Secrets(secretEngineClient.secretEngine.Namespace).Get(config.CredentialSecret, metav1.GetOptions{})
		if err != nil {
			return errors.Wrap(err, "failed to get aws credential secret")
		}

		if val, ok := sr.Data[api.AWSCredentialAccessKeyKey]; ok {
			payload["access_key"] = string(val)
		}
		if val, ok := sr.Data[api.AWSCredentialSecretKeyKey]; ok {
			payload["secret_key"] = string(val)
		}
	}

	if err := req.SetJSONBody(payload); err != nil {
		return errors.Wrap(err, "failed to load payload in config create request")
	}

	_, err := secretEngineClient.vaultClient.RawRequest(req)
	if err != nil {
		return errors.Wrap(err, "failed to create aws config")
	}

	// set lease config
	if config.LeaseConfig != nil {
		path := fmt.Sprintf("/v1/%s/config/lease", secretEngineClient.path)
		req := secretEngineClient.vaultClient.NewRequest("POST", path)

		payload := map[string]interface{}{
			"lease":     config.LeaseConfig.Lease,
			"lease_max": config.LeaseConfig.LeaseMax,
		}
		if err := req.SetJSONBody(payload); err != nil {
			return errors.Wrap(err, "failed to load payload in create lease config request")
		}

		_, err := secretEngineClient.vaultClient.RawRequest(req)
		if err != nil {
			return errors.Wrap(err, "failed to create aws lease config")
		}
	}
	return nil
}

// ref:
//	- https://www.vaultproject.io/api/secret/azure/index.html#configure-access

// Configures Azure secret engine at specified path
func (secretEngineClient *SecretEngine) CreateAzureConfig() error {
	config := secretEngineClient.secretEngine.Spec.Azure
	if config == nil {
		return errors.New("Azure config is nil")
	}

	if secretEngineClient.vaultClient == nil {
		return errors.New("vault client is nil")
	}

	path := fmt.Sprintf("/v1/%s/config", secretEngineClient.path)
	req := secretEngineClient.vaultClient.NewRequest("POST", path)

	payload := map[string]interface{}{}
	if config.CredentialSecret != "" {
		sr, err := secretEngineClient.kubeClient.CoreV1().Secrets(secretEngineClient.secretEngine.Namespace).Get(config.CredentialSecret, metav1.GetOptions{})
		if err != nil {
			return errors.Wrap(err, "failed to get azure credential secret")
		}

		if val, ok := sr.Data[api.AzureSubscriptionID]; ok && len(val) > 0 {
			payload["subscription_id"] = string(val)
		} else {
			return errors.New("azure secret engine configuration failed: subscription id missing")
		}

		if val, ok := sr.Data[api.AzureTenantID]; ok && len(val) > 0 {
			payload["tenant_id"] = string(val)
		} else {
			return errors.New("azure secret engine configuration failed: tenant id missing")
		}

		if val, ok := sr.Data[api.AzureClientID]; ok && len(val) > 0 {
			payload["client_id"] = string(val)
		}

		if val, ok := sr.Data[api.AzureClientSecret]; ok && len(val) > 0 {
			payload["client_secret"] = string(val)
		}
	}

	if config.Environment != "" {
		payload["environment"] = config.Environment
	}

	if err := req.SetJSONBody(payload); err != nil {
		return errors.Wrap(err, "failed to load payload in config create request")
	}

	_, err := secretEngineClient.vaultClient.RawRequest(req)
	if err != nil {
		return errors.Wrap(err, "failed to create azure config")
	}
	return nil
}

// ref:
//  - https://www.vaultproject.io/api/secret/gcp/index.html#write-config

// Configures GCP secret engine at specified path
func (secretEngineClient *SecretEngine) CreateGCPConfig() error {
	config := secretEngineClient.secretEngine.Spec.GCP
	if config == nil {
		return errors.New("GCP config is nil")
	}

	if secretEngineClient.vaultClient == nil {
		return errors.New("vault client is nil")
	}

	path := fmt.Sprintf("/v1/%s/config", secretEngineClient.path)
	req := secretEngineClient.vaultClient.NewRequest("POST", path)

	payload := map[string]interface{}{}
	if config.TTL != "" {
		payload["ttl"] = config.TTL
	}
	if config.MaxTTL != "" {
		payload["max_ttl"] = config.MaxTTL
	}

	if config.CredentialSecret != "" {
		sr, err := secretEngineClient.kubeClient.CoreV1().Secrets(secretEngineClient.secretEngine.Namespace).Get(config.CredentialSecret, metav1.GetOptions{})
		if err != nil {
			return errors.Wrap(err, "failed to get gcp credential secret")
		}

		if val, ok := sr.Data[api.GCPSACredentialJson]; ok {
			payload["credentials"] = string(val)
		}
	}

	if err := req.SetJSONBody(payload); err != nil {
		return errors.Wrap(err, "failed to load payload in config create request")
	}

	_, err := secretEngineClient.vaultClient.RawRequest(req)
	if err != nil {
		return errors.Wrap(err, "failed to create gcp config")
	}
	return nil
}
