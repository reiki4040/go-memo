package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
)

func main() {
	c := &PortForwardingConfig{}
	err := LoadFromYamlFile("./config.yml", c)
	if err != nil {
		log.Fatalf("can not load config file: %v", err)
	}
	c.DoBinding()

	// waiting Ctrl + C
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}

func LoadFromYamlFile(filePath string, p interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	yml, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(yml, p)
}

type PortForwardingConfig struct {
	Name                 string `yaml:"name"`
	SSHBastionHostPort   string `yaml:"ssh_bastion_host_port"`
	SSHUser              string `yaml:"ssh_user"`
	SSHKeyFilePath       string `yaml:"ssh_key_file_path"`
	LocalBindPort        string `yaml:"local_bind_port"`
	ForwardingRemotePort string `yaml:"forwarding_remote_port"`
}

func (c *PortForwardingConfig) DoBinding() {
	key, err := ioutil.ReadFile(c.SSHKeyFilePath)
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	// Create sshClientConfig
	sshConfig := &ssh.ClientConfig{
		User: c.SSHUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		//HostKeyCallback: ssh.FixedHostKey(hostKey),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Setup localListener (type net.Listener)
	localListener, err := net.Listen("tcp", c.LocalBindPort)
	if err != nil {
		fmt.Println("net.Listen failed: %v", err)
	} else {
		// go accept loop
		go func() {
			for {
				// Setup localConn (type net.Conn)
				localConn, err := localListener.Accept()
				if err != nil {
					fmt.Println("listen.Accept failed: %v", err)
					// maybe reconnection.
				}

				// go forwarding
				go forward(localConn, c.SSHBastionHostPort, c.ForwardingRemotePort, sshConfig)
			}
		}()
	}
}

func forward(localConn net.Conn, hostport, remoteport string, config *ssh.ClientConfig) {
	// Setup sshClientConn (type *ssh.ClientConn)
	sshClientConn, err := ssh.Dial("tcp", hostport, config)
	if err != nil {
		fmt.Println("ssh.Dial failed: %s", err)
	}
	// Setup sshConn (type net.Conn)
	sshConn, err := sshClientConn.Dial("tcp", remoteport)
	// Copy localConn.Reader to sshConn.Writer
	go func() {
		_, err = io.Copy(sshConn, localConn)
		if err != nil {
			fmt.Println("io.Copy failed: %v", err)
		}
	}()
	// Copy sshConn.Reader to localConn.Writer
	go func() {
		_, err = io.Copy(localConn, sshConn)
		if err != nil {
			fmt.Println("io.Copy failed: %v", err)
		}
	}()
}
