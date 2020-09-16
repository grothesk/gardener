// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-License-Identifier: Apache-2.0

package kubernetes

// DeleteManifestOption is some configuration that modifies options for a delete request.
type DeleteManifestOption interface {
	// MutateDeleteOptions applies this configuration to the given delete options.
	MutateDeleteManifestOptions(opts *DeleteManifestOptions)
}

// DeleteOptions contains options for delete requests
type DeleteManifestOptions struct {
	// TolerateErrorFuncs are functions for which errors are tolerated.
	TolerateErrorFuncs []TolerateErrorFunc
}
