package raft

//
// this is an outline of the API that raft must expose to
// the service (or tester). see comments below for
// each of these functions for more details.
//
// rf = Make(...)
//   create a new Raft server.
// rf.Start(command interface{}) (index, term, isleader)
//   start agreement on a new log entry
// rf.GetState() (term, isLeader)
//   ask a Raft for its current term, and whether it thinks it is leader
// ApplyMsg
//   each time a new entry is committed to the log, each Raft peer
//   should send an ApplyMsg to the service (or tester)
//   in the same server.
//

import (
//	"bytes"
	"sync"
	"sync/atomic"

//	"6.824/labgob"
	"6.824/labrpc"
	"math/rand"
	"time"
	"fmt"
)


//
// as each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make(). set
// CommandValid to true to indicate that the ApplyMsg contains a newly
// committed log entry.
//
// in part 2D you'll want to send other kinds of messages (e.g.,
// snapshots) on the applyCh, but set CommandValid to false for these
// other uses.
//
type ApplyMsg struct {
	CommandValid bool
	Command      interface{}
	CommandIndex int

	// For 2D:
	SnapshotValid bool
	Snapshot      []byte
	SnapshotTerm  int
	SnapshotIndex int
}

type Log struct {
	Term		int
	index		int
}

const LEADER int = 0
const FOLLOWER int = 1
const CANDIDATE int = 1
//
// A Go object implementing a single Raft peer.
//
type Raft struct {
	mu        sync.Mutex          // Lock to protect shared access to this peer's state
	peers     []*labrpc.ClientEnd // RPC end points of all peers
	persister *Persister          // Object to hold this peer's persisted state
	me        int                 // this peer's index into peers[]
	dead      int32               // set by Kill()

	// Your data here (2A, 2B, 2C).
	// Look at the paper's Figure 2 for a description of what
	// state a Raft server must maintain.
	// persistent
	cur_term   int
	vote_for   int
	log		   []Log

	// volatile
	commit_idx int
	last_apply int

	next_idx   []int
	match_idx  []int

	heartsbeats bool

	state	   int
	wait_time	int64
}

// return currentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (int, bool) {

	var term int
	var isleader bool
	// Your code here (2A).

	rf.mu.Lock()
	defer rf.mu.Unlock()
	term = rf.cur_term
	isleader = (rf.state == LEADER)

	return term, isleader
}

//
// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
//
func (rf *Raft) persist() {
	// Your code here (2C).
	// Example:
	// w := new(bytes.Buffer)
	// e := labgob.NewEncoder(w)
	// e.Encode(rf.xxx)
	// e.Encode(rf.yyy)
	// data := w.Bytes()
	// rf.persister.SaveRaftState(data)
}


//
// restore previously persisted state.
//
func (rf *Raft) readPersist(data []byte) {
	if data == nil || len(data) < 1 { // bootstrap without any state?
		return
	}
	// Your code here (2C).
	// Example:
	// r := bytes.NewBuffer(data)
	// d := labgob.NewDecoder(r)
	// var xxx
	// var yyy
	// if d.Decode(&xxx) != nil ||
	//    d.Decode(&yyy) != nil {
	//   error...
	// } else {
	//   rf.xxx = xxx
	//   rf.yyy = yyy
	// }
}


//
// A service wants to switch to snapshot.  Only do so if Raft hasn't
// have more recent info since it communicate the snapshot on applyCh.
//
func (rf *Raft) CondInstallSnapshot(lastIncludedTerm int, lastIncludedIndex int, snapshot []byte) bool {

	// Your code here (2D).

	return true
}

// the service says it has created a snapshot that has
// all info up to and including index. this means the
// service no longer needs the log through (and including)
// that index. Raft should now trim its log as much as possible.
func (rf *Raft) Snapshot(index int, snapshot []byte) {
	// Your code here (2D).

}


//
// example RequestVote RPC arguments structure.
// field names must start with capital letters!
//
type RequestVoteArgs struct {
	// Your data here (2A, 2B).
	Term			int
	Candidate_idx	int
	Last_log_idx	int
	Last_log_term	int
}

//
// example RequestVote RPC reply structure.
// field names must start with capital letters!
//
type RequestVoteReply struct {
	// Your data here (2A).
	Vote_granted	bool
	Term			int
}

type RequestAppendArgs struct {
	Term			int
	Leader_id		int
	Prev_log_idx	int
	Prev_log_term	int

	Entries			[]Log
	leader_commit	int
}

type RequestAppendReply struct {
	Term			int
	Success			bool
}

//
// example RequestVote RPC handler.
//
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	// Your code here (2A, 2B).
	rf.mu.Lock()
	defer rf.mu.Unlock()
	if args.Term < rf.cur_term {
		reply.Vote_granted = false
	} else if args.Term == rf.cur_term {
		rf.state = FOLLOWER
		if rf.vote_for == -1 || rf.vote_for == rf.me {
			if args.Last_log_idx >= rf.log[len(rf.log)-1].index && args.Last_log_term >= rf.log[len(rf.log)-1].Term {
				reply.Vote_granted = true
				rf.vote_for = args.Candidate_idx
			} else {
				reply.Vote_granted = false;
			}
		} else {
			reply.Vote_granted = false
		}
	} else {
		rf.state = FOLLOWER
		rf.vote_for = args.Candidate_idx
		rf.cur_term = args.Term
		reply.Vote_granted = true
	}
	reply.Term = rf.cur_term
}

//
// example RequestVote RPC handler.
//
func (rf *Raft) AppendEntries(args *RequestAppendArgs, reply *RequestAppendReply) {
	// Your code here (2A, 2B).
	rf.mu.Lock()
	defer rf.mu.Unlock()

	rf.heartsbeats = true
	if len(args.Entries) == 0 {
		rf.vote_for = -1
		return
	}

	if args.Term < rf.cur_term {
		reply.Success = false
	} else {
		if args.Prev_log_idx >= len(rf.log) || rf.log[args.Prev_log_idx].Term != args.Prev_log_term {
			reply.Success = false
		} else {

		}
	}

	reply.Term = rf.cur_term
}

//
// example code to send a RequestVote RPC to a server.
// server is the index of the target server in rf.peers[].
// expects RPC arguments in args.
// fills in *reply with RPC reply, so caller should
// pass &reply.
// the types of the args and reply passed to Call() must be
// the same as the types of the arguments declared in the
// handler function (including whether they are pointers).
//
// The labrpc package simulates a lossy network, in which servers
// may be unreachable, and in which requests and replies may be lost.
// Call() sends a request and waits for a reply. If a reply arrives
// within a timeout interval, Call() returns true; otherwise
// Call() returns false. Thus Call() may not return for a while.
// A false return can be caused by a dead server, a live server that
// can't be reached, a lost request, or a lost reply.
//
// Call() is guaranteed to return (perhaps after a delay) *except* if the
// handler function on the server side does not return.  Thus there
// is no need to implement your own timeouts around Call().
//
// look at the comments in ../labrpc/labrpc.go for more details.
//
// if you're having trouble getting RPC to work, check that you've
// capitalized all field names in structs passed over RPC, and
// that the caller passes the address of the reply struct with &, not
// the struct itself.
//
func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	return ok
}

func (rf *Raft) sendRequestAppend(server int, args *RequestAppendArgs, reply *RequestAppendReply) bool {
	ok := rf.peers[server].Call("Raft.AppendEntries", args, reply)
	return ok
}

//
// the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. if this
// server isn't the leader, returns false. otherwise start the
// agreement and return immediately. there is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election. even if the Raft instance has been killed,
// this function should return gracefully.
//
// the first return value is the index that the command will appear at
// if it's ever committed. the second return value is the current
// term. the third return value is true if this server believes it is
// the leader.
//
func (rf *Raft) Start(command interface{}) (int, int, bool) {
	index := -1
	term := -1
	isLeader := true

	// Your code here (2B).


	return index, term, isLeader
}

//
// the tester doesn't halt goroutines created by Raft after each test,
// but it does call the Kill() method. your code can use killed() to
// check whether Kill() has been called. the use of atomic avoids the
// need for a lock.
//
// the issue is that long-running goroutines use memory and may chew
// up CPU time, perhaps causing later tests to fail and generating
// confusing debug output. any goroutine with a long-running loop
// should call killed() to check whether it should stop.
//
func (rf *Raft) Kill() {
	atomic.StoreInt32(&rf.dead, 1)
	// Your code here, if desired.
}

func (rf *Raft) killed() bool {
	z := atomic.LoadInt32(&rf.dead)
	return z == 1
}

func (rf *Raft) isCandidate() bool {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	return rf.state == CANDIDATE
}

func (rf *Raft) calculate_vote(msg_queue chan RequestVoteReply) {
	agree := 0
	n := len(rf.peers)

	rf.mu.Lock()
	defer rf.mu.Unlock()
	if rf.state != CANDIDATE {
		return
	}
	
	for msg := range msg_queue {
		if msg.Vote_granted {
			agree++
		}
	}

	fmt.Printf("Get Agree: %d\n", agree)
	if agree >= n/2+1 {
		rf.state = LEADER
		fmt.Printf("Node %d becomde Leader, term %d\n", rf.me, rf.cur_term)
	}
}

func (rf *Raft) elect() {
	rf.mu.Lock()
	fmt.Printf("Node %d Electing, term %d\n", rf.me, rf.cur_term)
	rf.vote_for = rf.me
	req := RequestVoteArgs{}
	req.Candidate_idx = rf.me
	req.Term = rf.cur_term
	req.Last_log_idx = rf.log[len(rf.log)-1].index
	req.Last_log_term = rf.log[len(rf.log)-1].Term
	rf.mu.Unlock()

	request := make([]RequestVoteArgs, len(rf.peers))
	reply := make([]RequestVoteReply, len(rf.peers))
	msg_queue := make(chan RequestVoteReply, len(rf.peers))

	for idx, _ := range rf.peers {
		if idx == rf.me {
			continue
		}
		request[idx] = req
		go func(idx int) {
			rf.sendRequestVote(idx, &request[idx], &reply[idx])
			msg_queue <- reply[idx]
		}(idx)
	}

	time.AfterFunc(200 * time.Millisecond, func() {
		rf.calculate_vote(msg_queue)
	})
	time.Sleep(200* time.Millisecond)
}

// The ticker go routine starts a new election if this peer hasn't received
// heartsbeats recently.
func (rf *Raft) ticker() {
	for rf.killed() == false {

		// Your code here to check if a leader election should
		// be started and to randomize sleeping time using
		// time.Sleep().
		time.Sleep(time.Duration(rf.wait_time) * time.Millisecond)

		rf.mu.Lock()
		if rf.heartsbeats == false {
			rf.state = CANDIDATE
			rf.cur_term++
			rf.persist()
		}
		rf.mu.Unlock()

		if rf.isCandidate() {
			rf.elect()
		}
	}
}

//
// the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
//
func Make(peers []*labrpc.ClientEnd, me int,
	persister *Persister, applyCh chan ApplyMsg) *Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.persister = persister
	rf.me = me

	// Your initialization code here (2A, 2B, 2C).

	// persistent
	rf.cur_term = -1
	rf.vote_for = -1

	// volatile
	rf.commit_idx = -1
	rf.last_apply = -1

	rf.next_idx = make([]int, len(peers))
	rf.match_idx = make([]int, len(peers))

	rf.heartsbeats = false
	rf.log = append(rf.log, Log{Term: 0, index: 0})

	rf.wait_time = 200 + (rand.Int63() % 400);
	rf.state = CANDIDATE

	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())

	// start ticker goroutine to start elections
	go rf.ticker()


	return rf
}
