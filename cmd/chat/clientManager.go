package main

import (
	"fmt"
	"log"
	"net"
	"sort"
	"strings"

	"github.com/nmmh/chat/utils"
)

type readOp struct {
	key  net.Conn
	resp chan *ClientState
}

type readOpAllVals struct {
	resp chan []string
}
type writeOp struct {
	key  net.Conn
	val  *ClientState
	resp chan bool
}

type broadcastMsgOp struct {
	msg *message
}

// ClientManager ygu
type ClientManager struct {
	clients          map[net.Conn]*ClientState
	reads            chan *readOp
	readsAllVals     chan *readOpAllVals
	writes           chan *writeOp
	msgsForBroadcast chan *broadcastMsgOp
	kills            chan net.Conn
}

//NewCM - creates a ClientManger
func NewCM() *ClientManager {
	return &ClientManager{
		clients:          make(map[net.Conn]*ClientState),
		reads:            make(chan *readOp),
		readsAllVals:     make(chan *readOpAllVals),
		writes:           make(chan *writeOp),
		msgsForBroadcast: make(chan *broadcastMsgOp),
		kills:            make(chan net.Conn),
	}
}

// Start jhg
func (cm *ClientManager) Start() {
	for {
		select {
		case read := <-cm.reads:
			read.resp <- cm.clients[read.key]
		case readAllVals := <-cm.readsAllVals:
			s := make([]string, 0)
			for _, val := range cm.clients {
				s = append(s, val.username)
			}
			readAllVals.resp <- s
		case write := <-cm.writes:
			cm.clients[write.key] = write.val
			write.resp <- true
		case msgForBroadcast := <-cm.msgsForBroadcast:
			// Loop over all connected clients
			for conn := range cm.clients {
				if msgForBroadcast.msg.msgScope == "ALLEXCEPTSENDER" {
					if msgForBroadcast.msg.username == cm.clients[conn].username {
						continue
					}
				} else if msgForBroadcast.msg.msgScope == "SENDERONLY" {
					if msgForBroadcast.msg.username != cm.clients[conn].username {
						continue
					}
				}
				//send msg in  a goroutine
				go sendMessage(conn, msgForBroadcast.msg)
			}
			//message always logged at the server
			log.Printf("%s", msgForBroadcast.msg.text)
		case kill := <-cm.kills:
			msgChannel <- &message{cm.clients[kill].username, "CHANOP", "ALL", fmt.Sprintf(" * [%s] disconnected\r\n", cm.clients[kill].username)}
			delete(cm.clients, kill)
			//log.Printf(getUserList(clientsMap))
			kill.Close()
		}
	}
}

// ReadByKey - this reads a value (clientState) from the CM
func (cm *ClientManager) ReadByKey(conn net.Conn) *ClientState {
	read := &readOp{key: conn, resp: make(chan *ClientState)}
	cm.reads <- read
	return <-read.resp
}

// ReadAll This should probably return a copy of the entire client manager
func (cm *ClientManager) ReadAll() []string {
	readAllVals := &readOpAllVals{resp: make(chan []string)}
	cm.readsAllVals <- readAllVals
	return <-readAllVals.resp
}

// Write the connection and the client state to the ClientManager
func (cm *ClientManager) Write(conn net.Conn, username string) bool {
	write := &writeOp{key: conn, val: &ClientState{username: username}, resp: make(chan bool)}
	cm.writes <- write
	return <-write.resp
}

//FormatUserList extracts,sorts adn returns a userlist string
func (cm *ClientManager) FormatUserList(usernames []string) (string, error) {
	//sort first
	sort.Strings(usernames)
	ul := "UserList:{"
	for _, username := range usernames {
		ul += fmt.Sprintf("%s, ", username)
	}
	ul = strings.TrimSuffix(ul, ", ") + fmt.Sprintf("} Total:[%d]", len(usernames))
	return ul, nil
}

//UsernameInUse looksup a username returns true if found
func (cm *ClientManager) UsernameInUse(usernames []string, search string) (bool, error) {
	return utils.StringInSlice(usernames, search)
}
