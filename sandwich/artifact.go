package sandwich

import (
	"fmt"
	"log"

	"github.com/sandwichcloud/deli-cli/api/client"
)

type Artifact struct {
	ImageName      string
	BuilderIdValue string

	ImageClient client.ImageClientInterface
}

func (a *Artifact) BuilderId() string {
	return a.BuilderIdValue
}

func (a *Artifact) Files() []string {
	return nil
}

func (a *Artifact) Id() string {
	return a.ImageName
}

func (a *Artifact) String() string {
	return fmt.Sprintf("An image was created: %v", a.ImageName)
}

func (a *Artifact) State(name string) interface{} {
	return nil
}

func (a *Artifact) Destroy() error {
	log.Printf("Destroying image: %s", a.ImageName)
	return a.ImageClient.Delete(a.ImageName)
}
