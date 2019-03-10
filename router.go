package sam_router

type Tag string

type Request interface {
	SetNext(next Handler)
	GetNext() Handler
	GetTag() Tag
}

type Response interface {
}

type Handler func(request Request) Response

type Router struct {
	Routes map[Tag]Handler
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) AddRoute(tag Tag, handlers []Handler) {
	n := len(handlers)
	var f Handler
	for i := n - 1; i >= 0; i-- {
		oldF := f
		c := handlers[i]
		f = func(request Request) Response {
			request.SetNext(oldF)
			return c(request)
		}
	}
	if r.Routes == nil {
		r.Routes = map[Tag]Handler{}
	}
	r.Routes[tag] = f
}

func (r *Router) Dispatch(request Request, baseHandler Handler) Response {
	tag := request.GetTag()
	f := r.Routes[tag]
	if f != nil {
		return f(request)
	}
	return baseHandler(request)
}
