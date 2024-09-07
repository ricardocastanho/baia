package contracts

import (
	"errors"
	"strconv"
	"strings"
)

type RealState struct {
	ID              string
	Code            string
	Name            string
	Description     string
	Url             string
	Price           int
	Bedrooms        int
	Bathrooms       int
	Area            int
	GarageSpaces    int
	Type            string
	Neighborhood    string
	Furnished       bool
	YearBuilt       int
	Characteristics []string
	ForSale         bool
	ForRent         bool
}

func (r *RealState) SetCode(text string) error {
	if !strings.Contains(text, "Cód. V") {
		return errors.New("invalid format: keywork 'Cód. V' is missing")
	}

	code := strings.TrimSpace(strings.Replace(text, "Cód. V", "", 1))

	r.Code = code

	return nil
}

func (r *RealState) SetName(name string) error {
	r.Name = strings.TrimSpace(name)
	return nil
}

func (r *RealState) SetDescription(description string) error {
	r.Description = strings.TrimSpace(description)
	return nil
}

func (r *RealState) SetPrice(text string) error {
	if !strings.Contains(text, "R$") {
		return errors.New("invalid format: keywork 'R$' is missing")
	}

	priceText := strings.TrimSpace(strings.Split(text, "R$")[1])
	priceText = strings.ReplaceAll(priceText, ".", "")
	priceText = strings.Replace(priceText, ",", ".", 1)

	priceFloat, err := strconv.ParseFloat(priceText, 64)
	if err != nil {
		return errors.New("error while converting the price: " + err.Error())
	}

	r.Price = int(priceFloat)

	return nil
}

func (r *RealState) SetBathrooms(text string) error {
	bedRoomsFloat, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return errors.New("error while converting the bathroom field: " + err.Error())
	}

	r.Bedrooms = int(bedRoomsFloat)

	return nil
}
