package ws_hub

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

//type mockConn struct{}
//
//func (c *mockConn) Close() error { return nil }

//func newTestClient(t *testing.T, hub *Hub) *Client {
//
//	return &Client{
//		UserID: uuid.New(),
//		hub:    nil,
//		conn:   nil,
//		send:   nil,
//	}
//}

func TestInMemoryClients_SetAndGet(t *testing.T) {
	t.Parallel()
	store := NewInMemoryClients()

	client1 := &Client{UserID: uuid.New()}
	client2 := &Client{UserID: client1.UserID}
	client3 := &Client{UserID: uuid.New()}

	err := store.SetConn(client1)
	require.NoError(t, err)
	err = store.SetConn(client2)
	require.NoError(t, err)
	err = store.SetConn(client3)
	require.NoError(t, err)

	conns1, err := store.GetConn(client1.UserID)
	require.NoError(t, err)
	require.Len(t, conns1, 2, "Waiting for 2 connections for User 1")
	require.Contains(t, conns1, client1)
	require.Contains(t, conns1, client2)

	conns3, err := store.GetConn(client3.UserID)
	require.NoError(t, err)
	require.Len(t, conns3, 1, "Waiting for 1 connections for User 3")
	require.Contains(t, conns3, client3)
}

func TestInMemoryClients_Concurrency(t *testing.T) {
	t.Parallel()
	store := NewInMemoryClients()
	wg := sync.WaitGroup{}
	numGoroutines := 100
	numOpsPerG := 50

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numOpsPerG; j++ {
				client := &Client{UserID: uuid.New()}

				_ = store.SetConn(client)

				_, _ = store.GetConn(client.UserID)

				_ = store.Delete(client)
			}
		}()
	}

	wg.Wait()
}
