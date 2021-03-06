package tc

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
)

type Alert struct {
	Text  string `json:"text"`
	Level string `json:"level"`
}

type Alerts struct {
	Alerts []Alert `json:"alerts"`
}

func CreateErrorAlerts(errs ...error) Alerts {
	alerts := []Alert{}
	for _, err := range errs {
		if err != nil {
			alerts = append(alerts, Alert{err.Error(), ErrorLevel.String()})
		}
	}
	return Alerts{alerts}
}

func CreateAlerts(level AlertLevel, messages ...string) Alerts {
	alerts := []Alert{}
	for _, message := range messages {
		alerts = append(alerts, Alert{message, level.String()})
	}
	return Alerts{alerts}
}

func GetHandleErrorFunc(w http.ResponseWriter, r *http.Request) func(err error, status int) {
	return func(err error, status int) {
		log.Errorf("%v %v\n", r.RemoteAddr, err)
		errBytes, jsonErr := json.Marshal(CreateErrorAlerts(err))
		if jsonErr != nil {
			log.Errorf("failed to marshal error: %s\n", jsonErr)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, http.StatusText(http.StatusInternalServerError))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		fmt.Fprintf(w, "%s", errBytes)
	}
}
