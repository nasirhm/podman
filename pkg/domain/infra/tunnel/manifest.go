package tunnel

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/containers/podman/v2/pkg/bindings/manifests"
	"github.com/containers/podman/v2/pkg/domain/entities"
	"github.com/pkg/errors"
)

// ManifestCreate implements manifest create via ImageEngine
func (ir *ImageEngine) ManifestCreate(ctx context.Context, names, images []string, opts entities.ManifestCreateOptions) (string, error) {
	options := new(manifests.CreateOptions).WithAll(opts.All)
	imageID, err := manifests.Create(ir.ClientCxt, names, images, options)
	if err != nil {
		return imageID, errors.Wrapf(err, "error creating manifest")
	}
	return imageID, err
}

// ManifestInspect returns contents of manifest list with given name
func (ir *ImageEngine) ManifestInspect(ctx context.Context, name string) ([]byte, error) {
	list, err := manifests.Inspect(ir.ClientCxt, name, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting content of manifest list or image %s", name)
	}

	buf, err := json.MarshalIndent(list, "", "    ")
	if err != nil {
		return buf, errors.Wrapf(err, "error rendering manifest for display")
	}
	return buf, err
}

// ManifestAdd adds images to the manifest list
func (ir *ImageEngine) ManifestAdd(ctx context.Context, opts entities.ManifestAddOptions) (string, error) {
	options := new(manifests.AddOptions).WithAll(opts.All).WithArch(opts.Arch).WithVariant(opts.Variant)
	options.WithFeatures(opts.Features).WithImages(opts.Images).WithOS(opts.OS).WithOSVersion(opts.OSVersion)
	if len(opts.Annotation) != 0 {
		annotations := make(map[string]string)
		for _, annotationSpec := range opts.Annotation {
			spec := strings.SplitN(annotationSpec, "=", 2)
			if len(spec) != 2 {
				return "", errors.Errorf("no value given for annotation %q", spec[0])
			}
			annotations[spec[0]] = spec[1]
		}
		options.WithAnnotation(annotations)
	}

	listID, err := manifests.Add(ir.ClientCxt, opts.Images[1], options)
	if err != nil {
		return listID, errors.Wrapf(err, "error adding to manifest list %s", opts.Images[1])
	}
	return listID, nil
}

// ManifestAnnotate updates an entry of the manifest list
func (ir *ImageEngine) ManifestAnnotate(ctx context.Context, names []string, opts entities.ManifestAnnotateOptions) (string, error) {
	return "", errors.New("not implemented")
}

// ManifestRemove removes the digest from manifest list
func (ir *ImageEngine) ManifestRemove(ctx context.Context, names []string) (string, error) {
	updatedListID, err := manifests.Remove(ir.ClientCxt, names[0], names[1], nil)
	if err != nil {
		return updatedListID, errors.Wrapf(err, "error removing from manifest %s", names[0])
	}
	return fmt.Sprintf("%s :%s\n", updatedListID, names[1]), nil
}

// ManifestPush pushes a manifest list or image index to the destination
func (ir *ImageEngine) ManifestPush(ctx context.Context, name, destination string, opts entities.ManifestPushOptions) error {
	options := new(manifests.PushOptions).WithAll(opts.All)
	_, err := manifests.Push(ir.ClientCxt, name, destination, options)
	return err
}
