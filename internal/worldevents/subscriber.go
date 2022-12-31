package worldevents

var WorldEventSubscribers = make(map[string][]string)

// GetWorldEventSubscribers get the subscribed users for feast and siege
func GetWorldEventSubscribers() map[string][]string {
	m.RLock()
	defer m.RUnlock()
	return WorldEventSubscribers
}

// SetWorldEventSubscribers set the subscribed users for feast and siege
func SetWorldEventSubscribers(worldEventType string, userID string) {
	v := GetWorldEventSubscribers()[worldEventType]
	var exists = false
	for _, u := range v {
		if u == userID {
			exists = true
		}
	}
	if !exists {
		m.Lock()
		WorldEventSubscribers[worldEventType] = append(WorldEventSubscribers[worldEventType], userID)
		m.Unlock()
	}

}

// UnsetWorldEventSubscribers unregister users for feast and siege
func UnsetWorldEventSubscribers(worldEventType string, userID string) {
	// WorldEventSubscribers[worldEventType] = append(WorldEventSubscribers[worldEventType], userID)
	if v, ok := GetWorldEventSubscribers()[worldEventType]; ok {
		var found = false
		var index = 0
		for i, u := range v {
			if u == userID {
				found = true
				index = i
			}
		}
		if found {
			v[index] = v[len(v)-1] // Copy last element to index i.
			v[len(v)-1] = ""
			m.Lock()                                             // Erase last element (write zero value).
			WorldEventSubscribers[worldEventType] = v[:len(v)-1] // Truncate slice.
			m.Unlock()
		}

	}

}
