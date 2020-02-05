package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Yapo/goutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/domain"
)

func TestTrackHandlernput(t *testing.T) {
	var h TrackHandler
	mMockInputRequest := MockInputRequest{}
	mMockInputRequest.On("Set",
		mock.AnythingOfType("*handlers.trackHandlerInput")).Return(&mMockInputRequest)
	mMockInputRequest.On("FromCookies").Return(&mMockInputRequest)
	mMockInputRequest.On("FromPath").Return(&mMockInputRequest)
	input := h.Input(&mMockInputRequest)
	var expected *trackHandlerInput
	assert.IsType(t, expected, input)
	mMockInputRequest.AssertExpectations(t)
}

type mockTrackerHandlerLogger struct {
	mock.Mock
}

func (m *mockTrackerHandlerLogger) LogFromUser(userID, listID string) {
	m.Called(userID, listID)
}

func (m *mockTrackerHandlerLogger) LogFromVisitor(visitorID, listID string) {
	m.Called(visitorID, listID)
}

type mockTrackHandlerInteractor struct {
	mock.Mock
}

func (m *mockTrackHandlerInteractor) TrackViewedAds(user domain.User, ad domain.Ad) (domain.Ads, error) {
	args := m.Called(user, ad)
	return args.Get(0).(domain.Ads), args.Error(1)
}

type mockSessionRepo struct {
	mock.Mock
}

func (m *mockSessionRepo) GetUserID(accSession string) (string, error) {
	args := m.Called(accSession)
	return args.String(0), args.Error(1)
}

func TestTrackHandlerOK(t *testing.T) {
	mInteractor := &mockTrackHandlerInteractor{}
	mSessionRepo := &mockSessionRepo{}
	mLogger := &mockTrackerHandlerLogger{}
	mLogger.On("LogFromUser", mock.AnythingOfType("string"),
		mock.AnythingOfType("string"))
	mInteractor.On("TrackViewedAds", mock.AnythingOfType("domain.User"),
		mock.AnythingOfType("domain.Ad")).Return(domain.Ads{{ID: "1"}}, nil)
	mSessionRepo.On("GetUserID", mock.AnythingOfType("string")).Return("", nil)
	h := TrackHandler{
		Interactor:  mInteractor,
		SessionRepo: mSessionRepo,
		Logger:      mLogger,
	}
	var input trackHandlerInput
	input.AccSession = "123"
	getter := MakeMockInputGetter(&input, nil)
	r := h.Execute(getter)

	expected := &goutils.Response{
		Code: http.StatusOK,
		Body: trackRequestOutput{
			Ads: []adsOutput{{ID: "1", URL: "1"}},
		},
	}
	assert.Equal(t, expected, r)
	mInteractor.AssertExpectations(t)
	mSessionRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestTrackHandlerFromVisitorOK(t *testing.T) {
	mInteractor := &mockTrackHandlerInteractor{}
	mSessionRepo := &mockSessionRepo{}
	mLogger := &mockTrackerHandlerLogger{}
	mLogger.On("LogFromVisitor", mock.AnythingOfType("string"),
		mock.AnythingOfType("string"))
	mInteractor.On("TrackViewedAds", mock.AnythingOfType("domain.User"),
		mock.AnythingOfType("domain.Ad")).Return(domain.Ads{{ID: "1"}}, nil)
	mSessionRepo.On("GetUserID",
		mock.AnythingOfType("string")).Return("", fmt.Errorf("e"))
	h := TrackHandler{
		Interactor:  mInteractor,
		SessionRepo: mSessionRepo,
		Logger:      mLogger,
	}
	var input trackHandlerInput
	input.AccSession = "123"
	getter := MakeMockInputGetter(&input, nil)
	r := h.Execute(getter)

	expected := &goutils.Response{
		Code: http.StatusOK,
		Body: trackRequestOutput{
			Ads: []adsOutput{{ID: "1", URL: "1"}},
		},
	}
	assert.Equal(t, expected, r)
	mInteractor.AssertExpectations(t)
	mSessionRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestTrackHandlerNoData(t *testing.T) {
	mInteractor := &mockTrackHandlerInteractor{}
	mSessionRepo := &mockSessionRepo{}
	mLogger := &mockTrackerHandlerLogger{}
	mLogger.On("LogFromUser", mock.AnythingOfType("string"),
		mock.AnythingOfType("string"))
	mInteractor.On("TrackViewedAds", mock.AnythingOfType("domain.User"),
		mock.AnythingOfType("domain.Ad")).Return(domain.Ads{}, nil)
	mSessionRepo.On("GetUserID", mock.AnythingOfType("string")).Return("", nil)
	h := TrackHandler{
		Interactor:  mInteractor,
		SessionRepo: mSessionRepo,
		Logger:      mLogger,
	}
	var input trackHandlerInput
	input.AccSession = "123"
	getter := MakeMockInputGetter(&input, nil)
	r := h.Execute(getter)

	expected := &goutils.Response{
		Code: http.StatusNoContent,
	}
	assert.Equal(t, expected, r)
	mInteractor.AssertExpectations(t)
	mSessionRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestTrackHandlerTrackViewedAdError(t *testing.T) {
	mInteractor := &mockTrackHandlerInteractor{}
	mSessionRepo := &mockSessionRepo{}
	mLogger := &mockTrackerHandlerLogger{}
	e := fmt.Errorf("error")
	mLogger.On("LogFromUser", mock.AnythingOfType("string"),
		mock.AnythingOfType("string"))
	mInteractor.On("TrackViewedAds", mock.AnythingOfType("domain.User"),
		mock.AnythingOfType("domain.Ad")).Return(domain.Ads{}, e)
	mSessionRepo.On("GetUserID", mock.AnythingOfType("string")).Return("", nil)
	h := TrackHandler{
		Interactor:  mInteractor,
		SessionRepo: mSessionRepo,
		Logger:      mLogger,
	}
	var input trackHandlerInput
	input.AccSession = "123"
	getter := MakeMockInputGetter(&input, nil)
	r := h.Execute(getter)

	expected := &goutils.Response{
		Code: http.StatusNoContent,
	}
	assert.Equal(t, expected, r)
	mInteractor.AssertExpectations(t)
	mSessionRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestTrackHandlerNotValidVisitorOrUser(t *testing.T) {
	mInteractor := &mockTrackHandlerInteractor{}
	mSessionRepo := &mockSessionRepo{}
	h := TrackHandler{
		Interactor:  mInteractor,
		SessionRepo: mSessionRepo,
	}
	var input trackHandlerInput
	getter := MakeMockInputGetter(&input, nil)
	r := h.Execute(getter)

	expected := &goutils.Response{
		Code: http.StatusNoContent,
	}
	assert.Equal(t, expected, r)
	mInteractor.AssertExpectations(t)
	mSessionRepo.AssertExpectations(t)
}

func TestTrackHandlerBadInput(t *testing.T) {
	mInteractor := &mockTrackHandlerInteractor{}
	mSessionRepo := &mockSessionRepo{}
	h := TrackHandler{
		Interactor:  mInteractor,
		SessionRepo: mSessionRepo,
	}
	var input trackHandlerInput
	getter := MakeMockInputGetter(&input,
		&goutils.Response{Code: http.StatusNoContent})
	r := h.Execute(getter)

	expected := &goutils.Response{
		Code: http.StatusNoContent,
	}
	assert.Equal(t, expected, r)
	mInteractor.AssertExpectations(t)
	mSessionRepo.AssertExpectations(t)
}

func TestFillResponse(t *testing.T) {
	mInteractor := &mockTrackHandlerInteractor{}
	mSessionRepo := &mockSessionRepo{}
	h := TrackHandler{
		Interactor:          mInteractor,
		SessionRepo:         mSessionRepo,
		UnitOfAccountSymbol: "U.F.",
		CurrencySymbol:      "$",
		AdViewLink:          "test.com/vi/",
	}

	r := h.fillResponse(
		domain.Ads{
			{ID: "1", UnitOfAccount: 123.0},
			{ID: "2", Price: 124.0},
		},
	)
	expected := []adsOutput{
		{ID: "1", Price: 123.0, Currency: h.UnitOfAccountSymbol, URL: h.AdViewLink + "1"},
		{ID: "2", Price: 124.0, Currency: h.CurrencySymbol, URL: h.AdViewLink + "2"},
	}
	assert.Equal(t, expected, r)
	mInteractor.AssertExpectations(t)
	mSessionRepo.AssertExpectations(t)
}
