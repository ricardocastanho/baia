package database

import (
	"context"
	"fmt"
	"sync"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Neo4jClient is a class that manages the Neo4j driver instance
type Neo4jClient struct {
	uri      string
	username string
	password string
	driver   neo4j.DriverWithContext
	mutex    sync.Mutex // Ensures safe concurrent access
}

// NewNeo4jClient is a constructor that creates a new Neo4jClient instance
func NewNeo4jClient(uri, username, password string) *Neo4jClient {
	return &Neo4jClient{
		uri:      uri,
		username: username,
		password: password,
	}
}

// GetDriver returns a Neo4j driver instance, creating it if necessary
func (client *Neo4jClient) GetDriver() (neo4j.DriverWithContext, error) {
	client.mutex.Lock()
	defer client.mutex.Unlock()

	// If the driver is already created, return it
	if client.driver != nil {
		return client.driver, nil
	}

	// Create a new Neo4j driver
	driver, err := neo4j.NewDriverWithContext(client.uri, neo4j.BasicAuth(client.username, client.password, ""))
	if err != nil {
		return nil, fmt.Errorf("failed to create Neo4j driver: %w", err)
	}

	// Assign the driver to the client and return it
	client.driver = driver
	return driver, nil
}

// Close closes the Neo4j driver instance
func (client *Neo4jClient) Close() error {
	client.mutex.Lock()
	defer client.mutex.Unlock()

	if client.driver != nil {
		err := client.driver.Close(context.Background())
		if err != nil {
			return fmt.Errorf("failed to close Neo4j driver: %w", err)
		}
		client.driver = nil
	}
	return nil
}
