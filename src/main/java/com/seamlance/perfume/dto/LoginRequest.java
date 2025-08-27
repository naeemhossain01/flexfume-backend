package com.seamlance.perfume.dto;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.seamlance.perfume.constants.EntityConstant;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotEmpty;
import lombok.Data;

@JsonIgnoreProperties(ignoreUnknown = true)
public class LoginRequest {
    @NotBlank(message = EntityConstant.PHONE_NUMBER_REQUIRED)
    private String phoneNumber;

    @NotEmpty(message = EntityConstant.PASSWORD_REQUIRED)
    private String password;

    public String getPhoneNumber() {
        return phoneNumber;
    }

    public void setPhoneNumber(String phoneNumber) {
        this.phoneNumber = phoneNumber;
    }

    public String getPassword() {
        return password;
    }

    public void setPassword(String password) {
        this.password = password;
    }
}
