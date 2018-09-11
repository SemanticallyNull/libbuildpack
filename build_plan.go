/*
 * Copyright 2018 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package libbuildpack

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

// BuildPlan represents the dependencies contributed by a build.
type BuildPlan map[string]BuildPlanDependency

// String makes BuildPlan satisfy the Stringer interface.
func (p BuildPlan) String() string {
	var entries []string

	for k, v := range p {
		entries = append(entries, fmt.Sprintf("%s: %s", k, v))
	}

	return fmt.Sprintf("BuildPlan{ %s }", strings.Join(entries, ", "))
}

// BuildPlanDependency represents a dependency in a build.
type BuildPlanDependency map[string]interface{}

// String makes BuildPlanDependency satisfy the Stringer interface.
func (d BuildPlanDependency) String() string {
	var entries []string

	for k, v := range d {
		entries = append(entries, fmt.Sprintf("%s: %v", k, v))
	}

	return fmt.Sprintf("BuildPlanDependency{ %s }", strings.Join(entries, ", "))
}

// Version returns the version of this dependency.  Returns an error if the version is missing or invalid.
func (d BuildPlanDependency) Version() (string, error) {
	v, ok := d["version"].(string)
	if !ok {
		return "", fmt.Errorf("version key does not exist")
	}

	return v, nil
}

// DefaultBuildPlan creates a new instance of BuildPlan, extracting the contents from stdin.
func DefaultBuildPlan(logger Logger) (BuildPlan, error) {
	return NewBuildPlan(os.Stdin, logger)
}

// NewBuildPlan creates a new instance of BuildPlan from a specified io.Reader.  Returns an error if the contents of the
// Reader are not valid TOML.
func NewBuildPlan(in io.Reader, logger Logger) (BuildPlan, error) {
	var p BuildPlan

	if _, err := toml.DecodeReader(in, &p); err != nil {
		return BuildPlan{}, err
	}

	logger.Debug("BuildPlan: %s", p)
	return p, nil
}
