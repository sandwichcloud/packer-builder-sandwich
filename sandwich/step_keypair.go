package sandwich

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"golang.org/x/crypto/ssh"
)

type StepKeyPair struct {
	Debug                bool
	SSHAgentAuth         bool
	TemporaryKeyPairName string
	KeyPairName          string
	PrivateKeyFile       string

	doCleanup            bool
	temporaryKeyPairName string
}

func (s *StepKeyPair) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)

	if s.PrivateKeyFile != "" {
		privateKeyBytes, err := ioutil.ReadFile(s.PrivateKeyFile)
		if err != nil {
			state.Put("error", fmt.Errorf("Error loading configured private key file: %s", err))
			return multistep.ActionHalt
		}

		state.Put("keyPairName", s.KeyPairName)
		state.Put("privateKey", string(privateKeyBytes))
		return multistep.ActionContinue
	}

	if s.SSHAgentAuth && s.KeyPairName != "" {
		ui.Say(fmt.Sprintf("Using SSH Agent for existing key pair %s", s.KeyPairName))
		state.Put("keyPairName", s.KeyPairName)
		return multistep.ActionContinue
	}

	if s.TemporaryKeyPairName == "" {
		ui.Say("Not using temporary keypair")
		state.Put("keyPairName", "")
		return multistep.ActionContinue
	}

	config := state.Get("config").(Config)
	ui.Say(fmt.Sprintf("Creating temporary keypair: %s ...", s.TemporaryKeyPairName))

	privateKey, err := rsa.GenerateKey(rand.Reader, 2046)
	if err != nil {
		state.Put("error", fmt.Errorf("Error creating generating private key for temporary keypair: %s", err))
		return multistep.ActionHalt
	}
	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	privateKeyString := string(pem.EncodeToMemory(privateKeyPEM))
	pubKey, err := ssh.NewPublicKey(&privateKey.PublicKey)

	keypairClient := config.sandwichClient.Keypair(config.ProjectName)
	keypair, err := keypairClient.Create(s.TemporaryKeyPairName, string(ssh.MarshalAuthorizedKey(pubKey)))
	if err != nil {
		state.Put("error", fmt.Errorf("Error creating temporary keypair: %s", err))
		return multistep.ActionHalt
	}

	s.temporaryKeyPairName = keypair.Name
	ui.Say(fmt.Sprintf("Created temporary keypair: %s (%s)", s.TemporaryKeyPairName, s.temporaryKeyPairName))

	s.doCleanup = true
	state.Put("keyPairName", s.temporaryKeyPairName)
	state.Put("privateKey", privateKeyString)

	return multistep.ActionContinue
}

func (s *StepKeyPair) Cleanup(state multistep.StateBag) {
	if !s.doCleanup {
		return
	}

	config := state.Get("config").(Config)
	ui := state.Get("ui").(packer.Ui)

	keypairClient := config.sandwichClient.Keypair(config.ProjectName)

	ui.Say(fmt.Sprintf("Deleting temporary keypair: %s (%s)", s.TemporaryKeyPairName, s.temporaryKeyPairName))
	err := keypairClient.Delete(s.temporaryKeyPairName)
	if err != nil {
		ui.Error(fmt.Sprintf("Error cleaning up keypair. Please delete the key manually: %s", s.TemporaryKeyPairName))
	}
}
