package store

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

const path = "users.json"

type User struct {
	ID        int64  `json:"id"`
	ChatID    int64  `json:"chat_id"`
	Username  string `json:"username,omitempty"`
	Name      string `json:"name,omitempty"`
	FirstSeen string `json:"first_seen"`
}

var (
	mu     sync.Mutex
	users  map[int64]User
	loaded bool
)

func load() {
	if loaded {
		return
	}
	loaded = true
	users = map[int64]User{}

	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	var list []User
	if json.Unmarshal(data, &list) == nil {
		for _, u := range list {
			users[u.ID] = u
		}
	}
}

func save() {
	list := make([]User, 0, len(users))
	for _, u := range users {
		list = append(list, u)
	}
	if data, err := json.MarshalIndent(list, "", "  "); err == nil {
		os.WriteFile(path, data, 0644)
	}
}

func Remember(u User) {
	mu.Lock()
	defer mu.Unlock()
	load()

	if existing, ok := users[u.ID]; ok {
		u.FirstSeen = existing.FirstSeen
		if existing == u {
			return
		}
	} else {
		u.FirstSeen = time.Now().Format("2006-01-02 15:04")
	}
	users[u.ID] = u
	save()
}

func All() []User {
	mu.Lock()
	defer mu.Unlock()
	load()

	list := make([]User, 0, len(users))
	for _, u := range users {
		list = append(list, u)
	}
	return list
}

func Remove(id int64) {
	mu.Lock()
	defer mu.Unlock()
	load()

	if _, ok := users[id]; ok {
		delete(users, id)
		save()
	}
}

func Count() int {
	mu.Lock()
	defer mu.Unlock()
	load()
	return len(users)
}
