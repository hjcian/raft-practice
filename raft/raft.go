package raft

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

type role int

const (
	minTimeout             = 150 // ms
	maxTimeout             = 300 // ms
	heartBeatInterval      = time.Millisecond * 100
	unknown           role = iota
	follower
	candidate
	leader
)

func randElectionTimeout() time.Duration {
	return time.Millisecond * time.Duration(rand.Intn(minTimeout)+(maxTimeout-minTimeout))
}

type NodeHandler interface {
	RequestVote(reqName string, reqTerm int) (int, error)
	AppendEntries(reqName string, reqTerm int) (bool, error)
}

func NewNode(name string) *Raft {
	return &Raft{
		name:      name,
		heartbeat: make(chan struct{}),
		role:      follower,
		term:      0,
		votedTerm: make(map[int]struct{}),
	}
}

type Raft struct {
	name          string
	heartbeat     chan struct{}
	otherNodes    []NodeHandler
	electionTimer *time.Timer

	rwlock    sync.RWMutex
	role      role
	term      int
	votedTerm map[int]struct{}
}

func (r *Raft) ReplyRequestVote(reqName string, reqTerm int) (int, error) {
	r.rwlock.Lock()
	defer r.rwlock.Unlock()
	if reqTerm < r.term {
		log.Println(r.name, "got expired request from:", reqName, reqTerm)
		return 0, nil
	}

	if _, ok := r.votedTerm[reqTerm]; ok {
		log.Println(r.name, "already voted", reqTerm)
		return 0, nil
	}

	if r.role == follower {
		r.term = reqTerm
		r.votedTerm[r.term] = struct{}{}
		r.resetElectionTimer()
		log.Println(r.name, "vote to", reqName, "at term:", reqTerm)
		return 1, nil
	}

	return 0, nil
}

func (r *Raft) ReplyAppendEntries(reqName string, reqTerm int) (bool, error) {
	r.rwlock.Lock()
	defer r.rwlock.Unlock()
	if reqTerm < r.term {
		return false, fmt.Errorf("%s term (%v) is greater than you (%v)", r.name, r.term, reqTerm)
	}

	fmt.Println(r.name, ": got heartbeat => turns out follower")
	r.role = follower
	r.heartbeat <- struct{}{}
	return true, nil
}

func (r *Raft) startElection() (int, error) {
	votes := 0
	for _, node := range r.otherNodes {
		vote, err := node.RequestVote(r.name, r.term)
		if err != nil {
			log.Println(r.name, ": node error:", err)
		}
		votes += vote
	}
	return votes, nil
}

func (r *Raft) startHeartbeat() {
	go func() {
		timer := time.NewTimer(heartBeatInterval)
		for r.role == leader {
			<-timer.C
			log.Println(r.name, ": send heartbeat to other nodes")
			for _, node := range r.otherNodes {
				_, err := node.AppendEntries(r.name, r.term)
				if err != nil {
					log.Println(r.name, ": send heartbeat failed:", err)
				}
			}
			timer.Reset(heartBeatInterval)
		}
	}()
}

func (r *Raft) resetElectionTimer() {
	if r.electionTimer == nil {
		r.electionTimer = time.NewTimer(randElectionTimeout())
	}

	if !r.electionTimer.Stop() {
		<-r.electionTimer.C
	}
	r.electionTimer.Reset(randElectionTimeout())
}

func (r *Raft) startup() {
	go func() {
		for {
			r.resetElectionTimer()

			select {
			case <-r.electionTimer.C:
				fmt.Println(r.name, ": timed out, start election")
				r.rwlock.Lock()
				r.role = candidate
				r.term++
				r.votedTerm[r.term] = struct{}{}
				r.rwlock.Unlock()

				votes, err := r.startElection()
				if err != nil {
					log.Println(r.name, ": start election failed: ", err)
				} else if votes > (len(r.otherNodes)+1)/2 {
					r.role = leader
					r.startHeartbeat()
				}
			case <-r.heartbeat:
			}

		}
	}()
}
