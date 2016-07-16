package concurrent

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

const tagKey = "mu"

type ErrRegistrationError struct {
	Name string
}

func (e ErrRegistrationError) Error() string {
	return fmt.Sprintf("mutex already registered: %s", e.Name)
}

type ErrInvalidName struct {
	Name string
}

func (e ErrInvalidName) Error() string {
	return fmt.Sprintf("mutex not registered: %s", e.Name)
}

type ErrInvalidObject struct {
	IsNil bool
}

func (e ErrInvalidObject) Error() string {
	if e.IsNil {
		return fmt.Sprintf("object is nil")
	}
	return fmt.Sprintf("object is invalid")
}

type MutexGroup struct {
	mu sync.Mutex
	m  map[string]*sync.Mutex
}

// Register registers a mutex per given name.
// If any provided name is already registered, no mutexes are registered and an error is returned.
func (mutexGroup *MutexGroup) Register(names ...string) error {
	mutexGroup.mu.Lock()
	defer mutexGroup.mu.Unlock()
	if mutexGroup.m == nil {
		mutexGroup.m = make(map[string]*sync.Mutex)
	}
	// only allow any registration if ALL given names are free
	for _, n := range names {
		_, has := mutexGroup.m[n]
		if has {
			return ErrRegistrationError{Name: n}
		}
	}
	// init mutexes for each given name
	for _, n := range names {
		mutexGroup.m[n] = &sync.Mutex{}
	}
	return nil
}

// RegisterObjs registers and initializes mutexes for the given objects.
// This method panics if the provided arguments are not pointers to structs.
// If mutex names are duplicated within the provided objects, an error is returned
// and no mutexes will be initialized on the objects and registered in this MutexGroup.
func (mutexGroup *MutexGroup) RegisterObjs(objs ...interface{}) error {
	mutexGroup.mu.Lock()
	defer mutexGroup.mu.Unlock()

	if mutexGroup.m == nil {
		mutexGroup.m = make(map[string]*sync.Mutex)
	}

	names := []string{}

	// iterate over each object and register all mutexes
	// which are defined via field tag "mu:mutex_name"
	for _, obj := range objs {
		// retrieve type and reflected value of the given object
		ele := reflect.ValueOf(obj).Elem()
		ty := ele.Type()

		// iterate over all fields of the given object
		for i := 0; i < ty.NumField(); i++ {
			field := ty.FieldByIndex([]int{i})
			mutexName := strings.TrimSpace(field.Tag.Get(tagKey))
			if len(mutexName) == 0 {
				continue
			}
			names = append(names, mutexName)
		}
	}

	// only allow any registration if ALL given names are free
	for _, n := range names {
		_, has := mutexGroup.m[n]
		if has {
			return ErrRegistrationError{Name: n}
		}
	}

	// TODO: find a better way to init everything without having to iterate
	// twice over the given objects while still only initializing
	// if no duplicated names are provided.

	// now reflect again all provided objects and initialize the mutexes.
	// also register the mutexes inside the MutexGroup
	for _, obj := range objs {
		// retrieve type and reflected value of the given object
		ele := reflect.ValueOf(obj).Elem()
		ty := ele.Type()

		// iterate over all fields of the given object
		for i := 0; i < ty.NumField(); i++ {
			field := ty.FieldByIndex([]int{i})
			mutexName := strings.TrimSpace(field.Tag.Get(tagKey))
			if len(mutexName) == 0 {
				continue
			}
			m := &sync.Mutex{}
			ele.Field(i).Set(reflect.ValueOf(m))
			mutexGroup.m[mutexName] = m
		}
	}
	return nil
}

// Registered returns a string slice of all registered mutexes.
func (mutexGroup *MutexGroup) Registered() []string {
	mutexGroup.mu.Lock()
	defer mutexGroup.mu.Unlock()
	names := []string{}
	for n := range mutexGroup.m {
		names = append(names, n)
	}
	return names
}

// Unregister unregisters the given mutexes.
// If any provided name is not registered, no mutexes are unregistered and an error is returned.
func (mutexGroup *MutexGroup) Unregister(names ...string) error {
	mutexGroup.mu.Lock()
	defer mutexGroup.mu.Unlock()
	// only allow any un-registration if ALL names are registered
	for _, n := range names {
		_, has := mutexGroup.m[n]
		if has {
			return ErrInvalidName{Name: n}
		}
	}
	for _, n := range names {
		delete(mutexGroup.m, n)
	}
	return nil
}

// Lock locks the given mutexes.
// If a given name is not registered, an error is returned while possibly still having locked some mutexes.
func (mutexGroup *MutexGroup) Lock(names ...string) error {
	mutexGroup.mu.Lock()
	defer mutexGroup.mu.Unlock()
	for _, n := range names {
		m, has := mutexGroup.m[n]
		if !has {
			return ErrInvalidName{Name: n}
		}
		m.Lock()
	}
	return nil
}

// Unlock unlocks the given mutexes.
// If a given name is not registered, an error is returned while possibly still having unlocked some mutexes.
func (mutexGroup *MutexGroup) Unlock(names ...string) error {
	mutexGroup.mu.Lock()
	defer mutexGroup.mu.Unlock()
	for _, n := range names {
		m, has := mutexGroup.m[n]
		if !has {
			return ErrInvalidName{Name: n}
		}
		m.Unlock()
	}
	return nil
}
