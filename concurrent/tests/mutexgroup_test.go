package tests

import (
	"sync"
	"testing"

	"github.com/luca-moser/belt/concurrent"
)

var (
	mg1 = concurrent.MutexGroup{}
	mg2 = concurrent.MutexGroup{}

	validObj   = &obj{}
	invalidObj = &invobj{}
)

type (
	obj struct {
		A *sync.Mutex `mu:"mutex_a"`
		B *sync.Mutex `mu:"mutex_b"`
	}

	invobj struct {
		A *sync.Mutex `mu:"mutex_a"`
		B *sync.Mutex `mu:"mutex_a"` // whoops!
	}
)

func TestRegister(t *testing.T) {
	if err := mg1.Register("alice", "vanessa", "james"); err != nil {
		t.Fatal(err)
	}
}

func TestLocking(t *testing.T) {
	if err := mg1.Lock("alice"); err != nil {
		t.Fatal(err)
	}
}

func TestUnlocking(t *testing.T) {
	if err := mg1.Unlock("alice"); err != nil {
		t.Fatal(err)
	}
}

func TestLockingMultiple(t *testing.T) {
	if err := mg1.Lock("alice", "vanessa", "james"); err != nil {
		t.Fatal(err)
	}
}

func TestUnlockingMultiple(t *testing.T) {
	if err := mg1.Unlock("alice", "vanessa", "james"); err != nil {
		t.Fatal(err)
	}
}

func TestRegisterObjs(t *testing.T) {
	if err := mg2.RegisterObjs(validObj); err != nil {
		t.Fatal(err)
	}
	if len(mg2.Registered()) != 2 {
		t.Fatal("registered mutexes are not of length 2")
	}

	if validObj.A == nil || validObj.B == nil {
		t.Fatal("mutex was not initialized on object")
	}
}

func TestLockingObjs(t *testing.T) {
	if err := mg2.Lock("mutex_a", "mutex_b"); err != nil {
		t.Fatal(err)
	}
}

func TestUnlokcingObjs(t *testing.T) {
	if err := mg2.Unlock("mutex_a", "mutex_b"); err != nil {
		t.Fatal(err)
	}
}

func TestRegisterObjsInvalid(t *testing.T) {
	if err := mg2.RegisterObjs(invalidObj); err == nil {
		t.Fatal("no error was returned but there were duplicated mutex names")
	}
}
