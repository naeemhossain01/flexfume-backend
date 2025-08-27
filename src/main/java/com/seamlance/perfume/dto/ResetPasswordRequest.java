package com.seamlance.perfume.dto;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.seamlance.perfume.constants.ErrorConstant;
import jakarta.validation.constraints.NotNull;

@JsonIgnoreProperties(ignoreUnknown = true)
public class ResetPasswordRequest {
    @NotNull(message = ErrorConstant.PHONE_NUMBER_NOT_FOUND)
    private String phoneNumber;

    @NotNull(message = ErrorConstant.INVALID_OTP)
    private String otp;
    private String newPassword;
    private String confirmPassword;

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

    public String getNewPassword() {
        return newPassword;
    }

    public void setNewPassword(String newPassword) {
        this.newPassword = newPassword;
    }

    public String getConfirmPassword() {
        return confirmPassword;
    }

    public void setConfirmPassword(String confirmPassword) {
        this.confirmPassword = confirmPassword;
    }
}
