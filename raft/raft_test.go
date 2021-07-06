package raft

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type server struct {
	node *Raft
}

func (s *server) RequestVote(reqName string, reqTerm int) (int, error) {
	return s.node.ReplyRequestVote(reqName, reqTerm)
}

func (s *server) AppendEntries(reqName string, reqTerm int) (bool, error) {
	return s.node.ReplyAppendEntries(reqName, reqTerm)
}

func Test_Election(t *testing.T) {
	a := NewNode("A")
	b := NewNode("B")
	c := NewNode("C")

	sa := &server{a}
	sb := &server{b}
	sc := &server{c}

	a.otherNodes = append(a.otherNodes, sb, sc)
	b.otherNodes = append(b.otherNodes, sa, sc)
	c.otherNodes = append(c.otherNodes, sa, sb)

	a.startup()
	b.startup()
	c.startup()

	time.Sleep(time.Second)

	leaderCnt, candidateCnt, followerCnt, unknownCnt := 0, 0, 0, 0
	for _, role := range []role{a.role, b.role, c.role} {
		switch role {
		case leader:
			leaderCnt++
		case candidate:
			candidateCnt++
		case follower:
			followerCnt++
		default:
			unknownCnt++
		}
	}

	assert.Equal(t, 1, leaderCnt)
	assert.Equal(t, 2, followerCnt)
	assert.Equal(t, 0, candidateCnt)
	assert.Equal(t, 0, unknownCnt)
}
