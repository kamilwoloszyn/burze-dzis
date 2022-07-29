package app

import (
	"context"
	"fmt"
	"reflect"

	"github.com/kamilwoloszyn/burze-dzis/domain"
	"github.com/kamilwoloszyn/burze-dzis/domain/vxml"
)

// Service contains all components necessary to run app properly.
type Service struct {
	burzeDzisClient BurzeDzis
}

// NewService creates a new service
func NewService(burzeDzisClient BurzeDzis) *Service {
	return &Service{
		burzeDzisClient: burzeDzisClient,
	}
}

// IsValidKey checks validity provided key.
// If an error will occur, then returning value will be set to false
// with error specified. If value is set false without any error, then
// a provided key is wrong.
func (s *Service) IsValidKey(ctx context.Context, keyReq vxml.APIKeyRequest) (bool, error) {
	return s.burzeDzisClient.IsValidKey(ctx, keyReq)
}

// CityLocation returns a coordinate of a city. If the city doesn't exist, then empty response without error will be returned.
func (s *Service) CityLocation(ctx context.Context, locationReq vxml.CityLocationRequest) (domain.CityLocation, error) {
	return s.burzeDzisClient.CityLocation(ctx, locationReq)
}

// Cities returns a list of suggestion of cities. If a provided keyword won't match any city, then empty list will be returned.
func (s *Service) Cities(ctx context.Context, citiesReq vxml.CitiesRequest) (domain.Cities, error) {
	return s.burzeDzisClient.Cities(ctx, citiesReq)
}

// StormSearch returns some data about thunderstorm in / arround the provided city.
// If the city doesn't exist, expect an error.
func (s *Service) StormSearch(ctx context.Context, stormReq vxml.StormSearchRequest) (domain.Storm, error) {
	if stormReq.Body.StormSearch.CityName != "" {
		cityLocation, err := s.CityLocation(
			ctx,
			vxml.NewCityLocationRequest(
				stormReq.Body.StormSearch.CityName,
				stormReq.Body.StormSearch.APIKey,
			),
		)
		if err != nil {
			return domain.Storm{}, fmt.Errorf("StormSearch: couldn't obtain a city coords: %v", err)
		}
		if equal := reflect.DeepEqual(cityLocation, domain.CityLocation{}); equal {
			return domain.Storm{}, fmt.Errorf("StormSearch: wrong coords received. Is a correct city provided ? ")
		}
		stormReq.Body.StormSearch.CoordY = cityLocation.CoordY
		stormReq.Body.StormSearch.CoordX = cityLocation.CoordX
	}
	return s.burzeDzisClient.StormSearch(ctx, stormReq)
}

// WeatherAlert returns weather alerts based on a provided city. If the city does't exist
// expect an error.
func (s *Service) WeatherAlert(ctx context.Context, alertReq vxml.WeatherAlertRequest) ([]domain.Alert, error) {
	if alertReq.Body.WeatherAlert.CityName != "" {
		cityLocation, err := s.CityLocation(
			ctx,
			vxml.NewCityLocationRequest(
				alertReq.Body.WeatherAlert.CityName,
				alertReq.Body.WeatherAlert.APIKey,
			),
		)
		if err != nil {
			return nil, fmt.Errorf("WeatherAlert: couldn't obtain a city coords: %v", err)
		}
		if equal := reflect.DeepEqual(cityLocation, domain.CityLocation{}); equal {
			return nil, fmt.Errorf("WeatherAlert: wrong coords received. Is a correct city provided ? ")
		}
		alertReq.Body.WeatherAlert.CoordY = cityLocation.CoordY
		alertReq.Body.WeatherAlert.CoordX = cityLocation.CoordX
	}
	return s.burzeDzisClient.WeatherAlert(ctx, alertReq)
}
