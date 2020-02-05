package handlers

import (
	"math"
	"net/http"

	"github.com/Yapo/goutils"

	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/domain"
	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/usecases"
)

// GetUserAds implements the handler interface and responds to /getUser with getUsered
// viewed ads
type GetUserAdsHandler struct {
	Interactor            usecases.GetUserAdsInteractor
	GetUserDataInteractor usecases.GetUserDataInteractor
	Logger                GetUserAdsLogger
	UnitOfAccountSymbol   string
	CurrencySymbol        string
}

// GetUserAdsLogger logger for TrackerHandler
type GetUserAdsLogger interface {
}

type getUserAdsHandlerInput struct {
	SHA1Email string `path:"token"`
}

type getUserRequestOutput struct {
	Ads []adsOutput `json:"ads"`
}

type adsOutput struct {
	ID       string      `json:"id"`
	Title    string      `json:"title"`
	Price    float64     `json:"price"`
	Currency string      `json:"currency"`
	URL      string      `json:"url"`
	Image    imageOutput `json:"images"`
}

type imageOutput struct {
	Full   string `json:"full,omitempty"`
	Medium string `json:"medium,omitempty"`
	Small  string `json:"small,omitempty"`
}

// Input returns a fresh, empty instance of getUserAdsHandlerInput
func (*GetUserAdsHandler) Input(ir InputRequest) HandlerInput {
	input := getUserAdsHandlerInput{}
	ir.Set(&input).FromPath()
	return &input
}

// Execute getUsers current adview and returns getUsered viewed ads list
func (h *GetUserAdsHandler) Execute(ig InputGetter) *goutils.Response {
	input, response := ig()
	if response != nil {
		return response
	}
	in := input.(*getUserAdsHandlerInput)

	userData, err := h.GetUserDataInteractor.GetUserData(in.SHA1Email)
	if err != nil {
		return &goutils.Response{
			Code: http.StatusNoContent,
		}
	}
	resp, err := h.Interactor.GetUserAds(userData.UserID)
	if err != nil {
		return &goutils.Response{
			Code: http.StatusNoContent,
		}
	}
	body := getUserRequestOutput{
		Ads: h.fillResponse(resp),
	}
	if len(body.Ads) == 0 {
		return &goutils.Response{
			Code: http.StatusNoContent,
		}
	}
	return &goutils.Response{
		Code: http.StatusOK,
		Body: body,
	}
}

func (h *GetUserAdsHandler) fillResponse(ads domain.Ads) []adsOutput {
	resp := []adsOutput{}
	for _, ad := range ads {
		adOutTemp := adsOutput{
			ID:       ad.ID,
			Title:    ad.Subject,
			Price:    ad.Price,
			Currency: ad.Currency,
			Image: imageOutput{
				Full:   ad.Image.Full,
				Medium: ad.Image.Medium,
				Small:  ad.Image.Small,
			},
		}
		if ad.UnitOfAccount > 0 {
			adOutTemp.Price = ad.UnitOfAccount
			adOutTemp.Currency = h.UnitOfAccountSymbol
		} else {
			adOutTemp.Currency = h.CurrencySymbol
		}
		// Round the output with 2 decimals
		adOutTemp.Price = math.Round(adOutTemp.Price*100) / 100
		adOutTemp.URL = ad.URL
		resp = append(resp, adOutTemp)
	}
	return resp
}
