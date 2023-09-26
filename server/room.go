package main

import "sync"

type Room struct {
	name    string
	clients []*Client
	size    int
	m       sync.Mutex
}

func NewRoom(name string) *Room {
	return &Room{
		name:    name,
		clients: make([]*Client, 0, MAX_ROOM_SIZE),
		size:    0,
		m:       sync.Mutex{},
	}
}

func (r *Room) Join(client *Client) bool {
	r.m.Lock()
	defer r.m.Unlock()

	if r.size >= MAX_ROOM_SIZE {
		return false
	}

	r.clients = append(r.clients, client)
	r.size++

	return true
}

func (r *Room) Leave(client *Client) {
	r.m.Lock()
	defer r.m.Unlock()

	for i, c := range r.clients {
		if c == client {
			r.clients = append(r.clients[:i], r.clients[i+1:]...)
			r.size--
			return
		}
	}
}

func (r *Room) Broadcast(msg string) {
	for _, c := range r.clients {
		if c == nil {
			continue
		}

		_, _ = c.Write([]byte(msg))
	}
}

func (r *Room) Destroy() {
	r.m.Lock()
	defer r.m.Unlock()

	for _, c := range r.clients {
		if c == nil {
			continue
		}

		_ = c.Close()
	}
}
