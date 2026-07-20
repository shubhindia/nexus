package toolstate

import "sync"

var thoughtSignatures sync.Map

func StoreThoughtSignature(callID string, thoughtSignature string) {
	if callID == "" || thoughtSignature == "" {
		return
	}

	thoughtSignatures.Store(callID, thoughtSignature)
}

func LookupThoughtSignature(callID string) string {
	if callID == "" {
		return ""
	}

	value, ok := thoughtSignatures.Load(callID)
	if !ok {
		return ""
	}

	thoughtSignature, _ := value.(string)
	return thoughtSignature
}
