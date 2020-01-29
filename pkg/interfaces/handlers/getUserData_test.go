package handlers

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/Yapo/goutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/usecases"
)

type mockGetUserDataHandlerPrometheusDefaultLogger struct {
	mock.Mock
}

func (m *mockGetUserDataHandlerPrometheusDefaultLogger) LogBadRequest(input interface{}) {
	m.Called(input)
}

func (m *mockGetUserDataHandlerPrometheusDefaultLogger) LogErrorGettingInternalData(err error) {
	m.Called(err)
}

type mockUserProfileInteractor struct {
	mock.Mock
}

func (m *mockUserProfileInteractor) GetUser(mail string) (usecases.UserBasicData, error) {
	args := m.Called(mail)
	return args.Get(0).(usecases.UserBasicData), args.Error(1)
}

func TestGetUserDataHandlerInput(t *testing.T) {
	m := mockUserProfileInteractor{}
	mMockInputRequest := MockInputRequest{}
	mMockTargetRequest := MockTargetRequest{}
	mMockInputRequest.On("Set", mock.AnythingOfType("*handlers.getUserDataRequestInput")).Return(&mMockTargetRequest)
	mMockTargetRequest.On("FromQuery").Return()

	h := GetUserDataHandler{
		Interactor: &m,
	}
	input := h.Input(&mMockInputRequest)

	var expected *getUserDataRequestInput
	assert.IsType(t, expected, input)
	m.AssertExpectations(t)
}

func TestGetUserDataHandlerDataRunOK(t *testing.T) {
	mInteractor := &mockUserProfileInteractor{}
	var userb usecases.UserBasicData
	emailV := regexp.MustCompile("@")
	mInteractor.On("GetUser", "hola@mail.com").Return(userb, nil)
	h := GetUserDataHandler{
		Interactor:    mInteractor,
		EmailValidate: emailV,
	}
	input := &getUserDataRequestInput{
		Mail: "hola@mail.com",
	}
	getter := MakeMockInputGetter(input, nil)
	r := h.Execute(getter)

	expected := &goutils.Response{
		Code: http.StatusOK,
		Body: getUserDataRequestOutput{},
	}
	assert.Equal(t, expected, r)
	mInteractor.AssertExpectations(t)
}

func TestInternalUserDataHandlerForInternalDataRunError(t *testing.T) {
	mInteractor := &mockUserProfileInteractor{}
	mLogger := &mockGetUserDataHandlerPrometheusDefaultLogger{}
	err := fmt.Errorf("err")
	var userb usecases.UserBasicData
	emailV := regexp.MustCompile("@")

	mInteractor.On("GetUser", "hola@mail.com").Return(userb, err)
	mLogger.On("LogErrorGettingInternalData", err).Once()

	h := GetUserDataHandler{
		Interactor:    mInteractor,
		EmailValidate: emailV,
		Logger:        mLogger,
	}
	input := &getUserDataRequestInput{
		Mail: "hola@mail.com",
	}
	getter := MakeMockInputGetter(input, nil)
	r := h.Execute(getter)

	expected := &goutils.Response{
		Code: http.StatusNoContent,
	}
	assert.Equal(t, expected, r)
	mLogger.AssertExpectations(t)
	mInteractor.AssertExpectations(t)
}

func TestInternalUserDataHandlerForInternalDataBadRequest(t *testing.T) {
	mInteractor := &mockUserProfileInteractor{}
	mLogger := &mockGetUserDataHandlerPrometheusDefaultLogger{}

	mLogger.On("LogBadRequest", mock.AnythingOfType("*goutils.Response")).Once()

	h := GetUserDataHandler{
		Interactor: mInteractor,
		Logger:     mLogger,
	}
	input := &getUserDataRequestInput{
		Mail: "mail@chao.cl",
	}
	getter := MakeMockInputGetter(input, &goutils.Response{
		Code: http.StatusBadRequest,
	})
	r := h.Execute(getter)

	expected := &goutils.Response{
		Code: http.StatusBadRequest,
	}
	assert.Equal(t, expected, r)
	mLogger.AssertExpectations(t)
	mInteractor.AssertExpectations(t)
}

func TestGetUserDataHandlerEmptyMail(t *testing.T) {
	mInteractor := &mockUserProfileInteractor{}
	mLogger := &mockGetUserDataHandlerPrometheusDefaultLogger{}

	emailV := regexp.MustCompile("@")

	mLogger.On("LogBadRequest", fmt.Errorf("Email is empty\n")).Once()
	h := GetUserDataHandler{
		Interactor:    mInteractor,
		EmailValidate: emailV,
		Logger:        mLogger,
	}
	input := &getUserDataRequestInput{
		Mail: "",
	}
	getter := MakeMockInputGetter(input, nil)
	r := h.Execute(getter)

	expected := &goutils.Response{
		Code: http.StatusBadRequest,
	}
	assert.Equal(t, expected, r)
	mLogger.AssertExpectations(t)
	mInteractor.AssertExpectations(t)
}

func TestGetUserDataHandlerBadMail(t *testing.T) {
	mInteractor := &mockUserProfileInteractor{}
	mLogger := &mockGetUserDataHandlerPrometheusDefaultLogger{}

	emailV := regexp.MustCompile("@")

	mLogger.On("LogBadRequest", fmt.Errorf("Email is invalid\n")).Once()
	h := GetUserDataHandler{
		Interactor:    mInteractor,
		EmailValidate: emailV,
		Logger:        mLogger,
	}
	input := &getUserDataRequestInput{
		Mail: "asdfg",
	}
	getter := MakeMockInputGetter(input, nil)
	r := h.Execute(getter)

	expected := &goutils.Response{
		Code: http.StatusBadRequest,
	}
	assert.Equal(t, expected, r)
	mLogger.AssertExpectations(t)
	mInteractor.AssertExpectations(t)
}
