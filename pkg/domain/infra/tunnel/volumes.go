package tunnel

import (
	"context"

	"github.com/containers/podman/v2/pkg/bindings/volumes"
	"github.com/containers/podman/v2/pkg/domain/entities"
	"github.com/pkg/errors"
)

func (ic *ContainerEngine) VolumeCreate(ctx context.Context, opts entities.VolumeCreateOptions) (*entities.IDOrNameResponse, error) {
	response, err := volumes.Create(ic.ClientCxt, opts, nil)
	if err != nil {
		return nil, err
	}
	return &entities.IDOrNameResponse{IDOrName: response.Name}, nil
}

func (ic *ContainerEngine) VolumeRm(ctx context.Context, namesOrIds []string, opts entities.VolumeRmOptions) ([]*entities.VolumeRmReport, error) {
	if opts.All {
		vols, err := volumes.List(ic.ClientCxt, nil)
		if err != nil {
			return nil, err
		}
		for _, v := range vols {
			namesOrIds = append(namesOrIds, v.Name)
		}
	}
	reports := make([]*entities.VolumeRmReport, 0, len(namesOrIds))
	for _, id := range namesOrIds {
		options := new(volumes.RemoveOptions).WithForce(opts.Force)
		reports = append(reports, &entities.VolumeRmReport{
			Err: volumes.Remove(ic.ClientCxt, id, options),
			Id:  id,
		})
	}
	return reports, nil
}

func (ic *ContainerEngine) VolumeInspect(ctx context.Context, namesOrIds []string, opts entities.InspectOptions) ([]*entities.VolumeInspectReport, []error, error) {
	var (
		reports = make([]*entities.VolumeInspectReport, 0, len(namesOrIds))
		errs    = []error{}
	)
	if opts.All {
		vols, err := volumes.List(ic.ClientCxt, nil)
		if err != nil {
			return nil, nil, err
		}
		for _, v := range vols {
			namesOrIds = append(namesOrIds, v.Name)
		}
	}
	for _, id := range namesOrIds {
		data, err := volumes.Inspect(ic.ClientCxt, id, nil)
		if err != nil {
			errModel, ok := err.(entities.ErrorModel)
			if !ok {
				return nil, nil, err
			}
			if errModel.ResponseCode == 404 {
				errs = append(errs, errors.Errorf("no such volume %q", id))
				continue
			}
			return nil, nil, err
		}
		reports = append(reports, &entities.VolumeInspectReport{VolumeConfigResponse: data})
	}
	return reports, errs, nil
}

func (ic *ContainerEngine) VolumePrune(ctx context.Context, opts entities.VolumePruneOptions) ([]*entities.VolumePruneReport, error) {
	options := new(volumes.PruneOptions).WithFilters(opts.Filters)
	return volumes.Prune(ic.ClientCxt, options)
}

func (ic *ContainerEngine) VolumeList(ctx context.Context, opts entities.VolumeListOptions) ([]*entities.VolumeListReport, error) {
	options := new(volumes.ListOptions).WithFilters(opts.Filter)
	return volumes.List(ic.ClientCxt, options)
}
