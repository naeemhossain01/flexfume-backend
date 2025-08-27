package com.seamlance.perfume.controller;

import java.util.List;

import com.seamlance.perfume.dto.ResetPasswordRequest;
import com.seamlance.perfume.service.OtpService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.bind.annotation.RequestParam;

import com.seamlance.perfume.constants.Constant;
import com.seamlance.perfume.delegate.UserDelegate;
import com.seamlance.perfume.dto.ApiResponse;
import com.seamlance.perfume.dto.ChangePasswordRequest;
import com.seamlance.perfume.info.UserInfo;

import jakarta.validation.Valid;

@RestController
@RequestMapping("/api/v1/user")
public class UserController {
    @Autowired
    private UserDelegate userDelegate;

    @Autowired
    private OtpService otpService;

    @GetMapping("/{id}")
    public ApiResponse<UserInfo> getUser(@PathVariable String id) {
        return ApiResponse.success(userDelegate.getUser(id));
    }

    @GetMapping("/profile")
    public ApiResponse<UserInfo> getCurrentUserProfile() {
        return ApiResponse.success(userDelegate.getCurrentUserProfile());
    }

    @PutMapping("/{id}")
    public ApiResponse<UserInfo> updateUser(@PathVariable String id, @RequestBody UserInfo userInfo) {
        return ApiResponse.success(userDelegate.updateUser(id, userInfo));
    }

    @GetMapping()
    @PreAuthorize("hasAuthority('ADMIN')")
    public ApiResponse<List<UserInfo>> getAllUser() {
        return ApiResponse.success(userDelegate.getAllUser());
    }

    @PostMapping("/change-password")
    public ApiResponse<String> changePassword(@Valid @RequestBody ChangePasswordRequest changePasswordRequest) {
        userDelegate.changePassword(changePasswordRequest);
        return ApiResponse.success(Constant.PASSWORD_CHANGED_SUCCESSFULLY);
    }

    @GetMapping("/reset-password-request")
    public ApiResponse<String> resetPasswordRequest(@RequestParam String phoneNumber, @RequestParam String type) {
        userDelegate.getUserByPhoneNumber(phoneNumber);
        otpService.send(phoneNumber, type, Constant.USER_PASSWORD_RESET_TYPE);

        return ApiResponse.success(Constant.PASSWORD_OTP_SENT);
    }

    @PostMapping("/reset-password")
    public ApiResponse<String> resetPassword(@Valid @RequestBody ResetPasswordRequest resetPasswordRequest) {
        otpService.verifyResetPasswordOtp(resetPasswordRequest.getPhoneNumber(), resetPasswordRequest.getOtp());
        userDelegate.resetPassword(resetPasswordRequest.getPhoneNumber(), resetPasswordRequest.getNewPassword(), resetPasswordRequest.getConfirmPassword());

        return ApiResponse.success(Constant.PASSWORD_RESET);
    }
}
