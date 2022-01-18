package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	actionstypes "github.com/forbole/bdjuno/v2/cmd/actions/types"
)


func ValidatorCommission(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Content-Type", "application/json")

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	var actionPayload actionstypes.AddressPayload
	err = json.Unmarshal(reqBody, &actionPayload)
	if err != nil {
		http.Error(w, "invalid payload: failed to unmarshal json", http.StatusInternalServerError)
		return
	}

	result, err := getValidatorCommission(actionPayload.Input.Address)
	if err != nil {
		errorHandler(w, err)
		return
	}

	data, _ := json.Marshal(result)
	w.Write(data)
}


func getValidatorCommission(address string) (response []actionstypes.ValidatorCommission, err error) {
	parseCtx, sources, err := getCtxAndSources()
	if err != nil {
		return response, err
	}

	// Get latest node height
	height, err := parseCtx.Node.LatestHeight()
	if err != nil {
		return response, fmt.Errorf("error while getting chain latest block height: %s", err)
	}

	// Get validator total commission value 
	commission, err := sources.DistrSource.ValidatorCommission(address, height)
		if err != nil {
		return response, err
	}

	response = make([]actionstypes.ValidatorCommission, len(commission))
	for index, commission := range commission {
		response[index] = actionstypes.ValidatorCommission{
			DecCoin:   commission,
			ValAddress: address,
		}
	}
	
	return response, nil
}