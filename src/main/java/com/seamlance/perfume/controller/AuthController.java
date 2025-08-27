package com.seamlance.perfume.controller;

import com.seamlance.perfume.constants.Constant;
import com.seamlance.perfume.delegate.UserDelegate;
import com.seamlance.perfume.dto.ApiResponse;
import com.seamlance.perfume.dto.LoginRequest;
import com.seamlance.perfume.dto.OtpRequest;
import com.seamlance.perfume.info.UserInfo;
import com.seamlance.perfume.service.OtpService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/v1/auth")
public class AuthController {

    @Autowired
    private UserDelegate userDelegate;

    @Autowired
    private OtpService otpService;

    @PostMapping("/register")
    public ApiResponse<UserInfo> registerUser(@RequestBody UserInfo userInfo) {
        otpService.isNumberIsRegistered(userInfo.getPhoneNumber());
        UserInfo toReturn = userDelegate.registerUser(userInfo);
        return ApiResponse.success(toReturn);
    }

    @PostMapping("/login")
    public ApiResponse<String> login(@RequestBody LoginRequest loginRequest) {
        String token = userDelegate.loginUser(loginRequest);
        return ApiResponse.success(token);
    }

    @GetMapping("/init-register")
    public ApiResponse<String> sendOpt(@RequestParam String phoneNumber, @RequestParam String type) {
        otpService.isPhoneNumberVerified(phoneNumber);
        userDelegate.validateAlreadyHaveAccount(phoneNumber);
        otpService.send(phoneNumber, type, Constant.USER_REGISTRATION_TYPE);

        return ApiResponse.success(Constant.SMS_SEND_SUCCESSFULLY);
    }

    @PostMapping("/verify-otp")
    public ApiResponse<String> validateOtp(@RequestBody OtpRequest otpRequest) {
        otpService.verifyOtpForUser(otpRequest);
        otpService.markPhoneNumberAsValid(otpRequest);

        return ApiResponse.success(Constant.USER_VERIFICATION_SUCCESSFUL);
    }
}
