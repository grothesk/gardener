// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-License-Identifier: Apache-2.0

package secrets

import (
	"encoding/json"
	"fmt"

	"github.com/gardener/gardener/pkg/utils/infodata"
)

// StaticTokenDataType is the type used to denote an StaticTokenJSONData structure in the ShootState
const StaticTokenDataType = infodata.TypeVersion("staticToken")

func init() {
	infodata.Register(StaticTokenDataType, UnmarshalStaticToken)
}

// StaticTokenJSONData is the json representation of a StaticTokenInfoData
type StaticTokenJSONData struct {
	Tokens map[string]string `json:"tokens"`
}

// UnmarshalStaticToken unmarshals an StaticTokenJSONData into a StaticTokenInfoData.
func UnmarshalStaticToken(bytes []byte) (infodata.InfoData, error) {
	if bytes == nil {
		return nil, fmt.Errorf("no data given")
	}

	data := &StaticTokenJSONData{}
	if err := json.Unmarshal(bytes, data); err != nil {
		return nil, err
	}

	return NewStaticTokenInfoData(data.Tokens), nil
}

// StaticTokenInfoData holds an array of TokenInfoData.
type StaticTokenInfoData struct {
	Tokens map[string]string
}

// TypeVersion implements InfoData.
func (s *StaticTokenInfoData) TypeVersion() infodata.TypeVersion {
	return StaticTokenDataType
}

// Marshal implements InfoData
func (s *StaticTokenInfoData) Marshal() ([]byte, error) {
	return json.Marshal(&StaticTokenJSONData{s.Tokens})
}

// Append appends the tokens from the provided StaticTokenInfoData to this StaticTokenInfoData.
func (s *StaticTokenInfoData) Append(staticTokenInfoData *StaticTokenInfoData) {
	for username, token := range staticTokenInfoData.Tokens {
		s.Tokens[username] = token
	}
}

// RemoveTokens removes tokens with the provided usernames from this StaticTokenInfoData.
func (s *StaticTokenInfoData) RemoveTokens(usernames ...string) {
	for _, username := range usernames {
		delete(s.Tokens, username)
	}
}

// NewStaticTokenInfoData creates a new StaticTokenInfoData with the provided tokens.
func NewStaticTokenInfoData(tokens map[string]string) *StaticTokenInfoData {
	return &StaticTokenInfoData{tokens}
}
