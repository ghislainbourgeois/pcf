// Copyright 2019 free5GC.org
//
// SPDX-License-Identifier: Apache-2.0
//

package producer

import (
	"net/http"
	"strconv"

	"github.com/omec-project/http_wrapper"
	"github.com/omec-project/openapi/models"
	"github.com/omec-project/pcf/context"
	"github.com/omec-project/pcf/logger"
)

type UEAmPolicy struct {
	PolicyAssociationID string
	AccessType          models.AccessType
	Rfsp                string
	Triggers            []models.RequestTrigger
	/*Service Area Restriction */
	RestrictionType models.RestrictionType
	Areas           []models.Area
	MaxNumOfTAs     int32
}

type UEAmPolicys []UEAmPolicy

func HandleOAMGetAmPolicyRequest(request *http_wrapper.Request) *http_wrapper.Response {
	// step 1: log
	logger.OamLog.Infof("Handle OAMGetAmPolicy")

	// step 2: retrieve request
	supi := request.Params["supi"]

	// step 3: handle the message
	response, problemDetails := OAMGetAmPolicyProcedure(supi)

	// step 4: process the return value from step 3
	if response != nil {
		// status code is based on SPEC, and option headers
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(http.StatusForbidden, nil, problemDetails)
}

func OAMGetAmPolicyProcedure(supi string) (response *UEAmPolicys, problemDetails *models.ProblemDetails) {
	logger.OamLog.Infof("Handle OAM Get Am Policy")
	response = &UEAmPolicys{}
	pcfSelf := context.PCF_Self()

	if val, exists := pcfSelf.UePool.Load(supi); exists {
		ue := val.(*context.UeContext)
		for _, amPolicy := range ue.AMPolicyData {
			ueAmPolicy := UEAmPolicy{
				PolicyAssociationID: amPolicy.PolAssoId,
				AccessType:          amPolicy.AccessType,
				Rfsp:                strconv.Itoa(int(amPolicy.Rfsp)),
				Triggers:            amPolicy.Triggers,
			}
			if amPolicy.ServAreaRes != nil {
				servAreaRes := amPolicy.ServAreaRes
				ueAmPolicy.RestrictionType = servAreaRes.RestrictionType
				ueAmPolicy.Areas = servAreaRes.Areas
				ueAmPolicy.MaxNumOfTAs = servAreaRes.MaxNumOfTAs
			}
			*response = append(*response, ueAmPolicy)
		}
		return response, nil
	} else {
		problemDetails = &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "CONTEXT_NOT_FOUND",
		}
		return nil, problemDetails
	}
}
