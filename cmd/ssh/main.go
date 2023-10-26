package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"sync"

	"golang.org/x/crypto/ssh"
	kh "golang.org/x/crypto/ssh/knownhosts"
)

func tailFile(session *ssh.Session, file string, linec chan<- string) error {

	command := fmt.Sprintf("tail -F %s", file)

	var wg sync.WaitGroup
	errc := make(chan error, 3)
	consumeStream := func(r io.Reader) {
		log.Printf("Starting consumeStream\n")
		defer wg.Done()
		scan := bufio.NewScanner(r)
		scan.Split(bufio.ScanLines)
		for scan.Scan() {
			log.Printf("Sending line to channel %s\n", scan.Text())
			linec <- scan.Text()
		}
		if err := scan.Err(); err != nil {
			errc <- err
		}
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		return fmt.Errorf("opening stderr: %v", err)
	}
	stdout, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("opening stdout: %v", err)
	}

	wg.Add(1)
	go consumeStream(stderr)
	go consumeStream(stdout)
	log.Printf("Starting session with command: %s\n", command)
	if err := session.Start(command); err != nil {
		return err
	}
	wg.Add(1)
	go func() {
		if err := session.Wait(); err != nil {
			errc <- err
		}
	}()

	go func() {
		wg.Wait()
		close(errc)
	}()

	return <-errc
}

func main() {
	user := "root"
	address := "10.0.4.38"
	port := "22"

	key, err := ioutil.ReadFile("/Users/eliofrancesconi/.ssh/id_rsa")
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	hostKeyCallback, err := kh.New("/Users/eliofrancesconi/.ssh/known_hosts")
	if err != nil {
		log.Fatal("could not create hostkeycallback function: ", err)
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			// Add in password check here for moar security.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: hostKeyCallback,
	}
	// Connect to the remote server and perform the SSH handshake.
	client, err := ssh.Dial("tcp", address+":"+port, config)
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}
	defer client.Close()
	ss, err := client.NewSession()
	if err != nil {
		log.Fatal("unable to create SSH session: ", err)
	}
	defer ss.Close()

	response := make(chan string)
	go func() {
		err = tailFile(ss, "/var/log/containers/nxw-sv__avo.log", response)
		if err != nil {
			log.Fatal("unable to tail file: ", err)
			return
		}
	}()
	for {
		select {
		case line := <-response:
			fmt.Println(line)
		}
	}
}
