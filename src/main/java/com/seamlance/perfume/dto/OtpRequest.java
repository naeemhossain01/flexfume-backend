package com.seamlance.perfume.dto;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.seamlance.perfume.constants.Constant;
import com.seamlance.perfume.constants.ErrorConstant;
import jakarta.validation.constraints.NotNull;

@JsonIgnoreProperties(ignoreUnknown = true)
public class OtpRequest {

    @NotNull(message = ErrorConstant.PHONE_NUMBER_NOT_FOUND)
    private String phoneNumber;

    @NotNull(message = ErrorConstant.INVALID_OTP)
    private String otp;
    private boolean verified = false;

    public String getPhoneNumber() {
        return phoneNumber;
    }

    public void setPhoneNumber(String phoneNumber) {
        this.phoneNumber = phoneNumber;
    }

    public String getOtp() {
        return otp;
    }

    public void setOtp(String otp) {
        this.otp = otp;
    }

    public boolean isVerified() {
        return verified;
    }

    public void setVerified(boolean verified) {
        this.verified = verified;
    }
}
