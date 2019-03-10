package sam_router

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

type MethodType string

const (
	GET    MethodType = "GET"
	POST   MethodType = "POST"
	DELETE MethodType = "DELETE"
	PUT    MethodType = "PUT"
)

type APIGatewayRequest struct {
	events.APIGatewayProxyRequest
	Next Handler
}

type APIGatewayResponse struct {
	events.APIGatewayProxyResponse
}

func (r *APIGatewayRequest) SetNext(next Handler) {
	r.Next = next
}

func (r *APIGatewayRequest) GetNext() Handler {
	return r.Next
}

func (r *APIGatewayRequest) GetTag() Tag {
	return generateAPIGatewayTag(MethodType(r.RequestContext.HTTPMethod), r.RequestContext.ResourcePath)
}

type APIGatewayHandler func(request APIGatewayRequest) events.APIGatewayProxyResponse

type APIGatewayRouter struct {
	Router *Router
}

func (r *APIGatewayRouter) Get(path string, handlers ...Handler) {
	r.addPath(GET, path, handlers)
}

func (r *APIGatewayRouter) Post(path string, handlers ...Handler) {
	r.addPath(POST, path, handlers)
}

func (r *APIGatewayRouter) Delete(path string, handlers ...Handler) {
	r.addPath(DELETE, path, handlers)
}

func (r *APIGatewayRouter) Put(path string, handlers ...Handler) {
	r.addPath(PUT, path, handlers)
}

func emptyHandler(request Request) Response {
	return &APIGatewayResponse{
		APIGatewayProxyResponse: events.APIGatewayProxyResponse{
			StatusCode: 404,
		},
	}
}

func (r *APIGatewayRouter) addPath(method MethodType, path string, handlers []Handler) {
	if r.Router == nil {
		r.Router = NewRouter()
	}
	r.Router.AddRoute(generateAPIGatewayTag(method, path), handlers)
}

func (r *APIGatewayRouter) Dispatch(request events.APIGatewayProxyRequest, baseHandler APIGatewayHandler) events.APIGatewayProxyResponse {
	req := &APIGatewayRequest{APIGatewayProxyRequest: request}
	res := r.Router.Dispatch(req, emptyHandler)
	return res.(*APIGatewayResponse).APIGatewayProxyResponse
}

func generateAPIGatewayTag(method MethodType, path string) Tag {
	return Tag(fmt.Sprintf("%s::%s", method, path))
}
