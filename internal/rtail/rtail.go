package rtail

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	kh "golang.org/x/crypto/ssh/knownhosts"
)

func tailFile(ctx context.Context, session *ssh.Session, file string, linec chan<- string) error {

	command := fmt.Sprintf("tail -F %s", file)

	var wg sync.WaitGroup
	errc := make(chan error, 3)
	consumeStream := func(r io.Reader) {
		logrus.Debugf("Starting consumeStream\n")
		defer logrus.Debugf("Done consumeStream\n")
		defer wg.Done()
		scan := bufio.NewScanner(r)
		scan.Split(bufio.ScanLines)
		for scan.Scan() {
			logrus.Tracef("Sending line to channel %s\n", scan.Text())
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

	wg.Add(2)
	go consumeStream(stderr)
	go consumeStream(stdout)
	logrus.Printf("Starting session with command: %s\n", command)
	if err := session.Start(command); err != nil {
		return fmt.Errorf("exec command: %s. %v", command, err)
	}
	go func() {
		if err := session.Wait(); err != nil {
			errc <- err
		}
	}()

	for {
		end := false
		select {
		case <-ctx.Done():
			logrus.Errorf("Context cancelled")
		case err := <-errc:
			if err != nil {
				logrus.Errorf("Error: %s", err)
			}
		}
		if end {
			break
		}
	}
	session.Close()
	wg.Wait()
	close(errc)
	return nil
}

type Options struct {
	// The name of the file to tail
	Filename string
	// username to use when connecting to the remote host
	User string
	// address of the remote host
	Address string
	// port to use when connecting to the remote host
	Port string
	// The private key to use when connecting to the remote host
	Key string
	// The known hosts file to use when connecting to the remote host
	KnownHosts string
}

// Tail tails the file and sends the lines to the lines channel
func Tail(ctx context.Context, opts Options, lines chan string) error {
	// Read the private key file.
	key, err := os.ReadFile(opts.Key)
	if err != nil {
		return err
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return err
	}

	hostKeyCallback, err := kh.New(opts.KnownHosts)
	if err != nil {
		return err
	}

	config := &ssh.ClientConfig{
		User: opts.User,
		Auth: []ssh.AuthMethod{
			// Add in password check here for moar security.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: hostKeyCallback,
	}
	// Connect to the remote server and perform the SSH handshake.
	client, err := ssh.Dial("tcp", opts.Address+":"+opts.Port, config)
	if err != nil {
		return err
	}
	defer client.Close()
	ss, err := client.NewSession()
	if err != nil {
		return err
	}
	defer ss.Close()

	err = tailFile(ctx, ss, opts.Filename, lines)
	if err != nil {
		log.Fatal("unable to tail file: ", err)
		return err
	}
	return nil
}
