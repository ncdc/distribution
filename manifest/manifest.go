package manifest

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/docker/distribution/digest"
)

// Versioned provides a struct with just the manifest schemaVersion. Incoming
// content with unknown schema version can be decoded against this struct to
// check the version.
type Versioned struct {
	// SchemaVersion is the image manifest schema that this image follows
	SchemaVersion int `json:"schemaVersion"`
}

// Manifest provides the base accessible fields for working with V2 image
// format in the registry.
type Manifest struct {
	Versioned

	// Name is the name of the image's repository
	Name string `json:"name"`

	// Tag is the tag of the image specified by this manifest
	Tag string `json:"tag"`

	// Architecture is the host architecture on which this image is intended to
	// run
	Architecture string `json:"architecture"`

	// FSLayers is a list of filesystem layer blobSums contained in this image
	FSLayers []FSLayer `json:"fsLayers"`

	// History is a list of unstructured historical data for v1 compatibility
	History []History `json:"history"`
}

func toMap(in interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ToMap only accepts structs; got %T", v)
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)

		out[fi.Name] = v.Field(i).Interface()

	}
	return out, nil
}

func (m *Manifest) Digest() (digest.Digest, error) {
	manifestMap, err := toMap(m)
	if err != nil {
		return "", err
	}

	b, err := json.Marshal(manifestMap)
	if err != nil {
		return "", err
	}

	return digest.FromBytes(b)
}

// SignedManifest provides an envelope for a signed image manifest, including
// the format sensitive raw bytes. It contains fields to
type SignedManifest struct {
	Manifest

	// Raw is the byte representation of the ImageManifest, used for signature
	// verification. The value of Raw must be used directly during
	// serialization, or the signature check will fail. The manifest byte
	// representation cannot change or it will have to be re-signed.
	Raw []byte `json:"-"`
}

// UnmarshalJSON populates a new ImageManifest struct from JSON data.
func (sm *SignedManifest) UnmarshalJSON(b []byte) error {
	var manifest Manifest
	if err := json.Unmarshal(b, &manifest); err != nil {
		return err
	}

	sm.Manifest = manifest
	sm.Raw = make([]byte, len(b), len(b))
	copy(sm.Raw, b)

	return nil
}

// MarshalJSON returns the contents of raw. If Raw is nil, marshals the inner
// contents. Applications requiring a marshaled signed manifest should simply
// use Raw directly, since the the content produced by json.Marshal will be
// compacted and will fail signature checks.
func (sm *SignedManifest) MarshalJSON() ([]byte, error) {
	if len(sm.Raw) > 0 {
		return sm.Raw, nil
	}

	// If the raw data is not available, just dump the inner content.
	return json.Marshal(&sm.Manifest)
}

// FSLayer is a container struct for BlobSums defined in an image manifest
type FSLayer struct {
	// BlobSum is the tarsum of the referenced filesystem image layer
	BlobSum digest.Digest `json:"blobSum"`
}

// History stores unstructured v1 compatibility information
type History struct {
	// V1Compatibility is the raw v1 compatibility information
	V1Compatibility string `json:"v1Compatibility"`
}
