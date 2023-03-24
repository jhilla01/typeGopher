package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"sync"

	"github.com/creack/pty"
	"golang.org/x/crypto/ssh"
)

// Winsize stores terminal dimensions
type Winsize struct {
	Cols uint16
	Rows uint16
}

// main Sets up an SSH server, loads the private key, listens for connections, and processes them.
func main() {

	// Configure server and load private key
	// You can generate a keypair with 'ssh-keygen -t rsa'
	config := &ssh.ServerConfig{NoClientAuth: true}
	privateBytes, err := os.ReadFile("id_rsa")
	if err != nil {
		log.Fatal("Failed to load private key (./id_rsa)")
	}
	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal("Failed to parse private key")
	}
	config.AddHostKey(private)

	// Start listening for connections
	listener, err := net.Listen("tcp", "0.0.0.0:2200")
	if err != nil {
		log.Fatalf("Failed to listen on 2200 (%s)", err)
	}
	log.Print("Listening on 2200...")

	// Accept and handle connections
	for {
		tcpConn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept incoming connection (%s)", err)
			continue
		}
		// Before use, a handshake must be performed on the incoming net.Conn.
		sshConn, chans, reqs, err := ssh.NewServerConn(tcpConn, config)
		if err != nil {
			log.Printf("Failed to handshake (%s)", err)
			continue
		}

		log.Printf("New SSH connection from %s (%s)", sshConn.RemoteAddr(), sshConn.ClientVersion())
		// Discard all global out-of-band Requests
		go ssh.DiscardRequests(reqs)
		// Accept all channels
		go handleChannels(chans)
	}
}

// handleChannels Services incoming SSH channels in a separate goroutine.
func handleChannels(chans <-chan ssh.NewChannel) {
	// Service the incoming Channel channel in go routine
	for newChannel := range chans {
		go handleChannel(newChannel)
	}
}

// handleChannel Handles a single SSH channel, starts a shell session, and processes requests.
func handleChannel(newChannel ssh.NewChannel) {
	//  Expect a channel type of "session" with a shell
	if t := newChannel.ChannelType(); t != "session" {
		newChannel.Reject(ssh.UnknownChannelType, fmt.Sprintf("unknown channel type: %s", t))
		return
	}

	// Accept or reject client connection
	connection, requests, err := newChannel.Accept()
	if err != nil {
		log.Printf("Could not accept channel (%s)", err)
		return
	}

	// Fire up bash for this session
	bash := exec.Command("./main")

	// Prepare teardown function
	c := func() {
		connection.Close()
		_, err := bash.Process.Wait()
		if err != nil {
			log.Printf("Failed to exit bash (%s)", err)
		}
		log.Printf("Session closed")
	}

	// Allocate a terminal for this channel
	log.Print("Creating pty...")
	bashf, err := pty.Start(bash)
	if err != nil {
		log.Printf("Could not start pty (%s)", err)
		c()
		return
	}

	// Pipe session connection to bash and bash to connection
	var once sync.Once
	go func() {
		io.Copy(connection, bashf)
		once.Do(c)
	}()
	go func() {
		io.Copy(bashf, connection)
		once.Do(c)
	}()

	// Sessions have out-of-band requests such as "shell", "pty-req" and "env"
	// Provides a conceptually independent channel that allows for any data to be sent
	go func() {
		for req := range requests {
			switch req.Type {
			case "shell":
				// We only accept the default shell
				// (i.e. no command in the Payload)
				if len(req.Payload) == 0 {
					req.Reply(true, nil)
				}
			case "pty-req":
				termLen := req.Payload[3]
				w, h := parseDims(req.Payload[termLen+4:])
				SetWinsize(bashf.Fd(), w, h)
				// Responding true (OK) here will let the client
				// know we have a pty ready for input
				req.Reply(true, nil)
			case "window-change":
				w, h := parseDims(req.Payload)
				SetWinsize(bashf.Fd(), w, h)
			}
		}
	}()
}

// parseDims extracts terminal dimensions (width x height) from the provided buffer.
func parseDims(b []byte) (uint32, uint32) {
	w := binary.BigEndian.Uint32(b)
	h := binary.BigEndian.Uint32(b[4:])
	return w, h
}

// SetWinsize sets the size of the given pseudotty.
func SetWinsize(fd uintptr, cols, rows uint32) error {
	ws := &Winsize{
		Cols: uint16(cols),
		Rows: uint16(rows),
	}
	ptyWs := &pty.Winsize{
		Cols: ws.Cols,
		Rows: ws.Rows,
	}
	return pty.Setsize(os.NewFile(fd, "ptyFd"), ptyWs)
}
