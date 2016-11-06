package gomock_sample

import (
	"testing"

	"github.com/golang/mock/gomock"

	mock "./mock_request"
)

func TestFake200Check(t *testing.T) {
	checker := NewStatusChecker(new(Fake200Request))
	t.Log(checker.DoCheck("http://localhost/"))
}

func TestMockLocalhost500(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRequest := mock.NewMockRequest(ctrl)
	mockRequest.EXPECT().DoRequest("http://localhost/").Return(500)
	checker := NewStatusChecker(mockRequest)
	t.Log(checker.DoCheck("http://localhost/"))
}
