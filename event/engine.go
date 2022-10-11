package event

import (
	"context"
	"fmt"
)

var type2Engine map[string]iEngine
var validEngines []iEngine

type iEngine interface {
	GetType() string
	Emit(ctx context.Context, e *Event)
}

func registerEngine(eg iEngine) {
	if type2Engine == nil {
		type2Engine = make(map[string]iEngine)
	}
	if _, ok := type2Engine[eg.GetType()]; !ok {
		type2Engine[eg.GetType()] = eg
		validEngines = append(validEngines, eg)
		fmt.Printf("event engine: %s registered", eg.GetType())
	}
}

func getValidEngines() (ies []iEngine) {
	return validEngines
}
