package daemon

import (
	"net/http"
	"encoding/json"
	"calgo/backend/calculator"
)

type CalculationRequest struct {
    Expression string            `json:"expression"`
    Context    ContextRequest    `json:"context"`
}

type ContextRequest struct {
    Variables map[string]string  `json:"variables"`
    TimeZone  string            `json:"timezone"`
    PPI       float64           `json:"ppi"`
    EmSize    float64           `json:"emSize"`
}

type CalculationResponse struct {
    Result struct {
        Number   float64 `json:"Number"`   // Capital letters to match Go struct
        Unit     string  `json:"Unit,omitempty"`
        Type     int     `json:"Type"`
        IsError  bool    `json:"IsError"`
        ErrorMsg string  `json:"ErrorMsg,omitempty"`
    } `json:"result"`
    Variables map[string]calculator.Value `json:"variables"`
    Error    string                      `json:"error,omitempty"`
}

func parseRoutes() {
	// Enable CORS
	corsMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			
			next(w, r)
		}
	}

	// API endpoints
	http.HandleFunc("/api/calculate", corsMiddleware(handleCalculate))
}

func handleCalculate(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req CalculationRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest) 
        return
    }

    ctx := calculator.Context{
        Variables: make(map[string]calculator.Value),
        TimeZone: req.Context.TimeZone,
        PPI: req.Context.PPI,
        EmSize: req.Context.EmSize,
    }

    // Convert context variables
    for k, v := range req.Context.Variables {
        result, err := calculator.Calculate(v)
        if err != nil {
            ctx.Variables[k] = calculator.Value{
                IsError: true,
                ErrorMsg: err.Error(),
            }
        } else {
            ctx.Variables[k] = result
        }
    }

    // Calculate expression
    result, err := calculator.Calculate(req.Expression)
    response := CalculationResponse{}
    response.Result.Number = result.Number
    response.Result.Unit = result.Unit
    response.Result.Type = int(result.Type)
    response.Result.IsError = result.IsError
    response.Result.ErrorMsg = result.ErrorMsg
    response.Variables = ctx.Variables

    if err != nil {
        response.Error = err.Error()
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func Serve() {
	parseRoutes()
	http.ListenAndServe(":8888", nil)
}