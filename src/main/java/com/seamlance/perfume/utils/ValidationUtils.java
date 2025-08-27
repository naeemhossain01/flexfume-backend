package com.seamlance.perfume.utils;

import com.fasterxml.jackson.databind.JsonNode;
import com.seamlance.perfume.constants.Constant;
import com.seamlance.perfume.constants.ErrorConstant;
import com.seamlance.perfume.exception.InvalidRequestsException;

public class ValidationUtils {
    public static void otpResponseValidation(JsonNode response) {
        if(response == null) {
            throw new InvalidRequestsException(ErrorConstant.INTERNAL_SERVER_ERROR);
        }

        if(response.has(Constant.OTP_SMS_RESPONSE_CODE_KEY) && !(response.get(Constant.OTP_SMS_RESPONSE_CODE_KEY).intValue() == 202)) {
            throw new InvalidRequestsException(ErrorConstant.INTERNAL_SERVER_ERROR);
        }
    }
}
