package handlers

import (
	"net/http"

	"github.com/Yapo/goutils"

	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/domain"
	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/usecases"
)

// GetUserAdsHandler implements the handler interface and responds to /ads with
// related user ads
type GetUserAdsHandler struct {
	Interactor          usecases.GetUserAdsInteractor
	GetAdInteractor     usecases.GetAdInteractor
	Logger              GetUserAdsLogger
	UnitOfAccountSymbol string
	CurrencySymbol      string
}

// GetUserAdsLogger logger for GetUserAds Handler
type GetUserAdsLogger interface{}

// getUserAdsHandlerInput is the handler expected input
type getUserAdsHandlerInput struct {
	ListID string `path:"listID"`
}

// getUserRequestOutput is the handler output
type getUserRequestOutput struct {
	Ads []adsOutput `json:"ads"`
}

// adsOutput is the main output struct
type adsOutput struct {
	ID        string      `json:"id"`
	Category  string      `json:"category"`
	Title     string      `json:"title"`
	Price     float64     `json:"price"`
	Currency  string      `json:"currency"`
	Image     imageOutput `json:"images"`
	IsRelated bool        `json:"isRelated"`
	URL       string      `json:"url"`
}

// imageOutput is the output struct for images
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

// Execute get user ads for the current adview and returns related ads list
func (h *GetUserAdsHandler) Execute(ig InputGetter) *goutils.Response {
	input, response := ig()
	if response != nil {
		return response
	}
	in := input.(*getUserAdsHandlerInput)

	currentAd, err := h.GetAdInteractor.GetAd(in.ListID)
	if err != nil {
		return &goutils.Response{
			Code: http.StatusNoContent,
		}
	}

	resp, err := h.Interactor.GetUserAds(currentAd.UserID, currentAd.ID)
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

// fillResponse parses domain struct to expected handler output
func (h *GetUserAdsHandler) fillResponse(ads domain.Ads) []adsOutput {
	resp := []adsOutput{}
	for _, ad := range ads {
		adOutTemp := adsOutput{
			ID:       ad.ID,
			Title:    ad.Subject,
			Price:    ad.Price,
			Currency: ad.Currency,
			Category: ad.CategoryID,
			Image: imageOutput{
				Full:   ad.Image.Full,
				Medium: ad.Image.Medium,
				Small:  ad.Image.Small,
			},
			URL:       ad.URL,
			IsRelated: ad.IsRelated,
		}
		if ad.Currency == "uf" {
			adOutTemp.Currency = h.UnitOfAccountSymbol
			adOutTemp.Price = (adOutTemp.Price / 100)
		} else {
			adOutTemp.Currency = h.CurrencySymbol
		}
		resp = append(resp, adOutTemp)
	}
	return resp
}
