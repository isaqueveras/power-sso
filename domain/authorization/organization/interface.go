// Copyright (c) 2024 Isaque Veras
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package organization

// IOrganization define an interface for data layer access methods
type IOrganization interface {
	Create(*Organization) error
}