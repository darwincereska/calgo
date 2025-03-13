package calculator

import (
	"fmt"
)

var unitConversions = map[string]map[string]float64{
    "length": {
        "m":    1.0,
        "cm":   0.01,
        "mm":   0.001,
        "km":   1000.0,
        "inch": 0.0254,
        "ft":   0.3048,
        "yd":   0.9144,
        "mil":  0.0000254,
        "point": 0.000352778,
        "line": 0.002116667,
        "hand": 0.1016,
        "rod":  5.0292,
        "chain": 20.1168,
        "furlong": 201.168,
        "mile": 1609.344,
        "cable": 185.2,
        "nauticalmile": 1852.0,
        "league": 4828.032,
    },
    "area": {
        "m2": 1.0,
        "cm2": 0.0001,
        "mm2": 0.000001,
        "km2": 1000000.0,
        "hectare": 10000.0,
        "are": 100.0,
        "acre": 4046.8564224,
        "sqft": 0.092903,
        "sqyd": 0.836127,
        "sqin": 0.00064516,
    },
    "volume": {
        "m3": 1.0,
        "cm3": 0.000001,
        "mm3": 0.000000001,
        "liter": 0.001,
        "ml": 0.000001,
        "gallon": 0.003785412,
        "quart": 0.000946353,
        "pint": 0.000473176,
        "cup": 0.000236588,
        "floz": 0.0000295735,
        "tbsp": 0.0000147868,
        "tsp": 0.00000492892,
    },
    "weight": {
        "g": 1.0,
        "kg": 1000.0,
        "mg": 0.001,
        "tonne": 1000000.0,
        "lb": 453.59237,
        "oz": 28.349523125,
        "stone": 6350.29318,
        "carat": 0.2,
        "centner": 100000.0,
    },
    "time": {
        "second": 1.0,
        "minute": 60.0,
        "hour": 3600.0,
        "day": 86400.0,
        "week": 604800.0,
        "month": 2629746.0,
        "year": 31556952.0,
    },
    "data": {
        "bit": 1.0,
        "byte": 8.0,
        "kb": 8000.0,
        "mb": 8000000.0,
        "gb": 8000000000.0,
        "tb": 8000000000000.0,
        "kib": 8192.0,
        "mib": 8388608.0,
        "gib": 8589934592.0,
        "tib": 8796093022208.0,
    },
    "angle": {
        "radian": 1.0,
        "degree": 0.0174533,
    },
}

func convertUnit(value float64, from, to string) (float64, error) {
    // Find the category for the input unit
    var fromCategory, toCategory string
    for category, units := range unitConversions {
        if _, exists := units[from]; exists {
            fromCategory = category
        }
        if _, exists := units[to]; exists {
            toCategory = category
        }
    }

    // Check if units are from same category
    if fromCategory == "" || toCategory == "" || fromCategory != toCategory {
        return 0, fmt.Errorf("incompatible units: %s and %s", from, to)
    }

    // Convert to base unit then to target unit
    baseValue := value * unitConversions[fromCategory][from]
    result := baseValue / unitConversions[fromCategory][to]

    return result, nil
}