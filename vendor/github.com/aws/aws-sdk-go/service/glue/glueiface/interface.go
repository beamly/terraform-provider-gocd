// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

// Package glueiface provides an interface to enable mocking the AWS Glue service client
// for testing your code.
//
// It is important to note that this interface will have breaking changes
// when the service model is updated and adds new API operations, paginators,
// and waiters.
package glueiface

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/glue"
)

// GlueAPI provides an interface to enable mocking the
// glue.Glue service client's API operation,
// paginators, and waiters. This make unit testing your code that calls out
// to the SDK's service client's calls easier.
//
// The best way to use this interface is so the SDK's service client's calls
// can be stubbed out for unit testing your code with the SDK without needing
// to inject custom request handlers into the SDK's request pipeline.
//
//    // myFunc uses an SDK service client to make a request to
//    // AWS Glue.
//    func myFunc(svc glueiface.GlueAPI) bool {
//        // Make svc.BatchCreatePartition request
//    }
//
//    func main() {
//        sess := session.New()
//        svc := glue.New(sess)
//
//        myFunc(svc)
//    }
//
// In your _test.go file:
//
//    // Define a mock struct to be used in your unit tests of myFunc.
//    type mockGlueClient struct {
//        glueiface.GlueAPI
//    }
//    func (m *mockGlueClient) BatchCreatePartition(input *glue.BatchCreatePartitionInput) (*glue.BatchCreatePartitionOutput, error) {
//        // mock response/functionality
//    }
//
//    func TestMyFunc(t *testing.T) {
//        // Setup Test
//        mockSvc := &mockGlueClient{}
//
//        myfunc(mockSvc)
//
//        // Verify myFunc's functionality
//    }
//
// It is important to note that this interface will have breaking changes
// when the service model is updated and adds new API operations, paginators,
// and waiters. Its suggested to use the pattern above for testing, or using
// tooling to generate mocks to satisfy the interfaces.
type GlueAPI interface {
	BatchCreatePartition(*glue.BatchCreatePartitionInput) (*glue.BatchCreatePartitionOutput, error)
	BatchCreatePartitionWithContext(aws.Context, *glue.BatchCreatePartitionInput, ...request.Option) (*glue.BatchCreatePartitionOutput, error)
	BatchCreatePartitionRequest(*glue.BatchCreatePartitionInput) (*request.Request, *glue.BatchCreatePartitionOutput)

	BatchDeleteConnection(*glue.BatchDeleteConnectionInput) (*glue.BatchDeleteConnectionOutput, error)
	BatchDeleteConnectionWithContext(aws.Context, *glue.BatchDeleteConnectionInput, ...request.Option) (*glue.BatchDeleteConnectionOutput, error)
	BatchDeleteConnectionRequest(*glue.BatchDeleteConnectionInput) (*request.Request, *glue.BatchDeleteConnectionOutput)

	BatchDeletePartition(*glue.BatchDeletePartitionInput) (*glue.BatchDeletePartitionOutput, error)
	BatchDeletePartitionWithContext(aws.Context, *glue.BatchDeletePartitionInput, ...request.Option) (*glue.BatchDeletePartitionOutput, error)
	BatchDeletePartitionRequest(*glue.BatchDeletePartitionInput) (*request.Request, *glue.BatchDeletePartitionOutput)

	BatchDeleteTable(*glue.BatchDeleteTableInput) (*glue.BatchDeleteTableOutput, error)
	BatchDeleteTableWithContext(aws.Context, *glue.BatchDeleteTableInput, ...request.Option) (*glue.BatchDeleteTableOutput, error)
	BatchDeleteTableRequest(*glue.BatchDeleteTableInput) (*request.Request, *glue.BatchDeleteTableOutput)

	BatchGetPartition(*glue.BatchGetPartitionInput) (*glue.BatchGetPartitionOutput, error)
	BatchGetPartitionWithContext(aws.Context, *glue.BatchGetPartitionInput, ...request.Option) (*glue.BatchGetPartitionOutput, error)
	BatchGetPartitionRequest(*glue.BatchGetPartitionInput) (*request.Request, *glue.BatchGetPartitionOutput)

	BatchStopJobRun(*glue.BatchStopJobRunInput) (*glue.BatchStopJobRunOutput, error)
	BatchStopJobRunWithContext(aws.Context, *glue.BatchStopJobRunInput, ...request.Option) (*glue.BatchStopJobRunOutput, error)
	BatchStopJobRunRequest(*glue.BatchStopJobRunInput) (*request.Request, *glue.BatchStopJobRunOutput)

	CreateClassifier(*glue.CreateClassifierInput) (*glue.CreateClassifierOutput, error)
	CreateClassifierWithContext(aws.Context, *glue.CreateClassifierInput, ...request.Option) (*glue.CreateClassifierOutput, error)
	CreateClassifierRequest(*glue.CreateClassifierInput) (*request.Request, *glue.CreateClassifierOutput)

	CreateConnection(*glue.CreateConnectionInput) (*glue.CreateConnectionOutput, error)
	CreateConnectionWithContext(aws.Context, *glue.CreateConnectionInput, ...request.Option) (*glue.CreateConnectionOutput, error)
	CreateConnectionRequest(*glue.CreateConnectionInput) (*request.Request, *glue.CreateConnectionOutput)

	CreateCrawler(*glue.CreateCrawlerInput) (*glue.CreateCrawlerOutput, error)
	CreateCrawlerWithContext(aws.Context, *glue.CreateCrawlerInput, ...request.Option) (*glue.CreateCrawlerOutput, error)
	CreateCrawlerRequest(*glue.CreateCrawlerInput) (*request.Request, *glue.CreateCrawlerOutput)

	CreateDatabase(*glue.CreateDatabaseInput) (*glue.CreateDatabaseOutput, error)
	CreateDatabaseWithContext(aws.Context, *glue.CreateDatabaseInput, ...request.Option) (*glue.CreateDatabaseOutput, error)
	CreateDatabaseRequest(*glue.CreateDatabaseInput) (*request.Request, *glue.CreateDatabaseOutput)

	CreateDevEndpoint(*glue.CreateDevEndpointInput) (*glue.CreateDevEndpointOutput, error)
	CreateDevEndpointWithContext(aws.Context, *glue.CreateDevEndpointInput, ...request.Option) (*glue.CreateDevEndpointOutput, error)
	CreateDevEndpointRequest(*glue.CreateDevEndpointInput) (*request.Request, *glue.CreateDevEndpointOutput)

	CreateJob(*glue.CreateJobInput) (*glue.CreateJobOutput, error)
	CreateJobWithContext(aws.Context, *glue.CreateJobInput, ...request.Option) (*glue.CreateJobOutput, error)
	CreateJobRequest(*glue.CreateJobInput) (*request.Request, *glue.CreateJobOutput)

	CreatePartition(*glue.CreatePartitionInput) (*glue.CreatePartitionOutput, error)
	CreatePartitionWithContext(aws.Context, *glue.CreatePartitionInput, ...request.Option) (*glue.CreatePartitionOutput, error)
	CreatePartitionRequest(*glue.CreatePartitionInput) (*request.Request, *glue.CreatePartitionOutput)

	CreateScript(*glue.CreateScriptInput) (*glue.CreateScriptOutput, error)
	CreateScriptWithContext(aws.Context, *glue.CreateScriptInput, ...request.Option) (*glue.CreateScriptOutput, error)
	CreateScriptRequest(*glue.CreateScriptInput) (*request.Request, *glue.CreateScriptOutput)

	CreateTable(*glue.CreateTableInput) (*glue.CreateTableOutput, error)
	CreateTableWithContext(aws.Context, *glue.CreateTableInput, ...request.Option) (*glue.CreateTableOutput, error)
	CreateTableRequest(*glue.CreateTableInput) (*request.Request, *glue.CreateTableOutput)

	CreateTrigger(*glue.CreateTriggerInput) (*glue.CreateTriggerOutput, error)
	CreateTriggerWithContext(aws.Context, *glue.CreateTriggerInput, ...request.Option) (*glue.CreateTriggerOutput, error)
	CreateTriggerRequest(*glue.CreateTriggerInput) (*request.Request, *glue.CreateTriggerOutput)

	CreateUserDefinedFunction(*glue.CreateUserDefinedFunctionInput) (*glue.CreateUserDefinedFunctionOutput, error)
	CreateUserDefinedFunctionWithContext(aws.Context, *glue.CreateUserDefinedFunctionInput, ...request.Option) (*glue.CreateUserDefinedFunctionOutput, error)
	CreateUserDefinedFunctionRequest(*glue.CreateUserDefinedFunctionInput) (*request.Request, *glue.CreateUserDefinedFunctionOutput)

	DeleteClassifier(*glue.DeleteClassifierInput) (*glue.DeleteClassifierOutput, error)
	DeleteClassifierWithContext(aws.Context, *glue.DeleteClassifierInput, ...request.Option) (*glue.DeleteClassifierOutput, error)
	DeleteClassifierRequest(*glue.DeleteClassifierInput) (*request.Request, *glue.DeleteClassifierOutput)

	DeleteConnection(*glue.DeleteConnectionInput) (*glue.DeleteConnectionOutput, error)
	DeleteConnectionWithContext(aws.Context, *glue.DeleteConnectionInput, ...request.Option) (*glue.DeleteConnectionOutput, error)
	DeleteConnectionRequest(*glue.DeleteConnectionInput) (*request.Request, *glue.DeleteConnectionOutput)

	DeleteCrawler(*glue.DeleteCrawlerInput) (*glue.DeleteCrawlerOutput, error)
	DeleteCrawlerWithContext(aws.Context, *glue.DeleteCrawlerInput, ...request.Option) (*glue.DeleteCrawlerOutput, error)
	DeleteCrawlerRequest(*glue.DeleteCrawlerInput) (*request.Request, *glue.DeleteCrawlerOutput)

	DeleteDatabase(*glue.DeleteDatabaseInput) (*glue.DeleteDatabaseOutput, error)
	DeleteDatabaseWithContext(aws.Context, *glue.DeleteDatabaseInput, ...request.Option) (*glue.DeleteDatabaseOutput, error)
	DeleteDatabaseRequest(*glue.DeleteDatabaseInput) (*request.Request, *glue.DeleteDatabaseOutput)

	DeleteDevEndpoint(*glue.DeleteDevEndpointInput) (*glue.DeleteDevEndpointOutput, error)
	DeleteDevEndpointWithContext(aws.Context, *glue.DeleteDevEndpointInput, ...request.Option) (*glue.DeleteDevEndpointOutput, error)
	DeleteDevEndpointRequest(*glue.DeleteDevEndpointInput) (*request.Request, *glue.DeleteDevEndpointOutput)

	DeleteJob(*glue.DeleteJobInput) (*glue.DeleteJobOutput, error)
	DeleteJobWithContext(aws.Context, *glue.DeleteJobInput, ...request.Option) (*glue.DeleteJobOutput, error)
	DeleteJobRequest(*glue.DeleteJobInput) (*request.Request, *glue.DeleteJobOutput)

	DeletePartition(*glue.DeletePartitionInput) (*glue.DeletePartitionOutput, error)
	DeletePartitionWithContext(aws.Context, *glue.DeletePartitionInput, ...request.Option) (*glue.DeletePartitionOutput, error)
	DeletePartitionRequest(*glue.DeletePartitionInput) (*request.Request, *glue.DeletePartitionOutput)

	DeleteTable(*glue.DeleteTableInput) (*glue.DeleteTableOutput, error)
	DeleteTableWithContext(aws.Context, *glue.DeleteTableInput, ...request.Option) (*glue.DeleteTableOutput, error)
	DeleteTableRequest(*glue.DeleteTableInput) (*request.Request, *glue.DeleteTableOutput)

	DeleteTrigger(*glue.DeleteTriggerInput) (*glue.DeleteTriggerOutput, error)
	DeleteTriggerWithContext(aws.Context, *glue.DeleteTriggerInput, ...request.Option) (*glue.DeleteTriggerOutput, error)
	DeleteTriggerRequest(*glue.DeleteTriggerInput) (*request.Request, *glue.DeleteTriggerOutput)

	DeleteUserDefinedFunction(*glue.DeleteUserDefinedFunctionInput) (*glue.DeleteUserDefinedFunctionOutput, error)
	DeleteUserDefinedFunctionWithContext(aws.Context, *glue.DeleteUserDefinedFunctionInput, ...request.Option) (*glue.DeleteUserDefinedFunctionOutput, error)
	DeleteUserDefinedFunctionRequest(*glue.DeleteUserDefinedFunctionInput) (*request.Request, *glue.DeleteUserDefinedFunctionOutput)

	GetCatalogImportStatus(*glue.GetCatalogImportStatusInput) (*glue.GetCatalogImportStatusOutput, error)
	GetCatalogImportStatusWithContext(aws.Context, *glue.GetCatalogImportStatusInput, ...request.Option) (*glue.GetCatalogImportStatusOutput, error)
	GetCatalogImportStatusRequest(*glue.GetCatalogImportStatusInput) (*request.Request, *glue.GetCatalogImportStatusOutput)

	GetClassifier(*glue.GetClassifierInput) (*glue.GetClassifierOutput, error)
	GetClassifierWithContext(aws.Context, *glue.GetClassifierInput, ...request.Option) (*glue.GetClassifierOutput, error)
	GetClassifierRequest(*glue.GetClassifierInput) (*request.Request, *glue.GetClassifierOutput)

	GetClassifiers(*glue.GetClassifiersInput) (*glue.GetClassifiersOutput, error)
	GetClassifiersWithContext(aws.Context, *glue.GetClassifiersInput, ...request.Option) (*glue.GetClassifiersOutput, error)
	GetClassifiersRequest(*glue.GetClassifiersInput) (*request.Request, *glue.GetClassifiersOutput)

	GetClassifiersPages(*glue.GetClassifiersInput, func(*glue.GetClassifiersOutput, bool) bool) error
	GetClassifiersPagesWithContext(aws.Context, *glue.GetClassifiersInput, func(*glue.GetClassifiersOutput, bool) bool, ...request.Option) error

	GetConnection(*glue.GetConnectionInput) (*glue.GetConnectionOutput, error)
	GetConnectionWithContext(aws.Context, *glue.GetConnectionInput, ...request.Option) (*glue.GetConnectionOutput, error)
	GetConnectionRequest(*glue.GetConnectionInput) (*request.Request, *glue.GetConnectionOutput)

	GetConnections(*glue.GetConnectionsInput) (*glue.GetConnectionsOutput, error)
	GetConnectionsWithContext(aws.Context, *glue.GetConnectionsInput, ...request.Option) (*glue.GetConnectionsOutput, error)
	GetConnectionsRequest(*glue.GetConnectionsInput) (*request.Request, *glue.GetConnectionsOutput)

	GetConnectionsPages(*glue.GetConnectionsInput, func(*glue.GetConnectionsOutput, bool) bool) error
	GetConnectionsPagesWithContext(aws.Context, *glue.GetConnectionsInput, func(*glue.GetConnectionsOutput, bool) bool, ...request.Option) error

	GetCrawler(*glue.GetCrawlerInput) (*glue.GetCrawlerOutput, error)
	GetCrawlerWithContext(aws.Context, *glue.GetCrawlerInput, ...request.Option) (*glue.GetCrawlerOutput, error)
	GetCrawlerRequest(*glue.GetCrawlerInput) (*request.Request, *glue.GetCrawlerOutput)

	GetCrawlerMetrics(*glue.GetCrawlerMetricsInput) (*glue.GetCrawlerMetricsOutput, error)
	GetCrawlerMetricsWithContext(aws.Context, *glue.GetCrawlerMetricsInput, ...request.Option) (*glue.GetCrawlerMetricsOutput, error)
	GetCrawlerMetricsRequest(*glue.GetCrawlerMetricsInput) (*request.Request, *glue.GetCrawlerMetricsOutput)

	GetCrawlerMetricsPages(*glue.GetCrawlerMetricsInput, func(*glue.GetCrawlerMetricsOutput, bool) bool) error
	GetCrawlerMetricsPagesWithContext(aws.Context, *glue.GetCrawlerMetricsInput, func(*glue.GetCrawlerMetricsOutput, bool) bool, ...request.Option) error

	GetCrawlers(*glue.GetCrawlersInput) (*glue.GetCrawlersOutput, error)
	GetCrawlersWithContext(aws.Context, *glue.GetCrawlersInput, ...request.Option) (*glue.GetCrawlersOutput, error)
	GetCrawlersRequest(*glue.GetCrawlersInput) (*request.Request, *glue.GetCrawlersOutput)

	GetCrawlersPages(*glue.GetCrawlersInput, func(*glue.GetCrawlersOutput, bool) bool) error
	GetCrawlersPagesWithContext(aws.Context, *glue.GetCrawlersInput, func(*glue.GetCrawlersOutput, bool) bool, ...request.Option) error

	GetDatabase(*glue.GetDatabaseInput) (*glue.GetDatabaseOutput, error)
	GetDatabaseWithContext(aws.Context, *glue.GetDatabaseInput, ...request.Option) (*glue.GetDatabaseOutput, error)
	GetDatabaseRequest(*glue.GetDatabaseInput) (*request.Request, *glue.GetDatabaseOutput)

	GetDatabases(*glue.GetDatabasesInput) (*glue.GetDatabasesOutput, error)
	GetDatabasesWithContext(aws.Context, *glue.GetDatabasesInput, ...request.Option) (*glue.GetDatabasesOutput, error)
	GetDatabasesRequest(*glue.GetDatabasesInput) (*request.Request, *glue.GetDatabasesOutput)

	GetDatabasesPages(*glue.GetDatabasesInput, func(*glue.GetDatabasesOutput, bool) bool) error
	GetDatabasesPagesWithContext(aws.Context, *glue.GetDatabasesInput, func(*glue.GetDatabasesOutput, bool) bool, ...request.Option) error

	GetDataflowGraph(*glue.GetDataflowGraphInput) (*glue.GetDataflowGraphOutput, error)
	GetDataflowGraphWithContext(aws.Context, *glue.GetDataflowGraphInput, ...request.Option) (*glue.GetDataflowGraphOutput, error)
	GetDataflowGraphRequest(*glue.GetDataflowGraphInput) (*request.Request, *glue.GetDataflowGraphOutput)

	GetDevEndpoint(*glue.GetDevEndpointInput) (*glue.GetDevEndpointOutput, error)
	GetDevEndpointWithContext(aws.Context, *glue.GetDevEndpointInput, ...request.Option) (*glue.GetDevEndpointOutput, error)
	GetDevEndpointRequest(*glue.GetDevEndpointInput) (*request.Request, *glue.GetDevEndpointOutput)

	GetDevEndpoints(*glue.GetDevEndpointsInput) (*glue.GetDevEndpointsOutput, error)
	GetDevEndpointsWithContext(aws.Context, *glue.GetDevEndpointsInput, ...request.Option) (*glue.GetDevEndpointsOutput, error)
	GetDevEndpointsRequest(*glue.GetDevEndpointsInput) (*request.Request, *glue.GetDevEndpointsOutput)

	GetDevEndpointsPages(*glue.GetDevEndpointsInput, func(*glue.GetDevEndpointsOutput, bool) bool) error
	GetDevEndpointsPagesWithContext(aws.Context, *glue.GetDevEndpointsInput, func(*glue.GetDevEndpointsOutput, bool) bool, ...request.Option) error

	GetJob(*glue.GetJobInput) (*glue.GetJobOutput, error)
	GetJobWithContext(aws.Context, *glue.GetJobInput, ...request.Option) (*glue.GetJobOutput, error)
	GetJobRequest(*glue.GetJobInput) (*request.Request, *glue.GetJobOutput)

	GetJobRun(*glue.GetJobRunInput) (*glue.GetJobRunOutput, error)
	GetJobRunWithContext(aws.Context, *glue.GetJobRunInput, ...request.Option) (*glue.GetJobRunOutput, error)
	GetJobRunRequest(*glue.GetJobRunInput) (*request.Request, *glue.GetJobRunOutput)

	GetJobRuns(*glue.GetJobRunsInput) (*glue.GetJobRunsOutput, error)
	GetJobRunsWithContext(aws.Context, *glue.GetJobRunsInput, ...request.Option) (*glue.GetJobRunsOutput, error)
	GetJobRunsRequest(*glue.GetJobRunsInput) (*request.Request, *glue.GetJobRunsOutput)

	GetJobRunsPages(*glue.GetJobRunsInput, func(*glue.GetJobRunsOutput, bool) bool) error
	GetJobRunsPagesWithContext(aws.Context, *glue.GetJobRunsInput, func(*glue.GetJobRunsOutput, bool) bool, ...request.Option) error

	GetJobs(*glue.GetJobsInput) (*glue.GetJobsOutput, error)
	GetJobsWithContext(aws.Context, *glue.GetJobsInput, ...request.Option) (*glue.GetJobsOutput, error)
	GetJobsRequest(*glue.GetJobsInput) (*request.Request, *glue.GetJobsOutput)

	GetJobsPages(*glue.GetJobsInput, func(*glue.GetJobsOutput, bool) bool) error
	GetJobsPagesWithContext(aws.Context, *glue.GetJobsInput, func(*glue.GetJobsOutput, bool) bool, ...request.Option) error

	GetMapping(*glue.GetMappingInput) (*glue.GetMappingOutput, error)
	GetMappingWithContext(aws.Context, *glue.GetMappingInput, ...request.Option) (*glue.GetMappingOutput, error)
	GetMappingRequest(*glue.GetMappingInput) (*request.Request, *glue.GetMappingOutput)

	GetPartition(*glue.GetPartitionInput) (*glue.GetPartitionOutput, error)
	GetPartitionWithContext(aws.Context, *glue.GetPartitionInput, ...request.Option) (*glue.GetPartitionOutput, error)
	GetPartitionRequest(*glue.GetPartitionInput) (*request.Request, *glue.GetPartitionOutput)

	GetPartitions(*glue.GetPartitionsInput) (*glue.GetPartitionsOutput, error)
	GetPartitionsWithContext(aws.Context, *glue.GetPartitionsInput, ...request.Option) (*glue.GetPartitionsOutput, error)
	GetPartitionsRequest(*glue.GetPartitionsInput) (*request.Request, *glue.GetPartitionsOutput)

	GetPartitionsPages(*glue.GetPartitionsInput, func(*glue.GetPartitionsOutput, bool) bool) error
	GetPartitionsPagesWithContext(aws.Context, *glue.GetPartitionsInput, func(*glue.GetPartitionsOutput, bool) bool, ...request.Option) error

	GetPlan(*glue.GetPlanInput) (*glue.GetPlanOutput, error)
	GetPlanWithContext(aws.Context, *glue.GetPlanInput, ...request.Option) (*glue.GetPlanOutput, error)
	GetPlanRequest(*glue.GetPlanInput) (*request.Request, *glue.GetPlanOutput)

	GetTable(*glue.GetTableInput) (*glue.GetTableOutput, error)
	GetTableWithContext(aws.Context, *glue.GetTableInput, ...request.Option) (*glue.GetTableOutput, error)
	GetTableRequest(*glue.GetTableInput) (*request.Request, *glue.GetTableOutput)

	GetTableVersions(*glue.GetTableVersionsInput) (*glue.GetTableVersionsOutput, error)
	GetTableVersionsWithContext(aws.Context, *glue.GetTableVersionsInput, ...request.Option) (*glue.GetTableVersionsOutput, error)
	GetTableVersionsRequest(*glue.GetTableVersionsInput) (*request.Request, *glue.GetTableVersionsOutput)

	GetTableVersionsPages(*glue.GetTableVersionsInput, func(*glue.GetTableVersionsOutput, bool) bool) error
	GetTableVersionsPagesWithContext(aws.Context, *glue.GetTableVersionsInput, func(*glue.GetTableVersionsOutput, bool) bool, ...request.Option) error

	GetTables(*glue.GetTablesInput) (*glue.GetTablesOutput, error)
	GetTablesWithContext(aws.Context, *glue.GetTablesInput, ...request.Option) (*glue.GetTablesOutput, error)
	GetTablesRequest(*glue.GetTablesInput) (*request.Request, *glue.GetTablesOutput)

	GetTablesPages(*glue.GetTablesInput, func(*glue.GetTablesOutput, bool) bool) error
	GetTablesPagesWithContext(aws.Context, *glue.GetTablesInput, func(*glue.GetTablesOutput, bool) bool, ...request.Option) error

	GetTrigger(*glue.GetTriggerInput) (*glue.GetTriggerOutput, error)
	GetTriggerWithContext(aws.Context, *glue.GetTriggerInput, ...request.Option) (*glue.GetTriggerOutput, error)
	GetTriggerRequest(*glue.GetTriggerInput) (*request.Request, *glue.GetTriggerOutput)

	GetTriggers(*glue.GetTriggersInput) (*glue.GetTriggersOutput, error)
	GetTriggersWithContext(aws.Context, *glue.GetTriggersInput, ...request.Option) (*glue.GetTriggersOutput, error)
	GetTriggersRequest(*glue.GetTriggersInput) (*request.Request, *glue.GetTriggersOutput)

	GetTriggersPages(*glue.GetTriggersInput, func(*glue.GetTriggersOutput, bool) bool) error
	GetTriggersPagesWithContext(aws.Context, *glue.GetTriggersInput, func(*glue.GetTriggersOutput, bool) bool, ...request.Option) error

	GetUserDefinedFunction(*glue.GetUserDefinedFunctionInput) (*glue.GetUserDefinedFunctionOutput, error)
	GetUserDefinedFunctionWithContext(aws.Context, *glue.GetUserDefinedFunctionInput, ...request.Option) (*glue.GetUserDefinedFunctionOutput, error)
	GetUserDefinedFunctionRequest(*glue.GetUserDefinedFunctionInput) (*request.Request, *glue.GetUserDefinedFunctionOutput)

	GetUserDefinedFunctions(*glue.GetUserDefinedFunctionsInput) (*glue.GetUserDefinedFunctionsOutput, error)
	GetUserDefinedFunctionsWithContext(aws.Context, *glue.GetUserDefinedFunctionsInput, ...request.Option) (*glue.GetUserDefinedFunctionsOutput, error)
	GetUserDefinedFunctionsRequest(*glue.GetUserDefinedFunctionsInput) (*request.Request, *glue.GetUserDefinedFunctionsOutput)

	GetUserDefinedFunctionsPages(*glue.GetUserDefinedFunctionsInput, func(*glue.GetUserDefinedFunctionsOutput, bool) bool) error
	GetUserDefinedFunctionsPagesWithContext(aws.Context, *glue.GetUserDefinedFunctionsInput, func(*glue.GetUserDefinedFunctionsOutput, bool) bool, ...request.Option) error

	ImportCatalogToGlue(*glue.ImportCatalogToGlueInput) (*glue.ImportCatalogToGlueOutput, error)
	ImportCatalogToGlueWithContext(aws.Context, *glue.ImportCatalogToGlueInput, ...request.Option) (*glue.ImportCatalogToGlueOutput, error)
	ImportCatalogToGlueRequest(*glue.ImportCatalogToGlueInput) (*request.Request, *glue.ImportCatalogToGlueOutput)

	ResetJobBookmark(*glue.ResetJobBookmarkInput) (*glue.ResetJobBookmarkOutput, error)
	ResetJobBookmarkWithContext(aws.Context, *glue.ResetJobBookmarkInput, ...request.Option) (*glue.ResetJobBookmarkOutput, error)
	ResetJobBookmarkRequest(*glue.ResetJobBookmarkInput) (*request.Request, *glue.ResetJobBookmarkOutput)

	StartCrawler(*glue.StartCrawlerInput) (*glue.StartCrawlerOutput, error)
	StartCrawlerWithContext(aws.Context, *glue.StartCrawlerInput, ...request.Option) (*glue.StartCrawlerOutput, error)
	StartCrawlerRequest(*glue.StartCrawlerInput) (*request.Request, *glue.StartCrawlerOutput)

	StartCrawlerSchedule(*glue.StartCrawlerScheduleInput) (*glue.StartCrawlerScheduleOutput, error)
	StartCrawlerScheduleWithContext(aws.Context, *glue.StartCrawlerScheduleInput, ...request.Option) (*glue.StartCrawlerScheduleOutput, error)
	StartCrawlerScheduleRequest(*glue.StartCrawlerScheduleInput) (*request.Request, *glue.StartCrawlerScheduleOutput)

	StartJobRun(*glue.StartJobRunInput) (*glue.StartJobRunOutput, error)
	StartJobRunWithContext(aws.Context, *glue.StartJobRunInput, ...request.Option) (*glue.StartJobRunOutput, error)
	StartJobRunRequest(*glue.StartJobRunInput) (*request.Request, *glue.StartJobRunOutput)

	StartTrigger(*glue.StartTriggerInput) (*glue.StartTriggerOutput, error)
	StartTriggerWithContext(aws.Context, *glue.StartTriggerInput, ...request.Option) (*glue.StartTriggerOutput, error)
	StartTriggerRequest(*glue.StartTriggerInput) (*request.Request, *glue.StartTriggerOutput)

	StopCrawler(*glue.StopCrawlerInput) (*glue.StopCrawlerOutput, error)
	StopCrawlerWithContext(aws.Context, *glue.StopCrawlerInput, ...request.Option) (*glue.StopCrawlerOutput, error)
	StopCrawlerRequest(*glue.StopCrawlerInput) (*request.Request, *glue.StopCrawlerOutput)

	StopCrawlerSchedule(*glue.StopCrawlerScheduleInput) (*glue.StopCrawlerScheduleOutput, error)
	StopCrawlerScheduleWithContext(aws.Context, *glue.StopCrawlerScheduleInput, ...request.Option) (*glue.StopCrawlerScheduleOutput, error)
	StopCrawlerScheduleRequest(*glue.StopCrawlerScheduleInput) (*request.Request, *glue.StopCrawlerScheduleOutput)

	StopTrigger(*glue.StopTriggerInput) (*glue.StopTriggerOutput, error)
	StopTriggerWithContext(aws.Context, *glue.StopTriggerInput, ...request.Option) (*glue.StopTriggerOutput, error)
	StopTriggerRequest(*glue.StopTriggerInput) (*request.Request, *glue.StopTriggerOutput)

	UpdateClassifier(*glue.UpdateClassifierInput) (*glue.UpdateClassifierOutput, error)
	UpdateClassifierWithContext(aws.Context, *glue.UpdateClassifierInput, ...request.Option) (*glue.UpdateClassifierOutput, error)
	UpdateClassifierRequest(*glue.UpdateClassifierInput) (*request.Request, *glue.UpdateClassifierOutput)

	UpdateConnection(*glue.UpdateConnectionInput) (*glue.UpdateConnectionOutput, error)
	UpdateConnectionWithContext(aws.Context, *glue.UpdateConnectionInput, ...request.Option) (*glue.UpdateConnectionOutput, error)
	UpdateConnectionRequest(*glue.UpdateConnectionInput) (*request.Request, *glue.UpdateConnectionOutput)

	UpdateCrawler(*glue.UpdateCrawlerInput) (*glue.UpdateCrawlerOutput, error)
	UpdateCrawlerWithContext(aws.Context, *glue.UpdateCrawlerInput, ...request.Option) (*glue.UpdateCrawlerOutput, error)
	UpdateCrawlerRequest(*glue.UpdateCrawlerInput) (*request.Request, *glue.UpdateCrawlerOutput)

	UpdateCrawlerSchedule(*glue.UpdateCrawlerScheduleInput) (*glue.UpdateCrawlerScheduleOutput, error)
	UpdateCrawlerScheduleWithContext(aws.Context, *glue.UpdateCrawlerScheduleInput, ...request.Option) (*glue.UpdateCrawlerScheduleOutput, error)
	UpdateCrawlerScheduleRequest(*glue.UpdateCrawlerScheduleInput) (*request.Request, *glue.UpdateCrawlerScheduleOutput)

	UpdateDatabase(*glue.UpdateDatabaseInput) (*glue.UpdateDatabaseOutput, error)
	UpdateDatabaseWithContext(aws.Context, *glue.UpdateDatabaseInput, ...request.Option) (*glue.UpdateDatabaseOutput, error)
	UpdateDatabaseRequest(*glue.UpdateDatabaseInput) (*request.Request, *glue.UpdateDatabaseOutput)

	UpdateDevEndpoint(*glue.UpdateDevEndpointInput) (*glue.UpdateDevEndpointOutput, error)
	UpdateDevEndpointWithContext(aws.Context, *glue.UpdateDevEndpointInput, ...request.Option) (*glue.UpdateDevEndpointOutput, error)
	UpdateDevEndpointRequest(*glue.UpdateDevEndpointInput) (*request.Request, *glue.UpdateDevEndpointOutput)

	UpdateJob(*glue.UpdateJobInput) (*glue.UpdateJobOutput, error)
	UpdateJobWithContext(aws.Context, *glue.UpdateJobInput, ...request.Option) (*glue.UpdateJobOutput, error)
	UpdateJobRequest(*glue.UpdateJobInput) (*request.Request, *glue.UpdateJobOutput)

	UpdatePartition(*glue.UpdatePartitionInput) (*glue.UpdatePartitionOutput, error)
	UpdatePartitionWithContext(aws.Context, *glue.UpdatePartitionInput, ...request.Option) (*glue.UpdatePartitionOutput, error)
	UpdatePartitionRequest(*glue.UpdatePartitionInput) (*request.Request, *glue.UpdatePartitionOutput)

	UpdateTable(*glue.UpdateTableInput) (*glue.UpdateTableOutput, error)
	UpdateTableWithContext(aws.Context, *glue.UpdateTableInput, ...request.Option) (*glue.UpdateTableOutput, error)
	UpdateTableRequest(*glue.UpdateTableInput) (*request.Request, *glue.UpdateTableOutput)

	UpdateTrigger(*glue.UpdateTriggerInput) (*glue.UpdateTriggerOutput, error)
	UpdateTriggerWithContext(aws.Context, *glue.UpdateTriggerInput, ...request.Option) (*glue.UpdateTriggerOutput, error)
	UpdateTriggerRequest(*glue.UpdateTriggerInput) (*request.Request, *glue.UpdateTriggerOutput)

	UpdateUserDefinedFunction(*glue.UpdateUserDefinedFunctionInput) (*glue.UpdateUserDefinedFunctionOutput, error)
	UpdateUserDefinedFunctionWithContext(aws.Context, *glue.UpdateUserDefinedFunctionInput, ...request.Option) (*glue.UpdateUserDefinedFunctionOutput, error)
	UpdateUserDefinedFunctionRequest(*glue.UpdateUserDefinedFunctionInput) (*request.Request, *glue.UpdateUserDefinedFunctionOutput)
}

var _ GlueAPI = (*glue.Glue)(nil)
