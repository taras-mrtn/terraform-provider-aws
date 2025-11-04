// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pinpointsmsvoicev2

// Exports for use in tests only.
var (
	ResourceConfigurationSet = newConfigurationSetResource
	ResourceOptOutList       = newOptOutListResource
	ResourcePhoneNumber      = newPhoneNumberResource
	ResourceSenderId         = newSenderIdResource

	FindConfigurationSetByID = findConfigurationSetByID
	FindOptOutListByID       = findOptOutListByID
	FindPhoneNumberByID      = findPhoneNumberByID
	FindSenderIdByARN        = FindSenderIdByARN
)
