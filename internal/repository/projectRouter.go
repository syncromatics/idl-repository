package repository

import (
	"errors"
	"fmt"
)

type projectRouter struct {
	storage Storage
}

func newProjectRouter(storage Storage) *projectRouter {
	return &projectRouter{storage}
}

func (r *projectRouter) Register(router Muxer) {
	router.RegisterJson("/v1/projects", r.listHandler)
	router.RegisterJson("/v1/projects/{project:.*}/types", r.listTypeHandler)
	router.RegisterJson("/v1/projects/{project:.*}/types/{type:.*}/versions", r.listVersionHandler)
	router.RegisterData("/v1/projects/{project:.*}/types/{type:.*}/versions/{version:.*}/data.tar.gz", r.pullVersion)
	router.RegisterJson("/v1/projects/{project:.*}/types/{type:.*}/versions/{version:.*}", r.submitVersion)
}

func (r *projectRouter) listHandler(ctx HttpContext) (*JsonResponse, error) {
	projects, err := r.storage.ListFolders("/projects")
	if err != nil {
		return nil, err
	}

	return &JsonResponse{
		StatusCode: 200,
		Model:      projects,
	}, nil
}

func (r *projectRouter) listTypeHandler(ctx HttpContext) (*JsonResponse, error) {
	project, ok := ctx.Args["project"]
	if !ok {
		return nil, errors.New("failed to get project from args")
	}

	pth := fmt.Sprintf("/projects/%s", project)

	ok = r.storage.Exists(pth)
	if !ok {
		return &JsonResponse{
			StatusCode: 404,
			Model:      fmt.Sprintf("project '%s' does not exist", project),
		}, nil
	}

	types, err := r.storage.ListFolders(pth)
	if err != nil {
		return nil, err
	}

	return &JsonResponse{
		StatusCode: 200,
		Model:      types,
	}, nil
}

func (r *projectRouter) listVersionHandler(ctx HttpContext) (*JsonResponse, error) {
	project, ok := ctx.Args["project"]
	if !ok {
		return nil, errors.New("failed to get project from args")
	}

	idlType, ok := ctx.Args["type"]
	if !ok {
		return nil, errors.New("failed to get type from args")
	}

	pth := fmt.Sprintf("/projects/%s", project)

	ok = r.storage.Exists(pth)
	if !ok {
		return &JsonResponse{
			StatusCode: 404,
			Model:      fmt.Sprintf("project '%s' does not exist", project),
		}, nil
	}

	pth = fmt.Sprintf("/projects/%s/%s", project, idlType)

	ok = r.storage.Exists(pth)
	if !ok {
		return &JsonResponse{
			StatusCode: 404,
			Model:      fmt.Sprintf("project '%s' does not have type '%s'", project, idlType),
		}, nil
	}

	versions, err := r.storage.ListFolders(pth)
	if err != nil {
		return nil, err
	}

	return &JsonResponse{
		StatusCode: 200,
		Model:      versions,
	}, nil
}

func (r *projectRouter) submitVersion(ctx HttpContext) (*JsonResponse, error) {
	project, ok := ctx.Args["project"]
	if !ok {
		return nil, errors.New("failed to get project from args")
	}

	idlType, ok := ctx.Args["type"]
	if !ok {
		return nil, errors.New("failed to get type from args")
	}

	version, ok := ctx.Args["version"]
	if !ok {
		return nil, errors.New("failed to get version from args")
	}

	pth := fmt.Sprintf("/projects/%s/%s/%s", project, idlType, version)

	err := r.storage.MkDir(pth)
	if err != nil {
		return nil, err
	}

	pth = fmt.Sprintf("%s/data.tar.gz", pth)

	err = r.storage.CreateFile(pth, ctx.Body)
	if err != nil {
		return nil, err
	}
	return &JsonResponse{
		StatusCode: 201,
	}, nil
}

func (r *projectRouter) pullVersion(ctx HttpContext) (*DataResponse, error) {
	project, ok := ctx.Args["project"]
	if !ok {
		return nil, errors.New("failed to get project from args")
	}

	idlType, ok := ctx.Args["type"]
	if !ok {
		return nil, errors.New("failed to get type from args")
	}

	version, ok := ctx.Args["version"]
	if !ok {
		return nil, errors.New("failed to get version from args")
	}

	pth := fmt.Sprintf("/projects/%s", project)

	ok = r.storage.Exists(pth)
	if !ok {
		return &DataResponse{
			StatusCode: 404,
			Error:      fmt.Sprintf("project '%s' does not exist", project),
		}, nil
	}

	pth = fmt.Sprintf("/projects/%s/%s", project, idlType)

	ok = r.storage.Exists(pth)
	if !ok {
		return &DataResponse{
			StatusCode: 404,
			Error:      fmt.Sprintf("project '%s' does not have type '%s'", project, idlType),
		}, nil
	}

	pth = fmt.Sprintf("/projects/%s/%s/%s", project, idlType, version)

	ok = r.storage.Exists(pth)
	if !ok {
		return &DataResponse{
			StatusCode: 404,
			Error:      fmt.Sprintf("project '%s' with type '%s' does not have version '%s'", project, idlType, version),
		}, nil
	}

	f, err := r.storage.ReadFile(pth + "/data.tar.gz")
	if err != nil {
		return nil, err
	}

	return &DataResponse{
		StatusCode: 200,
		Data:       f,
	}, nil
}
