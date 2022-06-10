// Copyright © 2021 Alibaba Group Holding Ltd.
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

package utils

import (
	"fmt"
	"regexp"
	"time"

	v2 "github.com/sealerio/sealer/types/api/v2"
)

func Retry(tryTimes int, trySleepTime time.Duration, action func() error) error {
	var err error
	for i := 0; i < tryTimes; i++ {
		err = action()
		if err == nil {
			return nil
		}

		time.Sleep(trySleepTime * time.Duration(2*i+1))
	}
	return fmt.Errorf("retry action timeout: %v", err)
}

// ConfirmOperation confirm whether to continue with the operation, typing yes will return true.
func ConfirmOperation(promptInfo string) (bool, error) {
	var yesRx = regexp.MustCompile("^(?:y(?:es)?)$")
	var noRx = regexp.MustCompile("^(?:n(?:o)?)$")
	var input string
	for {
		fmt.Printf(promptInfo + " Yes [y/yes], No [n/no] : ")
		_, err := fmt.Scanln(&input)
		if err != nil {
			return false, err
		}
		if yesRx.MatchString(input) {
			break
		}
		if noRx.MatchString(input) {
			return false, nil
		}
	}
	return true, nil
}

// GetIPListFromHosts returns a list of IP for each 'Hosts' object
// IPs are obtained via:
// * ipList: []string, each string represents a single IP, such as 192.168.0.1
// Other formats are not supported.
func GetIPListFromHosts(hosts []v2.Host) []string {
	var hostsIPList []string
	for _, host := range hosts {
		// handle ipList
		for _, ip := range host.IPS {
			hostsIPList = append(hostsIPList, ip)
		}
	}

	return hostsIPList
}
