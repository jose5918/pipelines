package api_server

import (
	"context"
	"fmt"
	"os"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	apiclient "github.com/googleprivate/ml/backend/api/go_http_client/pipeline_upload_client"
	params "github.com/googleprivate/ml/backend/api/go_http_client/pipeline_upload_client/pipeline_upload_service"
	model "github.com/googleprivate/ml/backend/api/go_http_client/pipeline_upload_model"
	"github.com/googleprivate/ml/backend/src/common/util"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	pipelineUploadFieldName      = "uploadfile"
	pipelineUploadPath           = "pipelines/upload"
	pipelineUploadServerBasePath = "/api/v1/namespaces/%s/services/ml-pipeline:8888/proxy/apis/v1beta1/%s"
	pipelineUploadContentTypeKey = "Content-Type"
)

type PipelineUploadInterface interface {
	UploadFile(filePath string, parameters *params.UploadPipelineParams) (*model.APIPipeline, error)
}

type PipelineUploadClient struct {
	apiClient *apiclient.PipelineUpload
}

func NewPipelineUploadClient(clientConfig clientcmd.ClientConfig, debug bool) (
	*PipelineUploadClient, error) {

	runtime, err := NewHTTPRuntime(clientConfig, debug)
	if err != nil {
		return nil, err
	}

	apiClient := apiclient.New(runtime, strfmt.Default)

	// Creating upload client
	return &PipelineUploadClient{
		apiClient: apiClient,
	}, nil
}

func (c *PipelineUploadClient) UploadFile(filePath string, parameters *params.UploadPipelineParams) (
	*model.APIPipeline, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, util.NewUserErrorWithSingleMessage(err,
			fmt.Sprintf("Failed to open file '%s'", filePath))
	}
	defer file.Close()

	parameters.Uploadfile = runtime.NamedReader(filePath, file)
	return c.Upload(parameters)
}

func (c *PipelineUploadClient) Upload(parameters *params.UploadPipelineParams) (*model.APIPipeline,
	error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), apiServerDefaultTimeout)
	defer cancel()

	// Make service all
	parameters.Context = ctx
	response, err := c.apiClient.PipelineUploadService.UploadPipeline(parameters, PassThroughAuth)

	if err != nil {
		if defaultError, ok := err.(*params.UploadPipelineDefault); ok {
			err = CreateErrorFromAPIStatus(defaultError.Payload.Error, defaultError.Payload.Code)
		} else {
			err = CreateErrorCouldNotRecoverAPIStatus(err)
		}

		return nil, util.NewUserError(err,
			fmt.Sprintf("Failed to upload pipeline. Params: '%v'", parameters),
			fmt.Sprintf("Failed to upload pipeline"))
	}

	return response.Payload, nil
}