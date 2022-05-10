package test

// For more about the AWS SDK middleware, see pages below:
// - https://daisuzu.hatenablog.com/entry/2021/10/31/225356
// - https://aws.github.io/aws-sdk-go-v2/docs/middleware/

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go/middleware"
)

// MiddlewareUnprocessed returns a middleware of dynamodb:BatchWriteItem to a single table
// that returns an error and items that isUnprocessed() returns true.
func MiddlewareUnprocessed(
	tableName string,
	collection []types.WriteRequest,
	isUnprocessed func(request types.WriteRequest) bool,
	err error,
) func(*middleware.Stack) error {
	resp := make([]types.WriteRequest, 0, len(collection))
	for _, item := range collection {
		if isUnprocessed(item) {
			resp = append(resp, item)
		}
	}

	return func(stack *middleware.Stack) error {
		return stack.Deserialize.Add(
			middleware.DeserializeMiddlewareFunc(
				"unprocessed",
				func(context.Context, middleware.DeserializeInput, middleware.DeserializeHandler) (middleware.DeserializeOutput, middleware.Metadata, error) {
					return middleware.DeserializeOutput{
						// TODO: research what are the appropriate type for the members of DeserializeOutput
						Result: dynamodb.BatchWriteItemOutput{},
						// RawResponse: ??
					}, middleware.Metadata{}, nil
				},
			),
			middleware.Before,
		)
	}
}
