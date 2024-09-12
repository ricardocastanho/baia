package contracts

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
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

func (r *RealEstate) Save(ctx context.Context, driver neo4j.DriverWithContext) error {
	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			CREATE (r:RealEstate {
				id: $id,
				code: $code,
				type: $type,
				name: $name,
				description: $description,
				url: $url,
				price: $price,
				bedrooms: $bedrooms,
				bathrooms: $bathrooms,
				area: $area,
				garageSpaces: $garageSpaces,
				location: $location,
				furnished: $furnished,
				yearBuilt: $yearBuilt,
				photos: $photos,
				tags: $tags,
				forSale: $forSale,
				forRent: $forRent
			})
			RETURN r
		`

		_, err := tx.Run(ctx, query, map[string]any{
			"id":           r.ID,
			"code":         r.Code,
			"type":         r.Type,
			"name":         r.Name,
			"description":  r.Description,
			"url":          r.Url,
			"price":        r.Price,
			"bedrooms":     r.Bedrooms,
			"bathrooms":    r.Bathrooms,
			"area":         r.Area,
			"garageSpaces": r.GarageSpaces,
			"location":     r.Location,
			"furnished":    r.Furnished,
			"yearBuilt":    r.YearBuilt,
			"photos":       r.Photos,
			"tags":         r.Tags,
			"forSale":      r.ForSale,
			"forRent":      r.ForRent,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to execute query: %w", err)
		}

		return nil, nil
	})

	return err
}
