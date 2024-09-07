package contracts

import (
	"errors"
	"strconv"
	"strings"
)

type RealState struct {
	ID              string
	Cod             string
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

func (r *RealState) SetPrice(text string) error {
	if !strings.Contains(text, "R$") {
		return errors.New("invalid format: keywork 'R$' is missing")
	}

	priceText := strings.TrimSpace(strings.Replace(text, "R$", "", 1))
	priceText = strings.ReplaceAll(priceText, ".", "")
	priceText = strings.Replace(priceText, ",", ".", 1)

	priceFloat, err := strconv.ParseFloat(priceText, 64)
	if err != nil {
		return errors.New("error while converting the price: " + err.Error())
	}

	r.Price = int(priceFloat)

	return nil
}
