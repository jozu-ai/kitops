package artifact

import (
	"bytes"
	"context"
	"fmt"

	_ "crypto/sha256"
	_ "crypto/sha512"

	"github.com/opencontainers/go-digest"
	"github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2/content/oci"
)

type Store struct {
	Storage *oci.Store
}

func NewArtifactStore() *Store {

	store, error := oci.New(".jozuStore")
	if error != nil {
		panic(error)
	}

	return &Store{
		Storage: store,
	}
}

func (store *Store) SaveContentLayer (layer *Layer)  (*v1.Descriptor, error)  {
	 ctx := context.Background()
	 
	 buf := &bytes.Buffer{}
	 err := layer.Apply(buf)
	 if err != nil {
		return nil, err
	}
	
    // Create a descriptor for the layer
    desc := v1.Descriptor{
        MediaType: v1.MediaTypeImageLayer,
        Digest:    digest.FromBytes(buf.Bytes()),
        Size:      int64(buf.Len()),
    }
	err = store.Storage.Push(ctx, desc, buf)
	if err != nil {
		return nil, err
	}
	fmt.Println("Saved model layer: ", desc.Digest)
	return &desc, nil
}

func (store *Store) SaveModelFile (model *JozuFile) error {
	ctx := context.Background()
	buf, err := model.MarshalToJSON()
	if err != nil {
		return err
	}
	desc := v1.Descriptor{
		MediaType: "application/vnd.jozu.model",
		Digest:    digest.FromBytes(buf),
		Size:      int64(len(buf)),
	}
	err = store.Storage.Push(ctx, desc, bytes.NewReader(buf))
	if err != nil {
		return err
	}
	return nil
}