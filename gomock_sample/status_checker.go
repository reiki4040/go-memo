package gomock_sample

type Fake200Request struct{}

func (r *Fake200Request) DoRequest(url string) int {
	// dummy
	return 200
}

func NewStatusChecker(r Request) *StatusChecker {
	return &StatusChecker{
		request: r,
	}
}

type StatusChecker struct {
	request Request
}

func (c *StatusChecker) DoCheck(url string) bool {
	return (c.request.DoRequest(url) == 200)
}
