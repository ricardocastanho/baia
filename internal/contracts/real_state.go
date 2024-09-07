package contracts

import (
	"errors"
	"strconv"
	"strings"
)

type RealEstate struct {
	ID           string
	Code         string
	Name         string
	Description  string
	Url          string
	Price        int
	Bedrooms     int
	Bathrooms    int
	Area         int
	GarageSpaces int
	Type         string
	Neighborhood string
	Furnished    bool
	YearBuilt    int
	Tags         []string
	ForSale      bool
	ForRent      bool
}

func (r *RealEstate) SetCode(text string) error {
	if !strings.Contains(text, "Cód. V") {
		return errors.New("invalid format: keywork 'Cód. V' is missing")
	}

	code := strings.TrimSpace(strings.Replace(text, "Cód. V", "", 1))

	r.Code = code

	return nil
}

func (r *RealEstate) SetName(name string) error {
	r.Name = strings.TrimSpace(name)
	return nil
}

func (r *RealEstate) SetDescription(description string) error {
	r.Description = strings.TrimSpace(description)
	return nil
}

func (r *RealEstate) SetPrice(text string) error {
	if !strings.Contains(text, "R$") {
		return errors.New("invalid format: keywork 'R$' is missing")
	}

	priceText := strings.TrimSpace(strings.Split(text, "R$")[1])
	priceText = strings.ReplaceAll(priceText, ".", "")
	priceText = strings.Split(priceText, ",")[0]

	number, err := strconv.Atoi(priceText)
	if err != nil {
		return errors.New("error while converting the price: " + err.Error())
	}

	r.Price = number

	return nil
}

func (r *RealEstate) SetBedrooms(text string) error {
	number, err := strconv.Atoi(text)
	if err != nil {
		return errors.New("error while converting the bedroom field: " + err.Error())
	}

	r.Bedrooms = number

	return nil
}

func (r *RealEstate) SetBathrooms(text string) error {
	number, err := strconv.Atoi(text)
	if err != nil {
		return errors.New("error while converting the bathroom field: " + err.Error())
	}

	r.Bathrooms = number

	return nil
}
