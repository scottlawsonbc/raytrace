// // Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.
// // Automatically handle JSON marshalling and unmarshalling of interface types.
package phys

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sync"
)

// Type registry for interface types in a scene.
var (
	typeRegistry  = make(map[string]reflect.Type)
	registryMutex sync.RWMutex
)

// getInterfaceType retrieves a registered interface type by its name.
func getInterfaceType(typeName string) (reflect.Type, bool) {
	registryMutex.RLock()
	defer registryMutex.RUnlock()
	t, exists := typeRegistry[typeName]
	return t, exists
}

// RegisterInterfaceType registers a type with its name.
// This function should be called from each init function.
func RegisterInterfaceType(v any) {
	typ := reflect.TypeOf(v)
	var name string
	// Check if typ is a pointer and get the element type
	if typ.Kind() == reflect.Ptr {
		name = typ.Elem().Name()
		typ = typ.Elem()
	} else {
		name = typ.Name()
	}
	if name == "" {
		panic("cannot register a type with no name")
	}
	registryMutex.Lock()
	defer registryMutex.Unlock()
	if _, exists := typeRegistry[name]; exists {
		panic(fmt.Sprintf("Type '%s' is already registered", name))
	}
	typeRegistry[name] = typ
	log.Printf("RegisterInterfaceType: %s", name)
}

// Marshal an interface{} by wrapping it with its type name.
func marshalInterface(typ any) (json.RawMessage, error) {
	if typ == nil {
		return nil, fmt.Errorf("cannot marshal a nil value")
	}
	// Get the concrete type
	typType := reflect.TypeOf(typ)
	var typeName string
	if typType.Kind() == reflect.Ptr {
		// Get the name of the element type if it's a pointer
		typeName = typType.Elem().Name()
	} else {
		typeName = typType.Name()
	}
	if typeName == "" {
		return nil, fmt.Errorf("cannot marshal a type with no name: %T", typ)
	}
	// Marshal the interface's data.
	data, err := json.Marshal(typ)
	if err != nil {
		return nil, err
	}
	// Wrap with Type information.
	wrapped := map[string]interface{}{
		"Type": typeName,
		"Data": json.RawMessage(data),
	}
	return json.Marshal(wrapped)
}

// Unmarshal an interface{} by unwrapping it from its type name.
func unmarshalInterface(data json.RawMessage) (any, error) {
	var wrapper struct {
		Type string          `json:"Type"`
		Data json.RawMessage `json:"Data"`
	}
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, err
	}
	t, exists := getInterfaceType(wrapper.Type)
	if !exists {
		return nil, fmt.Errorf("unsupported type: `%s`; has it been registered?", wrapper.Type)
	}
	// Create a new instance of the type (pointer to t)
	s := reflect.New(t).Interface()
	if err := json.Unmarshal(wrapper.Data, s); err != nil {
		return nil, err
	}
	return s, nil
}
