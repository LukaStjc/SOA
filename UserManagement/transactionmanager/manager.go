// /transactionmanager/manager.go

package transactionmanager

import (
	"context"
	"sync"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"gorm.io/gorm"
)

type Manager struct {
	Transactions      map[uint32]*gorm.DB
	Neo4jTransactions map[uint32]neo4j.ExplicitTransaction
	lock              sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		Transactions:      make(map[uint32]*gorm.DB),
		Neo4jTransactions: make(map[uint32]neo4j.ExplicitTransaction),
	}
}

func (m *Manager) SavePendingTransaction(userID uint32, tx *gorm.DB, neo4j_tx neo4j.ExplicitTransaction) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.Transactions[userID] = tx
	m.Neo4jTransactions[userID] = neo4j_tx
}

func (m *Manager) Commit(userID uint32, ctx context.Context) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if tx, exists := m.Transactions[userID]; exists {
		tx.Commit()
		delete(m.Transactions, userID)
	}

	if tx, exists := m.Neo4jTransactions[userID]; exists {
		tx.Commit(ctx)
		delete(m.Neo4jTransactions, userID)
	}
}

func (m *Manager) Rollback(userID uint32, ctx context.Context) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if tx, exists := m.Transactions[userID]; exists {
		tx.Rollback()
		delete(m.Transactions, userID)
	}

	if tx, exists := m.Neo4jTransactions[userID]; exists {
		tx.Rollback(ctx)
		delete(m.Neo4jTransactions, userID)
	}
}
