package postgres

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestLock(t *testing.T) {
	t.Parallel()

	testDB := NewTestDatabase(t)
	ctx := context.Background()

	const (
		id1 = "test1"
		id2 = "test2"
	)

	mustLock := func(id string, ttl time.Duration) UnlockFn {
		t.Helper()
		unlock, err := testDB.Lock(ctx, id, ttl)
		if err != nil {
			t.Fatal(err)
		}
		return unlock
	}

	//  Grab a free lock.
	unlock1 := mustLock(id1, time.Hour)

	// Fail to grab a held lock.
	if _, err := testDB.Lock(ctx, id1, time.Hour); !errors.Is(err, ErrAlreadyLocked) {
		t.Fatalf("got %v, wanted ErrAlreadyLocked", err)
	}
	unlock2 := mustLock(id2, time.Hour)
	// Unlock the first lock.
	if err := unlock1(); err != nil {
		t.Fatal(err)
	}

	// Re-acquire the first lock, briefly.
	_ = mustLock(id1, time.Microsecond)

	// We can acquire the lock after it expires.
	time.Sleep(50 * time.Millisecond)
	unlock1 = mustLock(id1, time.Hour)

	// Unlock both locks.
	if err := unlock1(); err != nil {
		t.Fatal(err)
	}
	if err := unlock2(); err != nil {
		t.Fatal(err)
	}

	// Lock table should be empty.
	conn, err := testDB.Pool.Acquire(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Release()
	var count int
	if err := conn.QueryRow(ctx, `SELECT count(*) FROM Lock`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("got %d rows from Lock table, wanted zero", count)
	}
}
