package gomock_sample

type Request interface {
	// only return status code
	DoRequest(url string) int
}
