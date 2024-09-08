package contracts

import (
	"errors"
	"strconv"
	"strings"
)

const (
	House      string = "House"
	Apartment  string = "Apartment"
	Land       string = "Land"
	Commercial string = "Commercial"
	Industrial string = "Industrial"
)

type RealEstate struct {
	ID           string
	Code         string
	Type         string
	Name         string
	Description  string
	Url          string
	Price        int
	Bedrooms     int
	Bathrooms    int
	Area         int
	GarageSpaces int
	Location     string
	Furnished    bool
	YearBuilt    int
	Photos       []string
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

func (r *RealEstate) SetArea(text string) error {
	areaText := strings.Split(text, " m²")[0]
	areaText = strings.ReplaceAll(areaText, ".", "")
	areaText = strings.Split(areaText, ",")[0]
	number, err := strconv.Atoi(areaText)
	if err != nil {
		return errors.New("error while converting the area field: " + err.Error())
	}

	r.Area = number

	return nil
}

func (r *RealEstate) SetGarageSpaces(text string) error {
	number, err := strconv.Atoi(text)
	if err != nil {
		return errors.New("error while converting the garage spaces field: " + err.Error())
	}

	r.GarageSpaces = number

	return nil
}

func (r *RealEstate) SetLocation(text string) error {
	r.Location = strings.TrimSpace(strings.ReplaceAll(text, "/\t", ""))
	return nil
}

func (r *RealEstate) SetFurnished(is bool) error {
	r.Furnished = is
	return nil
}

func (r *RealEstate) SetYearBuilt(text string) error {
	number, err := strconv.Atoi(text)
	if err != nil {
		return errors.New("error while converting the year built field: " + err.Error())
	}

	r.YearBuilt = number

	return nil
}

func (r *RealEstate) SetPhoto(url string) error {
	r.Photos = append(r.Photos, url)
	return nil
}

func (r *RealEstate) SetTag(tag string) error {
	r.Tags = append(r.Tags, tag)
	return nil
}
