package utils

import (
	"errors"
	"strconv"
	"strings"
)

func ParsePriceToInt(text string) (int, error) {
	if !strings.Contains(text, "R$") {
		return 0, errors.New("formato inválido: falta o símbolo de real 'R$'")
	}

	priceText := strings.TrimSpace(strings.Replace(text, "R$", "", 1))
	priceText = strings.ReplaceAll(priceText, ".", "")
	priceText = strings.Replace(priceText, ",", ".", 1)

	priceFloat, err := strconv.ParseFloat(priceText, 64)
	if err != nil {
		return 0, errors.New("erro ao converter o preço: " + err.Error())
	}

	return int(priceFloat), nil
}
