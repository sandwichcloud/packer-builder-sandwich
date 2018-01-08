package sandwich

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/packer/packer"
	"github.com/mitchellh/multistep"
	"golang.org/x/crypto/ssh"
)

type StepKeyPair struct {
	Debug                bool
	SSHAgentAuth         bool
	TemporaryKeyPairName string
	KeyPairID            string
	PrivateKeyFile       string

	doCleanup          bool
	temporaryKeyPairID string
}

func (s *StepKeyPair) Run(state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)

	if s.PrivateKeyFile != "" {
		privateKeyBytes, err := ioutil.ReadFile(s.PrivateKeyFile)
		if err != nil {
			state.Put("error", fmt.Errorf("Error loading configured private key file: %s", err))
			return multistep.ActionHalt
		}

		state.Put("keyPairID", s.KeyPairID)
		state.Put("privateKey", string(privateKeyBytes))
		return multistep.ActionContinue
	}

	if s.SSHAgentAuth && s.KeyPairID != "" {
		ui.Say(fmt.Sprintf("Using SSH Agent for existing key pair %s", s.KeyPairID))
		state.Put("keyPairID", s.KeyPairID)
		return multistep.ActionContinue
	}

	if s.TemporaryKeyPairName == "" {
		ui.Say("Not using temporary keypair")
		state.Put("keyPairID", "")
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

	keypairClient := config.sandwichClient.Keypair()
	keypair, err := keypairClient.Create(s.TemporaryKeyPairName, string(ssh.MarshalAuthorizedKey(pubKey)))
	if err != nil {
		state.Put("error", fmt.Errorf("Error creating temporary keypair: %s", err))
		return multistep.ActionHalt
	}

	s.temporaryKeyPairID = keypair.ID.String()
	ui.Say(fmt.Sprintf("Created temporary keypair: %s (%s)", s.TemporaryKeyPairName, s.temporaryKeyPairID))

	s.doCleanup = true
	state.Put("keyPairID", s.temporaryKeyPairID)
	state.Put("privateKey", privateKeyString)

	return multistep.ActionContinue
}

func (s *StepKeyPair) Cleanup(state multistep.StateBag) {
	if !s.doCleanup {
		return
	}

	config := state.Get("config").(Config)
	ui := state.Get("ui").(packer.Ui)

	keypairClient := config.sandwichClient.Keypair()

	ui.Say(fmt.Sprintf("Deleting temporary keypair: %s (%s)", s.TemporaryKeyPairName, s.temporaryKeyPairID))
	err := keypairClient.Delete(s.temporaryKeyPairID)
	if err != nil {
		ui.Error(fmt.Sprintf("Error cleaning up keypair. Please delete the key manually: %s", s.TemporaryKeyPairName))
	}
}
