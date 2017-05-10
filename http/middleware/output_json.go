// Copyright 2017 Corey Scott http://www.sage42.org/
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
)

// OutputJSON will attempt to encode the supplied DTO into JSON bytes and add to the response.
// On success it will add status header of HTTP 200 (OK)
func OutputJSON(resp http.ResponseWriter, dto interface{}) error {
	if dto == nil {
		return errors.New("supplied DTO was empty")
	}

	encoder := json.NewEncoder(resp)
	err := encoder.Encode(dto)
	if err != nil {
		return err
	}

	resp.WriteHeader(http.StatusOK)
	return nil
}
