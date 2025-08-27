package com.seamlance.perfume.dto;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.seamlance.perfume.constants.EntityConstant;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;

@JsonIgnoreProperties(ignoreUnknown = true)
public class ChangePasswordRequest {
    @NotBlank(message = EntityConstant.PASSWORD_REQUIRED)
    private String currentPassword;

    @NotBlank(message = EntityConstant.PASSWORD_REQUIRED)
    @Size(min = 8, message = "New password must be at least 8 characters long")
    private String newPassword;

    @NotBlank(message = EntityConstant.PASSWORD_REQUIRED)
    private String confirmPassword;

    public String getCurrentPassword() {
        return currentPassword;
    }

    public void setCurrentPassword(String currentPassword) {
        this.currentPassword = currentPassword;
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