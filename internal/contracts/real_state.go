package contracts

import (
	"baia/internal/utils"
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
	ID             string
	Code           string
	Type           string
	Name           string
	NormalizedName string
	Description    string
	Url            string
	Price          int
	Bedrooms       int
	Bathrooms      int
	Area           int
	GarageSpaces   int
	City           string
	District       string
	Furnished      bool
	YearBuilt      int
	Photos         []string
	Tags           []string
	Agency         string
	ForSale        bool
	ForRent        bool
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

func (r *RealEstate) SetDistrict(text string) error {
	r.District = strings.TrimSpace(strings.ReplaceAll(text, "/\t", ""))
	return nil
}

func (r *RealEstate) SetCity(text string) error {
	r.City = strings.TrimSpace(strings.ReplaceAll(text, "/\t", ""))
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
		realEstateLabels := []string{"RealEstate"}
		priceLabels := []string{"Price"}

		if r.Type != "" {
			realEstateLabels = append(realEstateLabels, r.Type)
		}
		if r.ForSale {
			realEstateLabels = append(realEstateLabels, "ForSale")
			priceLabels = append(priceLabels, "SalePrice")
		}
		if r.ForRent {
			realEstateLabels = append(realEstateLabels, "ForRent")
			priceLabels = append(priceLabels, "RentalPrice")
		}

		realEstateLabelString := fmt.Sprintf("SET r:%s", strings.Join(realEstateLabels, ":"))
		priceLabelString := strings.Join(priceLabels, ":")

		query := fmt.Sprintf(`
			MERGE (r:RealEstate {code: $code})	
			ON CREATE SET
					r.id = randomUUID(),
					r.type = $type,
					r.name = $name,
					r.description = $description,
					r.url = $url,
					r.bedrooms = $bedrooms,
					r.bathrooms = $bathrooms,
					r.area = $area,
					r.garageSpaces = $garageSpaces,
					r.furnished = $furnished,
					r.yearBuilt = $yearBuilt,
					r.photos = $photos,
					r.tags = $tags,
					r.forSale = $forSale,
					r.forRent = $forRent,
					r.createdAt = datetime(),
					r.updatedAt = datetime()
			ON MATCH SET
					r.type = $type,
					r.name = $name,
					r.description = $description,
					r.url = $url,
					r.bedrooms = $bedrooms,
					r.bathrooms = $bathrooms,
					r.area = $area,
					r.garageSpaces = $garageSpaces,
					r.furnished = $furnished,
					r.yearBuilt = $yearBuilt,
					r.photos = $photos,
					r.tags = $tags,
					r.forSale = $forSale,
					r.forRent = $forRent,
					r.updatedAt = datetime()
			%s
			WITH r
			CALL {
				WITH r
				WITH r AS r2
				OPTIONAL MATCH (r2)-[old:LATEST_PRICE]->(oldPrice:%s)
				WHERE oldPrice IS NULL OR oldPrice.value <> $price
				DELETE old
				WITH r2, oldPrice
				CREATE (newPrice:%s {
					id: randomUUID(),
					value: $price,
					createdAt: datetime()
				})
				CREATE (r2)-[:LATEST_PRICE]->(newPrice)
				FOREACH (_ IN CASE WHEN oldPrice IS NOT NULL THEN [1] ELSE [] END |
					CREATE (newPrice)<-[:NEXT]-(oldPrice)
				)
				WITH r2, newPrice
				OPTIONAL MATCH (r2)-[:FIRST_PRICE]->(p:%s)
				WITH r2, newPrice, COUNT(p) AS existingFirst
				WHERE existingFirst = 0
				CREATE (r2)-[:FIRST_PRICE]->(newPrice)
			}
			MERGE (a:Agency {normalizedName: $normalizedAgencyName})
			ON CREATE SET
					a.id = randomUUID(),
					a.name = $agency,
					a.normalizedName = $normalizedAgencyName
			MERGE (r)-[:SELLED_BY]->(a)
			WITH r
			MERGE (e:Estate {name: "Rio Grande do Sul"})
			WITH r, e
			MERGE (c:City {normalizedName: $normalizedCityName})
			ON CREATE SET
					c.id = randomUUID(),
					c.name = $city,
					c.normalizedName = $normalizedCityName
			WITH r, e, c
			MERGE (c)-[:IN]->(e)
			MERGE (r)-[:IN]->(c)
			WITH r, c
			WHERE NOT $district = "" AND $district IS NOT NULL 
			MERGE (d:District {name: $district})
			ON CREATE SET
					d.id = randomUUID(),
					d.name = $district
			WITH r, d, c
			MERGE (d)-[:IN]->(c)
			WITH r, d
			CREATE (r)-[:IN]->(d)	
			RETURN r
		`, realEstateLabelString, priceLabelString, priceLabelString, priceLabelString)

		_, err := tx.Run(ctx, query, map[string]any{
			"code":                 r.Code,
			"type":                 r.Type,
			"name":                 r.Name,
			"description":          r.Description,
			"url":                  r.Url,
			"price":                r.Price,
			"bedrooms":             r.Bedrooms,
			"bathrooms":            r.Bathrooms,
			"area":                 r.Area,
			"city":                 r.City,
			"normalizedCityName":   utils.NormalizeCityName(r.City),
			"agency":               r.Agency,
			"normalizedAgencyName": utils.NormalizeCityName(r.Agency),
			"district":             r.District,
			"garageSpaces":         r.GarageSpaces,
			"furnished":            r.Furnished,
			"yearBuilt":            r.YearBuilt,
			"photos":               r.Photos,
			"tags":                 r.Tags,
			"forSale":              r.ForSale,
			"forRent":              r.ForRent,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to execute query: %w", err)
		}

		return nil, nil
	})

	return err
}
